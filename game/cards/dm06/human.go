package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func ArmoredDecimatorValkaizer(c *match.Card) {

	c.Name = "Armored Decimator Valkaizer"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: You may select 1 opponent's creature with 4000 or less power and destroy it", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 4000 },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})
	}))

}

func MigasaAdeptOfChaos(c *match.Card) {

	c.Name = "Migasa, Adept of Chaos"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s activated %s's tap ability", card.Player.Username(), card.Name))
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select 1 fire creature from your battlezone that will gain \"Double Breaker\"", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return x.Civ == civ.Fire },
			false,
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.DoubleBreaker, true, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was given \"Double Breaker\" power by %s until end of turn", x.Name, card.Name))
		})
	}

	c.Use(fx.Creature, fx.TapAbility)

}

func ChoyaTheUnheeding(c *match.Card) {

	c.Name = "Choya, the Unheeding"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker1000, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.Battle); ok {
			if !event.Blocked || event.Attacker != card {
				return
			}

			ctx.InterruptFlow()

			event.Attacker.Tapped = true
			event.Defender.Tapped = true

		}
	})
}
