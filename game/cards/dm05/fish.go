package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SolidskinFish ...
func SolidskinFish(c *match.Card) {

	c.Name = "Solidskin Fish"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.Fish}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.ReturnMyCardFromMZToHand))

}

// SpikestrikeIchthysQ ...
func SpikestrikeIchthysQ(c *match.Card) {

	c.Name = "Spikestrike Ichthys Q"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.Fish, family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.CantBeBlocked, true, card.ID)
			})
		})
	}))

}
