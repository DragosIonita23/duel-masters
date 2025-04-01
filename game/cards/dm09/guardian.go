package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GlenaVueleTheHypnotic ...
func GlenaVueleTheHypnotic(c *match.Card) {

	c.Name = "Glena Vuele, the Hypnotic"
	c.Power = 8500
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker,
		fx.When(fx.OpponentUsedShieldTrigger, func(card *match.Card, ctx *match.Context) {
			if fx.BinaryQuestion(
				card.Player,
				ctx.Match,
				fmt.Sprintf("%s's effect: do you want to add the top card of your deck to your shields?", card.Name)) {
				fx.TopCardToShield(card, ctx)
			}
		}))
}
