// Package lan finds hosted games on the local network (N-13, ADR 009).
// A joining app broadcasts a probe over UDP; a hosting app with a game
// waiting for players answers the sender.
package lan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"slices"
	"strconv"
	"time"
)

const (
	magic     = "four-souls"
	kindProbe = "probe"
	kindGames = "games"
	// maxMessage caps what is read and sent; anything longer is ignored.
	maxMessage = 8 << 10
)

// Game is one game that waits for players.
type Game struct {
	Host  string `json:"host"` // the nickname of who created it
	Seats int    `json:"seats"`
	Taken int    `json:"taken"`
	// Locked: joining needs the password (N-14).
	Locked bool `json:"locked,omitempty"`
}

// Answer is what a hosting app tells a probe.
type Answer struct {
	Magic    string `json:"magic"`
	Kind     string `json:"kind"`
	ID       string `json:"id"` // the server, so answers over two routes merge
	Protocol int    `json:"protocol"`
	App      string `json:"app"`
	Port     int    `json:"port"` // the TCP game port
	Games    []Game `json:"games"`
}

// Found is one server that answered, with the address to join it at.
type Found struct {
	Address  string `json:"address"`
	Protocol int    `json:"protocol"`
	App      string `json:"app"`
	Games    []Game `json:"games"`
}

type probe struct {
	Magic string `json:"magic"`
	Kind  string `json:"kind"`
}

// Listen answers probes on UDP port until ctx ends. games returns the
// answer now; ok false means no game waits, and nothing is sent.
func Listen(ctx context.Context, port int, games func() (Answer, bool)) error {
	var lc net.ListenConfig
	conn, err := lc.ListenPacket(ctx, "udp4", ":"+strconv.Itoa(port))
	if err != nil {
		return fmt.Errorf("lan: %w", err)
	}
	go func() {
		<-ctx.Done()
		_ = conn.Close() //nolint:errcheck // ends the read loop
	}()
	go answer(conn, games)
	return nil
}

func answer(conn net.PacketConn, games func() (Answer, bool)) {
	buf := make([]byte, maxMessage)
	for {
		n, from, err := conn.ReadFrom(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		var p probe
		if json.Unmarshal(buf[:n], &p) != nil || p.Magic != magic || p.Kind != kindProbe {
			continue
		}
		a, ok := games()
		if !ok {
			continue
		}
		a.Magic, a.Kind = magic, kindGames
		out, err := json.Marshal(a)
		if err != nil || len(out) > maxMessage {
			continue
		}
		_, _ = conn.WriteTo(out, from) //nolint:errcheck // the prober asks again
	}
}

// Find probes the network on UDP port and collects answers for wait.
func Find(ctx context.Context, port int, wait time.Duration) ([]Found, error) {
	var lc net.ListenConfig
	conn, err := lc.ListenPacket(ctx, "udp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("lan: %w", err)
	}
	defer conn.Close() //nolint:errcheck // read-only from here
	msg, err := json.Marshal(probe{Magic: magic, Kind: kindProbe})
	if err != nil {
		return nil, fmt.Errorf("lan: %w", err)
	}
	sent := 0
	for _, to := range targets(port) {
		if _, err := conn.WriteTo(msg, net.UDPAddrFromAddrPort(to)); err == nil {
			sent++
		}
	}
	if sent == 0 {
		return nil, errors.New("lan: no network to search")
	}
	deadline := time.Now().Add(wait)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	if err := conn.SetReadDeadline(deadline); err != nil {
		return nil, fmt.Errorf("lan: %w", err)
	}
	return collect(conn), nil
}

// collect reads answers until the read deadline; one server's answers
// merge, and a LAN address beats loopback.
func collect(conn net.PacketConn) []Found {
	var out []Found
	ids := map[string]int{}
	buf := make([]byte, maxMessage)
	for {
		n, from, err := conn.ReadFrom(buf)
		if err != nil {
			return out
		}
		var a Answer
		if json.Unmarshal(buf[:n], &a) != nil || a.Magic != magic || a.Kind != kindGames || a.Port < 1 || a.Port > 65535 {
			continue
		}
		udp, ok := from.(*net.UDPAddr)
		if !ok {
			continue
		}
		ip := udp.AddrPort().Addr().Unmap()
		f := Found{Address: net.JoinHostPort(ip.String(), strconv.Itoa(a.Port)), Protocol: a.Protocol, App: a.App, Games: a.Games}
		if i, seen := ids[a.ID]; seen {
			if !ip.IsLoopback() {
				out[i] = f
			}
			continue
		}
		ids[a.ID] = len(out)
		out = append(out, f)
	}
}

// targets are where a probe goes: the limited broadcast, each
// interface's own broadcast, and this computer (UC-16).
func targets(port int) []netip.AddrPort {
	p := uint16(port) //nolint:gosec // ports are 1 to 65535
	out := []netip.AddrPort{
		netip.AddrPortFrom(netip.AddrFrom4([4]byte{255, 255, 255, 255}), p),
		netip.AddrPortFrom(netip.AddrFrom4([4]byte{127, 0, 0, 1}), p),
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return out
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagBroadcast == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			if b, ok := broadcast(a); ok && !slices.Contains(out, netip.AddrPortFrom(b, p)) {
				out = append(out, netip.AddrPortFrom(b, p))
			}
		}
	}
	return out
}

// broadcast is the IPv4 broadcast address of an interface address.
func broadcast(a net.Addr) (netip.Addr, bool) {
	ipn, ok := a.(*net.IPNet)
	if !ok {
		return netip.Addr{}, false
	}
	ip4, mask := ipn.IP.To4(), ipn.Mask
	if len(mask) == net.IPv6len {
		mask = mask[12:] // an IPv4 mask in IPv6 form
	}
	if ip4 == nil || len(mask) != net.IPv4len {
		return netip.Addr{}, false
	}
	var b [4]byte
	for i := range b {
		b[i] = ip4[i] | ^mask[i]
	}
	return netip.AddrFrom4(b), true
}
