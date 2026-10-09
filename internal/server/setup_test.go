package server

import (
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// setupTable seats n players at a game with the given options and marks
// them ready.
func setupTable(t *testing.T, h *Hub, n int, opts protocol.Options) []*player {
	t.Helper()
	host, _ := arrive(t, h, "P0", "")
	host.say(protocol.TypeCreate, protocol.Create{Seats: n, Options: &opts})
	var tb protocol.Table
	host.next(protocol.TypeTable, &tb)
	players := []*player{host}
	for i := 1; i < n; i++ {
		p, _ := arrive(t, h, "P"+strconv.Itoa(i), "")
		p.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
		p.next(protocol.TypeTable, nil)
		players = append(players, p)
	}
	for _, p := range players {
		p.say(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	return players
}

// dealt returns the characters of the game's first update.
func dealt(p *player) []engine.CardRef {
	p.t.Helper()
	var u protocol.Update
	p.next(protocol.TypeUpdate, &u)
	var out []engine.CardRef
	for _, pl := range u.View.Players {
		out = append(out, pl.Character.Card)
	}
	return out
}

var baseCharacters = []engine.CardRef{"blue_baby", "cain", "eden", "eve", "isaac", "judas", "lazarus", "lilith", "maggy", "samson", "the_forgotten"}

// GS-02: host bans are never dealt.
func TestHostBans(t *testing.T) {
	h := startHub(t)
	// Ban all but Isaac and Samson: those two must be dealt.
	bans := slices.DeleteFunc(slices.Clone(baseCharacters), func(c engine.CardRef) bool { return c == "isaac" || c == "samson" })
	players := setupTable(t, h, 2, protocol.Options{HostBans: bans})
	got := dealt(players[0])
	for _, c := range got {
		if slices.Contains(bans, c) {
			t.Errorf("banned %s was dealt", c)
		}
	}
	if len(got) != 2 || !slices.Contains(got, "isaac") || !slices.Contains(got, "samson") {
		t.Fatalf("dealt %v, want Isaac and Samson", got)
	}
}

// GS-04: ban rounds go in snake order; bans leave the pool.
func TestBanRounds(t *testing.T) {
	h := startHub(t)
	players := setupTable(t, h, 3, protocol.Options{BanRounds: 2})
	want := []int{0, 1, 2, 2, 1, 0}
	var banned []engine.CardRef
	for i, seat := range want {
		var s protocol.Setup
		players[0].next(protocol.TypeSetup, &s)
		if s.Phase != protocol.PhaseBan || s.Turn != seat || s.Round != i/3+1 {
			t.Fatalf("ban %d: phase %s turn %d round %d, want turn %d", i, s.Phase, s.Turn, s.Round, seat)
		}
		if i == 0 {
			players[1].say(protocol.TypeBan, protocol.BanCard{Card: s.Pool[0]})
			if code := players[1].failure(); code != protocol.ErrNotYourTurn {
				t.Errorf("out of turn: %s", code)
			}
			players[0].say(protocol.TypeBan, protocol.BanCard{Card: "fly"})
			if code := players[0].failure(); code != protocol.ErrBadCard {
				t.Errorf("not a character: %s", code)
			}
		}
		players[seat].say(protocol.TypeBan, protocol.BanCard{Card: s.Pool[0]})
		banned = append(banned, s.Pool[0])
	}
	for _, c := range dealt(players[0]) {
		if slices.Contains(banned, c) {
			t.Errorf("banned %s was dealt", c)
		}
	}
}

// GS-04: when the ban timer runs out, that player makes no ban.
func TestBanTimer(t *testing.T) {
	old := banSecond
	banSecond = time.Millisecond
	defer func() { banSecond = old }()
	h := startHub(t)
	players := setupTable(t, h, 2, protocol.Options{BanRounds: 1, BanTimer: 5})
	var s protocol.Setup
	players[0].next(protocol.TypeSetup, &s)
	if s.Deadline == 0 {
		t.Error("no deadline shown")
	}
	// Nobody bans: both turns time out and the characters are dealt.
	if got := dealt(players[0]); len(got) != 2 {
		t.Fatalf("dealt %v", got)
	}
}

// GS-01, GS-03: a draft offers each player their own characters to pick
// from.
func TestDraft(t *testing.T) {
	h := startHub(t)
	players := setupTable(t, h, 2, protocol.Options{Picking: protocol.PickDraft, DraftSize: 3})
	var offers [2][]engine.CardRef
	for i, p := range players {
		var s protocol.Setup
		p.next(protocol.TypeSetup, &s)
		if s.Phase != protocol.PhasePick || len(s.Offers) != 3 {
			t.Fatalf("seat %d: phase %s offers %v", i, s.Phase, s.Offers)
		}
		offers[i] = s.Offers
	}
	for _, c := range offers[0] {
		if slices.Contains(offers[1], c) {
			t.Fatalf("offers overlap: %v and %v", offers[0], offers[1])
		}
	}
	players[0].say(protocol.TypePick, protocol.PickCard{Card: offers[1][0]})
	if code := players[0].failure(); code != protocol.ErrBadCard {
		t.Errorf("picking another's offer: %s", code)
	}
	players[0].say(protocol.TypePick, protocol.PickCard{Card: offers[0][2]})
	players[1].say(protocol.TypePick, protocol.PickCard{Card: offers[1][1]})
	got := dealt(players[1])
	if got[0] != offers[0][2] || got[1] != offers[1][1] {
		t.Errorf("characters %v, want the picks %s and %s", got, offers[0][2], offers[1][1])
	}
}

func TestOptionChecks(t *testing.T) {
	h := startHub(t)
	ann, _ := arrive(t, h, "Ann", "")
	for name, o := range map[string]protocol.Options{
		"draft of 9":       {Picking: protocol.PickDraft, DraftSize: 9},
		"unknown picking":  {Picking: "auction"},
		"4 ban rounds":     {BanRounds: 4},
		"1 second timer":   {BanRounds: 1, BanTimer: 1},
		"not a character":  {HostBans: []engine.CardRef{"fly"}},
		"too few to draft": {Picking: protocol.PickDraft, DraftSize: 5, BanRounds: 3},
	} {
		ann.say(protocol.TypeCreate, protocol.Create{Seats: 4, Options: &o})
		if code := ann.failure(); code != protocol.ErrBadSetup {
			t.Errorf("%s: %s", name, code)
		}
	}
}
