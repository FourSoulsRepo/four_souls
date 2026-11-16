package rulesengine

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

// cloneRules lists every field of Game and of each type inside it, in
// order, with how Game.Clone copies it (ADR 030, Decision 3.3):
//
//	Name   a plain value, copied with its struct; it holds no list, map
//	       or pointer
//	Name*  Clone gives the copy its own: a line in clone.go
//	Name!  shared or left out on purpose; Clone's comment says why
//
// A new field fails TestCloneCoversEveryField until it is here, and a
// field with a *, until Clone copies it.
var cloneRules = map[string]string{
	"Game": "RNG Players* Objects* Decks* Discards* Shop* Monsters* Rooms* " +
		"Bosses* Minions* Outside* BonusSouls* Queue* Attack* Purchase " +
		"Chain Applied Choice* StartChoices* Throw* FirstRoll* Activating* " +
		"FreeLoot PendingTriggers* Compost ExtraTurn Boosts* Shields* " +
		"Lingering* Markers* RewardChanges* ManyDice* Stack* StackSeq Turn* " +
		"Priority Waiting* MaxHand WinSouls Over Winners* Challenge* Sets* " +
		"Banned* AskSeq asks! maxHP* cards! events! forcedRolls!",
	"RNG": "State Inc",
	"Player": "ID Character Cents Hand* InPlay* Damage Dead CopyNextLoot " +
		"SkipTurns Activations LootPlayed TimesDamaged ExtraLootPlays Rolls " +
		"LastRoll LastAttack KilledBy Counters*",
	"Counter": "Name Count",
	"Object": "ID Card Zone Role Controller Charged Damage Counters* Eternal " +
		"CopyOf CopyThisTurn CopyTillYourTurn Gained GainedTill MoreGained* " +
		"ATKSet Flipped LostEternal NoAbilities StartingItem Under* " +
		"KilledBy KilledOn Chosen DoubleRewards ShareRewards* NoSpoils " +
		"HitsThisTurn",
	"Zone": "Kind Deck Slot Index",
	"Gain": "Card Till",
	"Slot": "Cards*",
	"Action": "Kind Player Object Amount From Ask* Reroll Penalty Reward Purchase " +
		"Applied*",
	"Asking": "Ability Mode RollResult Effect Key Controller Source Targets* " +
		"Answers* Order* EventPlayer EventAmount EventStack EventObject " +
		"EventSource EventBy",
	"AbilityRef":     "Card Index Lingering",
	"Chosen":         "Kind Player Object StackID Spent",
	"ReplacementRef": "Object Index",
	"AttackState": "LastRoll On Started Target Revealed Deck Rolls StopAfter Helpers* " +
		"Queue* Share*",
	"PurchaseState": "On",
	"Choice": "Purpose Player Rule Replacements* Objects* Previews* Deck Random " +
		"Slots* Targets* Indexes* Abilities* Ask* Question* Subject Look* " +
		"Owner To",
	"ReplacementPreview": "Actions* Unused* Again",
	"Msg":                "Key Args*",
	"Key":                "name",
	"Arg":                "Name Value!",
	"DeckLook":           "Deck Count",
	"Throw":              "StackID Reroll First",
	"FirstRoll":          "Round Players* Rolls*",
	"Activation":         "Player Source Ability Loot Mode Chosen* Window",
	"PendingTrigger": "Ability Source Controller On EventPlayer EventAmount EventStack " +
		"EventObject EventSource EventBy",
	"Boost":        "Stat Player Object Amount Source",
	"Shield":       "Target Amount Cap Source",
	"Target":       "Object Player IsPlayer",
	"Lingering":    "Ability Source Controller",
	"Marker":       "Object Kind From To",
	"RewardChange": "Object Set Kind Amount",
	"ManyDice":     "Chooser Dice Values* StackID",
	"StackItem": "ID Kind Controller Source Card Roll Amount Label* Target Attack " +
		"Ability Mode Targets* RollFor RollResult EventPlayer EventAmount " +
		"EventStack EventObject EventSource EventBy Reward Checked Text " +
		"SourceCard TargetCard Sources* Final",
	"StatSourceView": "Stat Source Card Amount Lasts When*",
	"Turn": "Active Number Step LootPlays Attacks Purchases AttackRolls " +
		"BonusLootUsed BonusAttacksUsed BonusPurchasesUsed MustAttack " +
		"MustAttacks MustAttackDeck RandomAttackMade DeckAttacks ShopBought " +
		"FreeShopItem EndDeclared DeathEnd NextHelpers* NextShare* " +
		"MonsterDied RoomReplaced LateRooms* Entered",
	"Priority": "Open Holder Passes Deferred",
	"Prompt":   "Kind Player Count Options* Purpose Text* Look* Players*",
	"ChallengeState": "ID Difficulty Timer Card Presents* Holidays Lost Reason* Starter " +
		"TurnFrom* RoundFrom* RoomTurn",
	"Present": "Kind Object Destroyed Flipped",
	"hpKey":   "p id",
}

