package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// JackViperShadowofDoom ...
func JackViperShadowofDoom(c *match.Card) {

	c.Name = "Jack Viper, Shadow of Doom"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, func(card *match.Card, ctx *match.Context) {
		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok &&
			event.Card.ID != card.ID &&
			event.Card.Player == card.Player &&
			event.Card.Civ == civ.Darkness {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				fmt.Sprintf("%s: You may return card to hand.", card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool { return event.Card.ID == x.ID },
				false,
			).Map(func(x *match.Card) {
				ctx.InterruptFlow()

				x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's hand by %s", x.Name, x.Player.Username(), card.Name))
			})
		}
	})

}

// WailingShadowBelbetphlo ...
func WailingShadowBelbetphlo(c *match.Card) {

	c.Name = "Wailing Shadow Belbetphlo"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Slayer)

}
