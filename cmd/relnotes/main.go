// Command relnotes prints the release notes of an app version: the app
// and rules engine versions, then the version's section of CHANGELOG.md
// (8.2). A trial version (1.2.3-test) uses the Unreleased section.
//
//	relnotes -version 1.2.3 > notes.md
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// semver is a version like 1.2.3 or 1.2.3-test.1.
var semver = regexp.MustCompile(`^\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$`)

// ErrVersion means the tag is not a semantic version.
var ErrVersion = errors.New("not a version like 1.2.3 or 1.2.3-test")

func main() {
	version := flag.String("version", "", "app version, e.g. 1.2.3 (the tag without v)")
	changelog := flag.String("changelog", "CHANGELOG.md", "changelog file")
	flag.Parse()
	log, err := os.ReadFile(*changelog)
	if err == nil {
		var out string
		if out, err = notes(log, *version, engine.Version); err == nil {
			_, err = fmt.Print(out)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "relnotes:", err)
		os.Exit(1)
	}
}

// notes builds the release notes from the changelog.
func notes(changelog []byte, version, engineVersion string) (string, error) {
	if !semver.MatchString(version) {
		return "", fmt.Errorf("%w: %q", ErrVersion, version)
	}
	heading := version
	if strings.Contains(version, "-") {
		heading = "Unreleased"
	}
	body, ok := section(changelog, heading)
	if !ok {
		return "", fmt.Errorf("CHANGELOG.md has no section %q", "## "+heading)
	}
	return fmt.Sprintf("**App %s · Rules engine %s**\n\n%s\n", version, engineVersion, body), nil
}

// section returns the text under "## <heading>" up to the next "## ".
func section(changelog []byte, heading string) (string, bool) {
	var b strings.Builder
	in, found := false, false
	sc := bufio.NewScanner(bytes.NewReader(changelog))
	for sc.Scan() {
		line := sc.Text()
		if rest, ok := strings.CutPrefix(line, "## "); ok {
			title, _, _ := strings.Cut(rest, " ")
			in = strings.Trim(title, "[]") == heading
			found = found || in
			continue
		}
		if in {
			b.WriteString(line + "\n")
		}
	}
	return strings.TrimSpace(b.String()), found
}
