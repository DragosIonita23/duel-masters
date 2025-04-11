package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SilvermoonTrailblazer ...
func SilvermoonTrailblazer(c *match.Card) {

	c.Name = "Silvermoon Trailblazer"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = silvermoonTrailblazerTapAbility

	c.Use(fx.Creature, fx.TapAbility)

}

func silvermoonTrailblazerTapAbility(card *match.Card, ctx *match.Context) {
	family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. Creatures of that race can't be blocked by creatures that have power 3000 or less this turn.", card.Name))

	fx.FindFilter(
		card.Player,
		match.BATTLEZONE,
		func(x *match.Card) bool {
			return x.HasFamily(family)
		},
	).Map(func(x *match.Card) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if x.Zone != match.BATTLEZONE {
				exit()
				return
			}

			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				exit()
				return
			}

			fx.CantBeBlockedByPowerUpTo3000(x, ctx2)
		})
	})
}

// StormWranglerTheFurious ...
func StormWranglerTheFurious(c *match.Card) {

	c.Name = "Storm Wrangler, the Furious"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution,
		func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.Battle); ok && event.Attacker == card && event.Blocked {
				event.AttackerPower += 3000
			}
		},
		fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: You may choose one of your opponent's untapped creatures that has 'Blocker'. This turn, that creature blocks %s if able and this creature can't be blocked by other creatures.", card.Name, card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool {
					return !x.Tapped && x.HasCondition(cnd.Blocker)
				},
				false,
			).Map(func(x *match.Card) {
				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
					if x.Zone != match.BATTLEZONE || card.Zone != match.BATTLEZONE {
						exit()
						return
					}

					if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
						exit()
						return
					}

					fx.CantBeBlockedByOtherCreaturesBesidesX(card, ctx2, x)

					if event, ok := ctx2.Event.(*match.Block); ok {
						if event.Attacker == card {
							for _, blocker := range event.Blockers {
								if blocker == x {
									// We force the opponent to block with this, i.e.
									// We cancel the Block event normal behaviour
									ctx2.InterruptFlow()

									// And we manually trigger the battle event
									// Between this creature and the selected opp blocker
									ctx2.Match.Battle(card, x, true)
								}
							}
						}
					}
				})
			})
		}))

}
