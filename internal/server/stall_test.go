package server

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"
)

// TestLogStall: a request that held the room past stallLimit is logged;
// a quick one is not (18.17).
func TestLogStall(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)
	r := &Room{}
	r.recInfo.game = "g7"
	r.logStall("a request", time.Now())
	if buf.Len() != 0 {
		t.Fatalf("a quick request was logged: %s", buf.String())
	}
	r.logStall("a request", time.Now().Add(-2*stallLimit))
	if got := buf.String(); !strings.Contains(got, "game g7: a request took") {
		t.Fatalf("no stall line: %q", got)
	}
}
