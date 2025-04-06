package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SimianWarriorGrash ...
func SimianWarriorGrash(c *match.Card) {

	c.Name = "Simian Warrior Grash"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature,
		func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && card.Zone == match.BATTLEZONE {
				if event.Card.HasFamily(family.Armorloid) &&
					event.Card.Player == card.Player {
					fx.Select(
						ctx.Match.Opponent(card.Player),
						ctx.Match,
						ctx.Match.Opponent(card.Player),
						match.MANAZONE,
						fmt.Sprintf("%s's effect: Choose a card from your mana zone and put it into your graveyard.", card.Name),
						1,
						1,
						false,
					).Map(func(x *match.Card) {
						ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
						ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was put into %s's graveyard from his manazone by %s's effect.", x.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
					})
				}
			}
		})

}

// SteamRumblerKain ...
func SteamRumblerKain(c *match.Card) {

	c.Name = "Steam Rumbler Kain"
	c.Power = 5000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature,
		func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.AttackConfirmed); ok && event.CardID == card.ID && card.Zone == match.BATTLEZONE {
				fx.SelectBackside(
					card.Player,
					ctx.Match,
					card.Player,
					match.SHIELDZONE,
					fmt.Sprintf("%s's effect: Choose one of your shields and put it into your graveyard.", card.Name),
					1,
					1,
					false,
				).Map(func(x *match.Card) {
					card.Player.MoveCard(x.ID, match.SHIELDZONE, match.GRAVEYARD, card.ID)
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was put into %s's graveyard from his shields by %s's effect.", x.Name, card.Player.Username(), card.Name))
				})
			}
		})

}
