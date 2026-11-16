package enginetest

import (
	"fmt"
	"reflect"
	"strings"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// CloneRules lists every field of Game and of each type inside it, in
// order, with how Game.Clone copies it (ADR 030, Decision 3.3):
//
//	Name   a plain value, copied with its struct; it holds no list, map
//	       or pointer
//	Name*  Clone gives the copy its own: a line in clone.go
//	Name!  shared or left out on purpose; Clone's comment says why
//
// A new field fails TestCloneCoversEveryField (clone_test.go) until it
// is here; a field marked * that Clone does not copy fails
// TestCloneCopiesEveryField.
var CloneRules = map[string]string{
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

// CloneMark is how Game.Clone copies field name of struct type rt:
// "" a plain value, "*" its own copy, "!" shared or left out. ok is
// false when CloneRules does not list the field.
func CloneMark(rt reflect.Type, name string) (mark string, ok bool) {
	mark, ok = cloneMarks[rt.Name()][name]
	return mark, ok
}

// cloneMarks is CloneRules read once: type, then field, then mark.
var cloneMarks = func() map[string]map[string]string {
	out := map[string]map[string]string{}
	for typ, fields := range CloneRules {
		out[typ] = map[string]string{}
		for _, f := range strings.Fields(fields) {
			name := strings.TrimRight(f, "*!")
			out[typ][name] = strings.TrimPrefix(f, name)
		}
	}
	return out
}()

// NoSharing walks a game and its copy side by side: no slice, map or
// pointer of the copy may be the game's, and a nil one stays nil.
// Fields marked "!" in CloneRules are skipped.
func NoSharing(g, c *engine.Game) error {
	return noSharing(reflect.ValueOf(g).Elem(), reflect.ValueOf(c).Elem(), "Game")
}

func noSharing(a, b reflect.Value, path string) error {
	switch a.Kind() {
	case reflect.Slice, reflect.Map, reflect.Pointer, reflect.Interface:
		if a.IsNil() != b.IsNil() {
			return fmt.Errorf("%s: nil %v, copy nil %v", path, a.IsNil(), b.IsNil())
		}
		if a.IsNil() {
			return nil
		}
		return sameShape(a, b, path)
	case reflect.Array:
		for i := range a.Len() {
			if err := noSharing(a.Index(i), b.Index(i), fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	case reflect.Struct:
		for i := range a.NumField() {
			name := a.Type().Field(i).Name
			mark, ok := CloneMark(a.Type(), name)
			if !ok {
				return fmt.Errorf("%s.%s: not in CloneRules", path, name)
			}
			if mark == "!" {
				continue
			}
			if err := noSharing(a.Field(i), b.Field(i), path+"."+name); err != nil {
				return err
			}
		}
	default: // a plain value
	}
	return nil
}

// sameShape checks a non-nil slice, map, pointer or interface value
// and what it holds.
func sameShape(a, b reflect.Value, path string) error {
	if a.Kind() != reflect.Interface && (a.Kind() != reflect.Slice || a.Cap() > 0) && a.Pointer() == b.Pointer() {
		return fmt.Errorf("%s: the copy shares it", path)
	}
	switch a.Kind() {
	case reflect.Pointer, reflect.Interface:
		return noSharing(a.Elem(), b.Elem(), path)
	case reflect.Slice:
		if a.Len() != b.Len() {
			return fmt.Errorf("%s: %d elements, copy %d", path, a.Len(), b.Len())
		}
		for i := range a.Len() {
			if err := noSharing(a.Index(i), b.Index(i), fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		iter := a.MapRange()
		for iter.Next() {
			bv := b.MapIndex(iter.Key())
			if !bv.IsValid() {
				return fmt.Errorf("%s: the copy lacks key %v", path, iter.Key())
			}
			if err := noSharing(iter.Value(), bv, fmt.Sprintf("%s[%v]", path, iter.Key())); err != nil {
				return err
			}
		}
	default:
	}
	return nil
}
