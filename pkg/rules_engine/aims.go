package rulesengine

// Aim is what an allowed play or use asks to target first, so a client
// can light the targets while the card is dragged (15.25). Intent is the
// index of that intent in Allowed.
type Aim struct {
	Intent  int          `json:"intent"`
	Targets []OptionView `json:"targets"`
}

// Aims lists the first target options of the allowed intents. An intent
// that asks no target first (no target, a mode, a copied ability) has no
// aim; one whose target has no options has an aim with no targets.
func (g *Game) Aims(allowed []Intent) []Aim {
	var out []Aim
	for i, in := range allowed {
		spec, ok := g.firstTarget(in)
		if !ok {
			continue
		}
		opts, _ := g.targetOptions(spec, in.Player, in.Objects[0])
		targets := make([]OptionView, 0, len(opts))
		for _, o := range opts {
			targets = append(targets, targetView(o))
		}
		out = append(out, Aim{Intent: i, Targets: targets})
	}
	return out
}

// firstTarget is the target an intent asks for right away, as apply and
// startActivation would ask it.
func (g *Game) firstTarget(in Intent) (TargetSpec, bool) {
	var ref AbilityRef
	switch in.Kind { //nolint:exhaustive // only plays and uses take targets
	case IntentPlayLoot:
		d := g.def(in.Objects[0])
		i := lootIndex(d)
		if d.CopiesTapOnPlay || i < 0 {
			return TargetSpec{}, false
		}
		ref = AbilityRef{Card: d.Ref, Index: i}
	case IntentActivate:
		ref = g.AbilitiesOf(in.Objects[0])[in.Choice]
	default:
		return TargetSpec{}, false
	}
	ab := g.ability(ref)
	if len(ab.Modes) > 0 || len(ab.Targets) == 0 {
		return TargetSpec{}, false
	}
	return ab.Targets[0], true
}
