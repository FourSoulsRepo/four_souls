package lan

import (
	"context"
	"net"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// freePort is a UDP port nobody uses now.
func freePort(t *testing.T) int {
	t.Helper()
	var lc net.ListenConfig
	c, err := lc.ListenPacket(t.Context(), "udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	a, ok := c.LocalAddr().(*net.UDPAddr)
	if !ok {
		t.Fatalf("local address %v", c.LocalAddr())
	}
	port := a.Port
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	return port
}

func TestFindsAWaitingGame(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	port := freePort(t)
	var waiting atomic.Bool
	waiting.Store(true)
	games := func() (Answer, bool) {
		return Answer{ID: "s1", Protocol: 1, App: "1.0.0", Port: 4774, Games: []Game{{Host: "Ann", Seats: 3, Taken: 1}}}, waiting.Load()
	}
	if err := Listen(ctx, port, games); err != nil {
		t.Fatal(err)
	}
	found, err := Find(ctx, port, 300*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	// One server, however many routes the probe took.
	if len(found) != 1 {
		t.Fatalf("found %+v, want one server", found)
	}
	f := found[0]
	host, p, err := net.SplitHostPort(f.Address)
	if err != nil || p != "4774" || net.ParseIP(host) == nil {
		t.Errorf("address %q, want an IP and the game port 4774", f.Address)
	}
	if len(f.Games) != 1 || f.Games[0].Host != "Ann" || f.Protocol != 1 || f.App != "1.0.0" {
		t.Errorf("found %+v", f)
	}

	// No game waits: no answer (the game started).
	waiting.Store(false)
	if found, err := Find(ctx, port, 200*time.Millisecond); err != nil || len(found) != 0 {
		t.Errorf("found %+v (%v) with no waiting game, want none", found, err)
	}
}

func TestIgnoresStrangers(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	port := freePort(t)
	var asked atomic.Int32
	if err := Listen(ctx, port, func() (Answer, bool) {
		asked.Add(1)
		return Answer{ID: "s1", Port: 4774}, true
	}); err != nil {
		t.Fatal(err)
	}
	var d net.Dialer
	c, err := d.DialContext(ctx, "udp4", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close() //nolint:errcheck // test
	for _, msg := range []string{"hello", `{"magic":"other","kind":"probe"}`, `{"magic":"four-souls","kind":"games"}`} {
		if _, err := c.Write([]byte(msg)); err != nil {
			t.Fatal(err)
		}
	}
	// A real probe after the junk shows the junk was read and dropped.
	if found, err := Find(ctx, port, 200*time.Millisecond); err != nil || len(found) != 1 {
		t.Fatalf("found %+v (%v)", found, err)
	}
	// Only the real probe asks for the games, once per route it took.
	if n := asked.Load(); n < 1 {
		t.Errorf("asked %d times", n)
	}
}

func TestBroadcastAddress(t *testing.T) {
	_, ipn, err := net.ParseCIDR("192.168.1.37/24")
	if err != nil {
		t.Fatal(err)
	}
	ipn.IP = net.ParseIP("192.168.1.37")
	b, ok := broadcast(ipn)
	if !ok || b.String() != "192.168.1.255" {
		t.Errorf("broadcast %v %v, want 192.168.1.255", b, ok)
	}
}
