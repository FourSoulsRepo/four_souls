package server

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// startSoft starts a Soft game of n players (ADR 020, items 3 to 5).
func startSoft(t *testing.T, h *Hub, n int, opts protocol.Options) *driver {
	t.Helper()
	opts.Rollback = protocol.RollbackSoft
	players := startTableWith(t, h, protocol.Create{Seats: n, Options: &opts})
	tb := &driver{t: t, players: players, last: make([]protocol.Update, n)}
	tb.updates(-1) // the start
	return tb
}

// toTarget plays until a Soft rollback has a target and returns the
// newest one, with its author.
func (tb *driver) toTarget() protocol.RollbackStep {
	tb.t.Helper()
	for range 200 {
		if s := tb.last[0].RollbackLast; s != nil {
			return *s
		}
		tb.act()
	}
	tb.t.Fatal("no Soft target came")
	return protocol.RollbackStep{}
}

// others are the seats other than the author, in seat order.
func (tb *driver) others(author int) []int {
	var out []int
	for i := range tb.players {
		if i != author {
			out = append(out, i)
		}
	}
	return out
}

// openVote expects every player's update to show the open vote.
func (tb *driver) openVote(step, seat, yes int) protocol.RollbackVote {
	tb.t.Helper()
	got := tb.updates(-1)
	for i, u := range got {
		v := u.RollbackVote
		if v == nil || v.Step != step || v.Seat != seat || u.Rollback != nil {
			tb.t.Fatalf("player %d: vote %+v, rollback %+v; want a vote by %d on step %d", i, v, u.Rollback, seat, step)
		}
		n := 0
		for _, b := range v.Votes {
			if b.Yes {
				n++
			}
		}
		if n != yes || v.Need != need(len(tb.players)-1) || len(v.Voters) != len(tb.players)-1 {
			tb.t.Fatalf("player %d: %d yes (want %d), need %d of %v", i, n, yes, v.Need, v.Voters)
		}
		if u.Deadline != 0 {
			tb.t.Errorf("player %d: the response timer runs during the vote", i)
		}
	}
	return *got[0].RollbackVote
}

// rolledBack expects the rollback update: seat rolled back to before step.
func (tb *driver) rolledBack(step, seat int) {
	tb.t.Helper()
	for i, u := range tb.updates(-1) {
		if u.Rollback == nil || u.Rollback.To != step || u.Rollback.Seat != seat || u.RollbackVote != nil {
			tb.t.Fatalf("player %d: rollback %+v, vote %+v; want %d to before %d", i, u.Rollback, u.RollbackVote, seat, step)
		}
	}
}

// ended expects the update of a vote that ended with no rollback.
func (tb *driver) ended(step int, why string) {
	tb.t.Helper()
	for i, u := range tb.updates(-1) {
		if u.RollbackEnd == nil || u.RollbackEnd.Step != step || u.RollbackEnd.Why != why || u.RollbackVote != nil || u.Rollback != nil {
			tb.t.Fatalf("player %d: end %+v, vote %+v, rollback %+v; want %s", i, u.RollbackEnd, u.RollbackVote, u.Rollback, why)
		}
	}
}

// TestSoftMajority: in Soft the author never rolls back alone; another
// player opens the vote with their yes, and more than half of the
// others decide: 2 players, the opponent; 3, both others; 4, two of
// three (ADR 020, items 5.1 to 5.6).
func TestSoftMajority(t *testing.T) {
	for _, n := range []int{2, 3, 4} {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			h := startHub(t)
			tb := startSoft(t, h, n, protocol.Options{})
			if u := tb.last[0]; u.RollbackMode != protocol.RollbackSoft {
				t.Fatalf("mode %q", u.RollbackMode)
			}
			target := tb.toTarget()
			author, others := target.Seat, tb.others(target.Seat)
			tb.players[author].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step})
			if code := tb.players[author].failure(); code != protocol.ErrRollback {
				t.Errorf("the author rolled back alone: %s", code)
			}
			opener := others[0]
			tb.players[opener].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step})
			switch n {
			case 3:
				tb.openVote(target.Step, opener, 1)
				tb.players[others[1]].say(protocol.TypeRollbackVote, protocol.RollbackVoteAnswer{Yes: true})
			case 4:
				// One no leaves the majority in reach: two of three.
				tb.openVote(target.Step, opener, 1)
				tb.players[others[1]].say(protocol.TypeRollbackVote, protocol.RollbackVoteAnswer{Yes: false})
				tb.openVote(target.Step, opener, 1)
				tb.players[others[2]].say(protocol.TypeRollbackVote, protocol.RollbackVoteAnswer{Yes: true})
			}
			tb.rolledBack(target.Step, opener)
		})
	}
}

