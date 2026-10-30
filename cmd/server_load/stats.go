package main

import (
	"cmp"
	"slices"
	"sync"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// collector gathers what the bots of one stage see. Its methods are safe
// for concurrent use.
type collector struct {
	mu          sync.Mutex
	handshake   []time.Duration // connect and hello until welcome
	response    []time.Duration // an intent until the next update its seat gets
	ping        []time.Duration // WebSocket ping until its pong
	updates     int
	bytes       int
	largest     int
	reconnects  int            // seats taken back on a new socket
	failedDials int            // tries to take a seat back that failed
	errors      map[string]int // error replies and broken orders, by code
	notices     map[string]int // notices, by kind
	server      *protocol.Welcome
}

func newCollector() *collector {
	return &collector{errors: map[string]int{}, notices: map[string]int{}}
}

func (c *collector) addHandshake(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handshake = append(c.handshake, d)
}

func (c *collector) addResponse(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.response = append(c.response, d)
}

func (c *collector) addPing(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ping = append(c.ping, d)
}

// addUpdate counts one update of size bytes.
func (c *collector) addUpdate(size int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.updates++
	c.bytes += size
	c.largest = max(c.largest, size)
}

// addReconnect counts a seat taken back on a new socket.
func (c *collector) addReconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reconnects++
}

// addFailedDial counts a try to take a seat back that failed.
func (c *collector) addFailedDial() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failedDials++
}

func (c *collector) addError(code string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errors[code]++
}

func (c *collector) addNotice(kind string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.notices[kind]++
}

// setServer keeps the first welcome: the server's name and versions.
func (c *collector) setServer(w protocol.Welcome) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.server == nil {
		c.server = &w
	}
}

// spread holds the figures of a list of samples.
type spread[T cmp.Ordered] struct {
	n                          int
	min, median, p95, p99, max T
}

// spreadOf sorts a copy of xs and picks the smallest, the median (the upper
// middle one for an even count), the 95th and 99th percentiles and the largest.
func spreadOf[T cmp.Ordered](xs []T) spread[T] {
	if len(xs) == 0 {
		return spread[T]{}
	}
	s := slices.Clone(xs)
	slices.Sort(s)
	n := len(s)
	rank := (n*95 + 99) / 100 // nearest rank: at least one
	rank99 := (n*99 + 99) / 100
	return spread[T]{
		n:      n,
		min:    s[0],
		median: s[n/2],
		p95:    s[rank-1],
		p99:    s[rank99-1],
		max:    s[n-1],
	}
}

// responses and pings are the two latency figures of a stage.
func (c *collector) responses() spread[time.Duration] {
	c.mu.Lock()
	defer c.mu.Unlock()
	return spreadOf(c.response)
}

func (c *collector) pings() spread[time.Duration] {
	c.mu.Lock()
	defer c.mu.Unlock()
	return spreadOf(c.ping)
}

// redials are the reconnects of a stage and its failed attempts.
func (c *collector) redials() (reconnects, failed int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reconnects, c.failedDials
}
