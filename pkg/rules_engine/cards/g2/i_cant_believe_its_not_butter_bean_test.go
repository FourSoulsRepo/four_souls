package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestICantBelieveItsNotButterBean(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"lil_battery"}},
		{Character: "cain", Items: []engine.CardRef{"i_cant_believe_its_not_butter_bean", "the_bible"}, Deactivated: []engine.CardRef{"the_bible"}},
	}})
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "lil_battery")}}); err != nil {
		t.Fatal(err)
	}
	tb.Choose("the_bible") // Lil Battery targets Cain's item
	tb.Pass(0)
	tb.Activate(1, "i_cant_believe_its_not_butter_bean", 0, "lil_battery")
	if tb.G.Object(tb.Find(1, "the_bible")).Charged {
		t.Error("Lil Battery resolved; it should be cancelled")
	}
}

func TestICantBelieveItsNotButterBeanOnlyYours(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"the_bible"}, Deactivated: []engine.CardRef{"the_bible"}, Hand: []engine.CardRef{"lil_battery"}},
		{Character: "cain", Items: []engine.CardRef{"i_cant_believe_its_not_butter_bean"}},
	}})
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "lil_battery")}}); err != nil {
		t.Fatal(err)
	}
	tb.Choose("the_bible") // Isaac's own item: not Cain's
	tb.Pass(0)
	if rule := tb.Refused(1, "i_cant_believe_its_not_butter_bean", 0); rule == "" {
		t.Error("it must not cancel an ability that targets another player's item")
	}
}