// TestSoftDeclinedAndAskedAgain: a no from one of two voters declines
// the vote at once; the author then asks, the game waits, a voter
// changes their mind, and the rollback passes (items 5.2, 5.8, 5.12,
// 5.14).
func TestSoftDeclinedAndAskedAgain(t *testing.T) {
	h := startHub(t)
	tb := startSoft(t, h, 3, protocol.Options{})
	target := tb.toTarget()
	author, others := target.Seat, tb.others(target.Seat)
	tb.players[others[0]].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step})
	tb.openVote(target.Step, others[0], 1)
	tb.players[others[1]].say(protocol.TypeRollbackVote, protocol.RollbackVoteAnswer{Yes: false})
	tb.ended(target.Step, protocol.VoteDeclined)
	if s := tb.last[0].RollbackLast; s == nil || s.Step < target.Step {
		t.Fatalf("the declined target is gone: %+v", s)
	}
	// The author asks: a vote with no yes yet.
	tb.players[author].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step, Ask: true})
	v := tb.openVote(target.Step, author, 0)
	if v.Author != author || slices.Contains(v.Voters, author) {
		t.Errorf("vote %+v: the author %d votes", v, author)
	}
	// The game waits; the author has no vote; nobody opens a second one.
	for i, u := range tb.last {
		if len(u.Allowed) > 0 {
			tb.players[i].say(protocol.TypeIntent, protocol.Intent{Intent: u.Allowed[0]})
			if code := tb.players[i].failure(); code != protocol.ErrRollback {
				t.Errorf("an intent during the vote: %s", code)
			}
		}
	}
	tb.players[author].say(protocol.TypeRollbackVote, protocol.RollbackVoteAnswer{Yes: true})
	if code := tb.players[author].failure(); code != protocol.ErrRollback {
		t.Errorf("the author voted: %s", code)
	}
	tb.players[others[0]].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step})
	if code := tb.players[others[0]].failure(); code != protocol.ErrRollback {
		t.Errorf("a second vote: %s", code)
	}
	tb.players[others[0]].say(protocol.TypeRollbackVote, protocol.RollbackVoteAnswer{Yes: true})
	tb.openVote(target.Step, author, 1)
	tb.players[others[1]].say(protocol.TypeRollbackVote, protocol.RollbackVoteAnswer{Yes: true})
	tb.rolledBack(target.Step, author)
}

// TestSoftCancel: only whoever opened the vote cancels it (item 5.13).
func TestSoftCancel(t *testing.T) {
	h := startHub(t)
	tb := startSoft(t, h, 3, protocol.Options{})
	target := tb.toTarget()
	author, others := target.Seat, tb.others(target.Seat)
	tb.players[others[0]].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step})
	tb.openVote(target.Step, others[0], 1)
	for _, s := range []int{author, others[1]} {
		tb.players[s].say(protocol.TypeRollbackCancel, struct{}{})
		if code := tb.players[s].failure(); code != protocol.ErrRollback {
			t.Errorf("seat %d cancelled another's vote: %s", s, code)
		}
	}
	tb.players[others[0]].say(protocol.TypeRollbackCancel, struct{}{})
	tb.ended(target.Step, protocol.VoteCancelled)
	// The game goes on.
	tb.act()
}

// TestSoftVoteTimeout: with the response timer on, the vote gets the
// same seconds and running out declines it; then the response timer
// runs again (items 5.9 to 5.11).
func TestSoftVoteTimeout(t *testing.T) {
	h := startHubWith(t, 20*time.Millisecond) // 60 "seconds" are 1.2 s
	tb := startSoft(t, h, 2, protocol.Options{ResponseTimer: 60})
	target := tb.toTarget()
	author := target.Seat
	opened := time.Now()
	tb.players[author].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step, Ask: true})
	v := tb.openVote(target.Step, author, 0)
	if v.Deadline == 0 {
		t.Fatal("the vote has no deadline with the response timer on")
	}
	tb.ended(target.Step, protocol.VoteTimeout)
	if d := time.Since(opened); d < 1100*time.Millisecond {
		t.Errorf("the vote ran out after %v, before its 60 seconds", d)
	}
	if tb.last[0].Deadline == 0 {
		t.Error("the response timer did not start again")
	}
}

// TestSoftNoVoteWithoutTimer: with no response timer the vote stays
// open until it is decided (item 5.9).
func TestSoftNoVoteWithoutTimer(t *testing.T) {
	h := startHubWith(t, time.Millisecond)
	tb := startSoft(t, h, 2, protocol.Options{})
	target := tb.toTarget()
	tb.players[target.Seat].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step, Ask: true})
	if v := tb.openVote(target.Step, target.Seat, 0); v.Deadline != 0 {
		t.Errorf("a deadline with no response timer: %d", v.Deadline)
	}
	time.Sleep(100 * time.Millisecond) // 100 "seconds"
	if u, ok := tb.players[0].latest(); ok {
		t.Errorf("the vote moved on by itself: %+v", u.RollbackEnd)
	}
}

