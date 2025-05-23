package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AquaHulcus ...
func AquaHulcus(c *match.Card) {

	c.Name = "Aqua Hulcus"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.MayDraw1))

}

// AquaKnight ...
func AquaKnight(c *match.Card) {

	c.Name = "Aqua Knight"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, fx.ReturnToHand))

}

// AquaSniper ...
func AquaSniper(c *match.Card) {

	c.Name = "Aqua Sniper"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		cards := make(map[string][]*match.Card)

		myCards, err := card.Player.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		opponentCards, err := ctx.Match.Opponent(card.Player).Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		if len(myCards) < 1 && len(opponentCards) < 1 {
			return
		}

		cards["Your creatures"] = myCards
		cards["Opponent's creatures"] = opponentCards

		ctx.Match.NewMultipartAction(card.Player, cards, 1, 2, fmt.Sprintf("%s: Choose up to 2 creatures in the battle zone and return them to their owners' hands", card.Name), true)

		for {

			action := <-card.Player.Action

			if action.Cancel {
				break
			}

			if len(action.Cards) < 1 || len(action.Cards) > 2 {
				break
			}

			for _, vid := range action.Cards {

				ref, err := c.Player.MoveCard(vid, match.BATTLEZONE, match.HAND, card.ID)

				if err != nil {

					ref, err := ctx.Match.Opponent(c.Player).MoveCard(vid, match.BATTLEZONE, match.HAND, card.ID)

					if err == nil {
						ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand", ref.Name, ctx.Match.PlayerRef(ref.Player).Socket.User.Username))
					}

				} else {
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand", ref.Name, ctx.Match.PlayerRef(ref.Player).Socket.User.Username))
				}

			}

			break

		}

		ctx.Match.CloseAction(c.Player)

	}))

}

// AquaSoldier ...
func AquaSoldier(c *match.Card) {

	c.Name = "Aqua Soldier"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, fx.ReturnToHand))

}

// AquaVehicle ...
func AquaVehicle(c *match.Card) {

	c.Name = "Aqua Vehicle"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature)

}
