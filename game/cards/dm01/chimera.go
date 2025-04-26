package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Gigaberos ...
func Gigaberos(c *match.Card) {

	c.Name = "Gigaberos"
	c.Power = 8000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		creatures, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		otherCreatures := make([]*match.Card, 0)
		for _, creature := range creatures {
			if creature.ID != card.ID {
				otherCreatures = append(otherCreatures, creature)
			}
		}
		this := make([]*match.Card, 0)
		this = append(this, card)

		options := make(map[string][]*match.Card)

		options["This creature"] = this
		options["Your other creatures"] = otherCreatures

		ctx.Match.NewMultipartAction(card.Player, options, 1, 2, "Choose 2 of your other creatures in the battle zone that will be destroyed or destroy this creature", false)

		defer ctx.Match.CloseAction(card.Player)

		for {

			action := <-card.Player.Action

			if len(action.Cards) < 1 || len(action.Cards) > 2 {
				ctx.Match.DefaultActionWarning(card.Player)
				continue
			}

			// must be an attempt to destroy this creature
			if len(action.Cards) == 1 {

				if action.Cards[0] != card.ID {
					ctx.Match.DefaultActionWarning(card.Player)
					continue
				}

				ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
				ctx.InterruptFlow()

				break

			}

			if !match.AssertCardsIn(creatures, action.Cards...) {
				ctx.Match.DefaultActionWarning(card.Player)
				continue
			}

			for _, id := range action.Cards {

				creature, err := card.Player.GetCard(id, match.BATTLEZONE)

				if err != nil {
					continue
				}

				ctx.Match.Destroy(creature, card, match.DestroyedByMiscAbility)

			}

			break

		}

	}))

}

// Gigagiele ...
func Gigagiele(c *match.Card) {

	c.Name = "Gigagiele"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Slayer)

}

// Gigargon ...
func Gigargon(c *match.Card) {

	c.Name = "Gigargon"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s: Select up to 2 creatures from your graveyard that will be returned to your hand", card.Name),
			1,
			2,
			true,
			func(x *match.Card) bool {
				return x.HasCondition(cnd.Creature)
			},
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was returned to %s's hand from their graveyard", x.Name, x.Player.Username()))
		})
	}))

}
