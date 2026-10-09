package rulesengine

// Question is one question an Ask effect puts to the ability's
// controller while it resolves. Options gets the answers so far and
// returns the option labels; no options skips the question (answer -1).
type Question struct {
	Text    string
	Options func(c *Ctx, answers []int) []string
	// Player, if set, is who answers; by default the controller.
	Player func(c *Ctx, answers []int) PlayerID
	// Random: nobody answers; the game picks an option at random
	// ("choose a player at random").
	Random bool
}

type askEffect struct {
	questions []Question
	do        func(c *Ctx, answers []int)
}

// Ask asks the controller questions in order while the ability resolves,
// then calls do with the answers: Ask(do, q1, q2). Targets and "choose
// one-" are not questions: they are picked on activation (R-ABIL-04).
// What do queues happens before the ability's later effects.
func Ask(do func(c *Ctx, answers []int), questions ...Question) Effect {
	return askEffect{questions: questions, do: do}
}

// Asking is an Ask effect waiting for answers. It is plain data, so a
// saved game can continue; the effect itself is found again by ref.
type Asking struct {
	Ability    AbilityRef `json:"ability"`
	Mode       int        `json:"mode,omitempty"`
	RollResult int        `json:"roll_result,omitempty"`
	Effect     int        `json:"effect"` // index in the effect list
	// Key finds the effect while the game runs, also when it is nested
	// in another effect; Effect is the fallback after a Load.
	Key        int      `json:"key"`
	Controller PlayerID `json:"controller"`
	Source     ObjectID `json:"source"`
	Targets    []Chosen `json:"targets,omitempty"`
	Answers    []int    `json:"answers,omitempty"`
	// EventPlayer and EventAmount come from the trigger, if any.
	EventPlayer PlayerID `json:"event_player,omitempty"`
	EventAmount int      `json:"event_amount,omitempty"`
	EventStack  int      `json:"event_stack,omitempty"`
	EventObject ObjectID `json:"event_object,omitempty"`
}

func (e askEffect) apply(c *Ctx) {
	a := &Asking{
		Ability: c.ref, Mode: c.mode, RollResult: c.roll, Effect: c.effect,
		Controller: c.Controller, Source: c.Source, Targets: c.Targets,
		EventPlayer: c.EventPlayer, EventAmount: c.EventAmount, EventStack: c.EventStack, EventObject: c.EventObject,
	}
	c.G.AskSeq++
	a.Key = c.G.AskSeq
	if c.G.asks == nil {
		c.G.asks = map[int]askEffect{}
	}
	c.G.asks[a.Key] = e
	c.G.enqueue(Action{Kind: ActAsk, Player: c.Controller, Ask: a})
}

// Do runs more effects in the same context, e.g. from an Ask's do.
func (c *Ctx) Do(effects ...Effect) {
	for _, e := range effects {
		e.apply(c)
	}
}

func (a *Asking) ctx(g *Game) *Ctx {
	return &Ctx{
		G: g, Controller: a.Controller, Source: a.Source, Targets: a.Targets,
		EventPlayer: a.EventPlayer, EventAmount: a.EventAmount, EventStack: a.EventStack, EventObject: a.EventObject,
		ref: a.Ability, mode: a.Mode, roll: a.RollResult, effect: a.Effect,
	}
}

func (g *Game) askEffectOf(a *Asking) askEffect {
	if e, ok := g.asks[a.Key]; ok {
		return e
	}
	effects := g.effectsOf(a.Ability, a.Mode, a.RollResult)
	e, ok := effects[a.Effect].(askEffect)
	if !ok {
		panic("rulesengine: asking a non-ask effect of " + string(a.Ability.Card))
	}
	return e
}

// continueAsk asks the next question, or finishes the effect.
func (g *Game) continueAsk(a *Asking) {
	e := g.askEffectOf(a)
	c := a.ctx(g)
	for len(a.Answers) < len(e.questions) {
		q := e.questions[len(a.Answers)]
		opts := q.Options(c, a.Answers)
		if len(opts) == 0 {
			a.Answers = append(a.Answers, -1)
			continue
		}
		if q.Random {
			i := g.RNG.Intn(len(opts))
			g.emit(Event{Kind: EvRandomPick, Player: a.Controller, Amount: i, Text: opts[i]})
			a.Answers = append(a.Answers, i)
			continue
		}
		who := a.Controller
		if q.Player != nil {
			who = q.Player(c, a.Answers)
		}
		g.ask(Choice{Purpose: ChooseAnswer, Player: who, Rule: "R-ABIL-05", Ask: a, Question: q.Text}, opts)
		return
	}
	// What the effect queues goes before the rest of the queue.
	rest := g.Queue
	g.Queue = nil
	delete(g.asks, a.Key)
	e.do(c, a.Answers)
	g.Queue = append(g.Queue, rest...)
}

func (g *Game) answerAsk(c Choice, i int) {
	c.Ask.Answers = append(c.Ask.Answers, i)
	g.continueAsk(c.Ask)
}

// DeckQuestion asks for a deck ("a deck"); the options are the decks in
// play. Read the answer with c.Deck.
func DeckQuestion(text string) Question {
	return Question{Text: text, Options: func(c *Ctx, _ []int) []string {
		var out []string
		for _, d := range c.G.DecksInPlay() {
			out = append(out, d.String()+" deck")
		}
		return out
	}}
}

// Deck is the deck an answer to DeckQuestion picked.
func (c *Ctx) Deck(answer int) DeckKind { return c.G.DecksInPlay()[answer] }

// HandQuestion asks the controller for a card in their hand. Read the
// answer with c.HandCard.
func HandQuestion(text string) Question {
	return Question{Text: text, Options: func(c *Ctx, _ []int) []string {
		return c.G.labels(c.G.Players[c.Controller].Hand)
	}}
}

// HandCard is the card an answer to HandQuestion picked.
func (c *Ctx) HandCard(answer int) ObjectID { return c.G.Players[c.Controller].Hand[answer] }
