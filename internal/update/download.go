package update

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Download fetches the archive of the last Check, checks its SHA-256
// against SHA256SUMS.txt and unpacks it. progress may be nil.
// A wrong or missing hash drops the file.
func (c *Client) Download(ctx context.Context, progress Progress) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.identity(); err != nil {
		return err
	}
	if !c.last.CanInstall || c.assetURL == "" {
		return ErrNotReady
	}
	c.payload = ""
	if c.sumsURL == "" {
		return ErrNoSums
	}
	asset := c.last.Asset
	sums, err := c.getSmall(ctx, c.sumsURL)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNoSums, err)
	}
	want, err := checksumFor(string(sums), asset)
	if err != nil {
		return err
	}
	dir, err := c.newWorkDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, asset)
	got, err := c.getFile(ctx, c.assetURL, path, progress)
	if err != nil {
		return errors.Join(err, removeFile(path))
	}
	if got != want {
		return errors.Join(fmt.Errorf("%w: %s: got %s, want %s", ErrChecksum, asset, got, want), removeFile(path))
	}
	payload, err := unpack(path, filepath.Join(dir, "payload"), c.cfg.Name, c.cfg.GOOS)
	if rmErr := removeFile(path); err == nil && rmErr != nil {
		err = rmErr
	}
	if err != nil {
		return err
	}
	c.payload = payload
	return nil
}

func removeFile(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

// checksumFor finds the hash of exactly this file name.
// Lines are "<hash>  <name>" as sha256sum writes them.
func checksumFor(sums, name string) (string, error) {
	sc := bufio.NewScanner(strings.NewReader(sums))
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) != 2 {
			continue
		}
		sum, file := fields[0], strings.TrimPrefix(fields[1], "*")
		if file != name {
			continue
		}
		if _, err := hex.DecodeString(sum); err != nil || len(sum) != 64 {
			return "", fmt.Errorf("%w: bad line for %s in %s", ErrChecksum, name, sumsName)
		}
		return strings.ToLower(sum), nil
	}
	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("update: %s: %w", sumsName, err)
	}
	return "", fmt.Errorf("%w: %s is not in %s", ErrChecksum, name, sumsName)
}

func (c *Client) getSmall(ctx context.Context, url string) ([]byte, error) {
	res, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer closeBody(res)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", res.StatusCode)
	}
	return io.ReadAll(io.LimitReader(res.Body, 1<<20))
}

// getFile saves url to dest and returns its SHA-256.
func (c *Client) getFile(ctx context.Context, url, dest string, progress Progress) (string, error) {
	res, err := c.get(ctx, url)
	if err != nil {
		return "", fmt.Errorf("update: download: %w", err)
	}
	defer closeBody(res)
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("update: download HTTP %d", res.StatusCode)
	}
	total := max(res.ContentLength, 0)
	if total > maxAssetBytes {
		return "", fmt.Errorf("%w: %d bytes", ErrTooLarge, total)
	}
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600) //nolint:gosec // our temp folder
	if err != nil {
		return "", fmt.Errorf("update: %w", err)
	}
	h := sha256.New()
	pw := &progressWriter{report: progress, total: total, every: c.ProgressEvery}
	pw.call()
	n, err := io.Copy(io.MultiWriter(f, h, pw), io.LimitReader(res.Body, maxAssetBytes+1))
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return "", fmt.Errorf("update: download: %w", err)
	}
	if n > maxAssetBytes {
		return "", fmt.Errorf("%w: over %d bytes", ErrTooLarge, int64(maxAssetBytes))
	}
	pw.call()
	return hex.EncodeToString(h.Sum(nil)), nil
}

// progressWriter counts bytes and reports them at most once per every.
type progressWriter struct {
	report Progress
	total  int64
	done   int64
	every  time.Duration
	last   time.Time
}

func (p *progressWriter) Write(b []byte) (int, error) {
	p.done += int64(len(b))
	if p.report != nil && time.Since(p.last) >= p.every {
		p.call()
	}
	return len(b), nil
}

func (p *progressWriter) call() {
	if p.report == nil {
		return
	}
	p.last = time.Now()
	p.report(p.done, p.total)
}

// newWorkDir drops the last download's folder and makes a new private one.
func (c *Client) newWorkDir() (string, error) {
	if c.workDir != "" {
		if err := os.RemoveAll(c.workDir); err != nil {
			return "", fmt.Errorf("update: %w", err)
		}
		c.workDir = ""
	}
	parent := c.TempDir
	if parent == "" {
		parent = os.TempDir()
	}
	dir, err := os.MkdirTemp(parent, c.cfg.Name+"-update-*")
	if err != nil {
		return "", fmt.Errorf("update: %w", err)
	}
	c.workDir = dir
	return dir, nil
}
