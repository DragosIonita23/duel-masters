package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// VreemahFreakyMojoTotem ...
func VreemahFreakyMojoTotem(c *match.Card) {

	c.Name = "Vreemah, Freaky Mojo Totem"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature,
		fx.When(fx.AnotherOwnCreatureSummoned, func(card *match.Card, ctx *match.Context) {
			beastFolks := fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasFamily(family.BeastFolk)
				},
			)

			beastFolks = append(beastFolks, fx.FindFilter(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.HasFamily(family.BeastFolk)
				},
			)...)

			beastFolks.Map(func(x *match.Card) {
				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

					if x.Zone != match.BATTLEZONE {
						x.RemoveConditionBySource(card.ID)
						exit()
						return
					}

					if _, ok := ctx2.Event.(*match.EndStep); ok {
						x.RemoveConditionBySource(card.ID)
						exit()
						return
					}

					x.AddUniqueSourceCondition(cnd.PowerAmplifier, 2000, card.ID)
					x.AddUniqueSourceCondition(cnd.DoubleBreaker, true, card.ID)

				})
			})
		}))
}