// TestCloneCoversEveryField walks the types of Game: each struct type
// is in cloneRules with exactly its fields, a plain field holds no list,
// map or pointer, and a marked one does.
func TestCloneCoversEveryField(t *testing.T) {
	seen := map[string]bool{}
	var walk func(rt reflect.Type)
	walk = func(rt reflect.Type) {
		for {
			switch rt.Kind() {
			case reflect.Pointer, reflect.Slice, reflect.Array:
				rt = rt.Elem()
				continue
			case reflect.Map:
				walk(rt.Key())
				rt = rt.Elem()
				continue
			default:
			}
			break
		}
		if rt.Kind() != reflect.Struct || seen[rt.Name()] {
			return
		}
		seen[rt.Name()] = true
		rule, ok := cloneRules[rt.Name()]
		if !ok {
			t.Errorf("type %s is not in cloneRules: copy it in Game.Clone and list its fields", rt.Name())
			return
		}
		listed := strings.Fields(rule)
		var fields []string
		for i := range rt.NumField() {
			fields = append(fields, rt.Field(i).Name)
		}
		names := make([]string, len(listed))
		for i, f := range listed {
			names[i] = strings.TrimRight(f, "*!")
		}
		if !slices.Equal(names, fields) {
			t.Errorf("%s: fields %v, cloneRules lists %v: a new field needs a copy rule in Game.Clone and here", rt.Name(), fields, names)
			return
		}
		for i, f := range listed {
			ft := rt.Field(i).Type
			switch {
			case strings.HasSuffix(f, "!"):
			case strings.HasSuffix(f, "*"):
				if !holdsRefs(ft) {
					t.Errorf("%s.%s is a plain value: drop its *", rt.Name(), names[i])
				}
				walk(ft)
			default:
				if holdsRefs(ft) {
					t.Errorf("%s.%s holds a list, map or pointer: copy it in Game.Clone and mark it *", rt.Name(), names[i])
				}
				walk(ft)
			}
		}
	}
	walk(reflect.TypeFor[Game]())
	if t.Failed() {
		return // a type left unwalked is no stale entry
	}
	for name := range cloneRules {
		if !seen[name] {
			t.Errorf("cloneRules lists %s, which Game no longer holds", name)
		}
	}
}

// holdsRefs says whether a value of type rt shares memory when copied.
// A string never changes, so it is a plain value.
func holdsRefs(rt reflect.Type) bool {
	switch rt.Kind() {
	case reflect.Slice, reflect.Map, reflect.Pointer, reflect.Interface, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return true
	case reflect.Array:
		return holdsRefs(rt.Elem())
	case reflect.Struct:
		for i := range rt.NumField() {
			if holdsRefs(rt.Field(i).Type) {
				return true
			}
		}
	default:
	}
	return false
}