// TestSoftBarrierInGame: once a barrier comes (a roll, here), the
// steps before it are no targets and a rollback to them is refused
// (item 3.1).
func TestSoftBarrierInGame(t *testing.T) {
	h := startHub(t)
	tb := startSoft(t, h, 2, protocol.Options{})
	target := tb.toTarget()
	for range 200 {
		if tb.last[0].RollbackFrom > target.Step {
			break
		}
		tb.act()
	}
	u := tb.last[0]
	if u.RollbackFrom <= target.Step || u.RollbackLast != nil && u.RollbackLast.Step <= target.Step {
		t.Fatalf("from %d, last %+v: the target %d stays", u.RollbackFrom, u.RollbackLast, target.Step)
	}
	other := tb.players[1-target.Seat]
	other.say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step})
	if code := other.failure(); code != protocol.ErrRollback {
		t.Errorf("a rollback past a barrier: %s", code)
	}
}

// TestSoftBarrierKinds: each event kind the engine calls a barrier
// moves the first Soft target past its step; every other kind leaves
// it (ADR 020, item 4). A lone pass or throw is never a target, and a
// private event is a barrier, but for a hand card put on the loot deck
// (item 4.22).
func TestSoftBarrierKinds(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "pkg", "rules_engine", "testdata", "event_kinds.txt"))
	if err != nil {
		t.Fatal(err)
	}
	attack := engine.Intent{Player: 1, Kind: engine.IntentAttack}
	room := func(events ...engine.Event) *Room {
		r := &Room{opts: protocol.Options{Rollback: protocol.RollbackSoft}, base: []byte("{}")}
		for i, ev := range [][]engine.Event{nil, events} {
			r.steps = i + 1
			r.keepLive(engine.StepRecord{Intent: attack, Events: ev})
		}
		return r
	}
	barriers := 0
	for _, k := range strings.Fields(string(data)) {
		r := room(engine.Event{Kind: engine.EventKind(k)})
		from, want := r.rollbackFrom(), 1
		if engine.Barrier(engine.Event{Kind: engine.EventKind(k)}) {
			want = 3
			barriers++
		}
		if from != want || len(allTargets(r)) != 3-want {
			t.Errorf("%s: from %d, targets %+v", k, from, allTargets(r))
		}
	}
	if barriers == 0 {
		t.Error("no event kind is a barrier")
	}
	if r := room(engine.Event{Kind: engine.EvLooted, Private: true}); r.rollbackFrom() != 3 {
		t.Error("a private event is no barrier")
	}
	if r := room(engine.Event{Kind: engine.EvMovedToDeck, Private: true}); r.rollbackFrom() != 1 {
		t.Error("a hand card put on the loot deck is a barrier")
	}
	for _, k := range []engine.IntentKind{engine.IntentPass, engine.IntentThrow} {
		r := room()
		r.live[1].rec.Intent.Kind = k
		if got := allTargets(r); len(got) != 1 || got[0] != (protocol.RollbackStep{Step: 1, Seat: 1}) {
			t.Errorf("intent %d: targets %+v", k, got)
		}
	}
}

// TestSoftLeaveEndsVote: a player who leaves pauses the game, and the
// pause ends the vote; once the game goes on without them, a vote
// needs a voter who is still there (items 5.4 and 5.7).
func TestSoftLeaveEndsVote(t *testing.T) {
	h := startHub(t)
	tb := startSoft(t, h, 2, protocol.Options{})
	target := tb.toTarget()
	author, other := target.Seat, 1-target.Seat
	tb.players[author].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step, Ask: true})
	tb.openVote(target.Step, author, 0)
	tb.players[other].c.Leave()
	u := tb.players[author].untilUpdate(func(u protocol.Update) bool { return u.RollbackEnd != nil })
	if u.RollbackEnd.Why != protocol.VoteLeft || u.RollbackVote != nil || u.Pause == nil {
		t.Fatalf("end %+v, vote %+v, pause %+v", u.RollbackEnd, u.RollbackVote, u.Pause)
	}
	// Kicked, the other leaves no voter: the author cannot even ask.
	tb.players[author].say(protocol.TypeVote, protocol.Vote{Kick: true})
	u = tb.players[author].untilUpdate(func(u protocol.Update) bool { return u.Pause == nil })
	if u.RollbackLast == nil || u.RollbackLast.Seat != author {
		t.Skip("the kicked seat's forced steps left no target") // a barrier came
	}
	tb.players[author].say(protocol.TypeRollback, protocol.RollbackTo{Step: u.RollbackLast.Step, Ask: true})
	var er protocol.Error
	tb.players[author].next(protocol.TypeError, &er)
	if er.Code != protocol.ErrRollback || er.Message != "error.rollback.no_voters" {
		t.Errorf("a vote with no voter: %+v", er)
	}
}

// allTargets lists every Soft target of the room now, oldest first. The
// protocol sends only the latest (rollback_last); tests that check the
// older ones read them here.
func allTargets(r *Room) []protocol.RollbackStep {
	var out []protocol.RollbackStep
	for i := range r.live {
		if r.isTarget(r.live[i]) {
			out = append(out, r.targetAt(i))
		}
	}
	return out
}
