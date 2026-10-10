package server

import (
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// untilUpdate reads updates until ok says so, or fails after 5 seconds.
func (p *player) untilUpdate(ok func(protocol.Update) bool) protocol.Update {
	p.t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var u protocol.Update
		p.next(protocol.TypeUpdate, &u)
		if ok(u) {
			return u
		}
	}
	p.t.Fatal("the expected update did not come")
	return protocol.Update{}
}

// TestDropPausesAndReconnectResumes (N-08).
func TestDropPausesAndReconnectResumes(t *testing.T) {
	h := startHub(t)
	ann, w := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)
	bob, _ := arrive(t, h, "Bob", "")
	bob.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bob.next(protocol.TypeTable, nil)
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	bob.say(protocol.TypeReady, protocol.Ready{Ready: true})
	ann.next(protocol.TypeUpdate, nil)

	ann.c.Leave()
	u := bob.untilUpdate(func(u protocol.Update) bool { return u.Pause != nil })
	if len(u.Pause.Away) != 1 || u.Pause.Away[0] != 0 {
		t.Fatalf("pause %+v", u.Pause)
	}
	bob.say(protocol.TypeIntent, protocol.Intent{Intent: engine.Intent{Kind: engine.IntentPass}})
	if code := bob.failure(); code != protocol.ErrPaused {
		t.Errorf("an intent during the pause: %s", code)
	}

	back, _ := arrive(t, h, "Ann", w.Token)
	u = back.untilUpdate(func(u protocol.Update) bool { return u.Pause == nil })
	if !u.Seats[0].Connected || u.View.Players[0].HandSize == 0 {
		t.Errorf("Ann did not get the full view back: %+v", u.Seats)
	}
	bob.untilUpdate(func(u protocol.Update) bool { return u.Pause == nil })
}

// TestVoteToKick: more than half of the connected players must vote
// kick; the kicked seat stays empty and the game goes on (N-08).
func TestVoteToKick(t *testing.T) {
	h := startHub(t)
	players := startTable(t, h, 3)
	for _, p := range players {
		p.next(protocol.TypeUpdate, nil)
	}
	players[2].c.Leave()
	players[0].untilUpdate(func(u protocol.Update) bool { return u.Pause != nil })

	players[0].say(protocol.TypeVote, protocol.Vote{Kick: true})
	u := players[0].untilUpdate(func(u protocol.Update) bool { return u.Pause != nil && len(u.Pause.Votes) == 1 })
	if u.Seats[2].Kicked {
		t.Fatal("one vote of two kicked")
	}
	players[1].say(protocol.TypeVote, protocol.Vote{Kick: false})
	players[1].untilUpdate(func(u protocol.Update) bool { return u.Pause != nil && len(u.Pause.Votes) == 2 })
	players[1].say(protocol.TypeVote, protocol.Vote{Kick: true}) // changes the vote
	u = players[0].untilUpdate(func(u protocol.Update) bool { return u.Pause == nil })
	if !u.Seats[2].Kicked {
		t.Fatalf("seats %+v", u.Seats)
	}

	// The game goes on: the kicked seat's turn passes by itself, inside
	// one run of automatic answers, so the turn number moves past it.
	start := u.View.Turn.Number
	passed := false
	var last [2]protocol.Update
	pending := map[int]protocol.Update{0: u} // player 0 already read this one
	deadline := time.Now().Add(5 * time.Second)
	for !passed && time.Now().Before(deadline) {
		for i, p := range players[:2] {
			up, ok := p.latest()
			if !ok {
				if up, ok = pending[i]; !ok {
					continue
				}
			}
			delete(pending, i)
			last[i] = up
			if up.View.Turn.Number >= start+3 { // every seat had a turn
				passed = true
			}
			if len(up.Allowed) == 0 || up.View.Over || up.Pause != nil {
				continue
			}
			in := up.Allowed[0]
			for _, x := range up.Allowed {
				if x.Kind == engine.IntentEndTurn {
					in = x // move turns along
				}
			}
			if in.Kind == engine.IntentDiscard {
				in.Objects = in.Objects[:up.View.Waiting.Count]
			}
			p.say(protocol.TypeIntent, protocol.Intent{Intent: in})
		}
		time.Sleep(5 * time.Millisecond)
	}
	if !passed {
		for i, up := range last {
			t.Errorf("player %d last saw step %d turn %d (from %d), active %d, waiting %+v, allowed %d, pause %v",
				i, up.Step, up.View.Turn.Number, start, up.View.Turn.Active, up.View.Waiting, len(up.Allowed), up.Pause)
		}
	}
}

// latest drains the pending messages and returns the newest update.
func (p *player) latest() (protocol.Update, bool) {
	var last protocol.Update
	got := false
	for {
		select {
		case msg := <-p.conn.out:
			e, err := protocol.Decode(msg)
			if err != nil || e.Type != protocol.TypeUpdate {
				continue
			}
			var u protocol.Update
			if e.Unpack(&u) == nil {
				last, got = u, true
			}
		default:
			return last, got
		}
	}
}
