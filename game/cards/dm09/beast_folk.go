package dm09

import (
	"duel-masters/game/civ"
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

			fx.CantBeBlockedByPowerUpTo3000(x, ctx2) //TODO refactor after merging PR #315 from SindaccoFamily
		})
	})
}
