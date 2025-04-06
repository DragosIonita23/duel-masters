package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// BatDoctorShadowOfUndeath ...
func BatDoctorShadowOfUndeath(c *match.Card) {

	c.Name = "Bat Doctor, Shadow of Undeath"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s's effect: You may return another creature from your graveyard to your hand.", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool {
				return x.ID != card.ID
			},
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to %s's hand from its graveyard by %s's effect.", x.Name, card.Player.Username(), card.Name))
		})

	}))

}

// IceVaporShadowOfAnguish ...
func IceVaporShadowOfAnguish(c *match.Card) {

	c.Name = "Ice Vapor, Shadow of Anguish"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.OppSpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.HAND,
			fmt.Sprintf("%s's effect: Choose a card from your hand and discard it.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was discarded from %s's hand by %s's effect.", x.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
		})

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
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was put into the graveyard from %s's mana zone by %s's effect.", x.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
		})

	}))

}
