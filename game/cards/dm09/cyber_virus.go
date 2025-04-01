package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MarchingMotherboard ...
func MarchingMotherboard(c *match.Card) {

	c.Name = "Marching Motherboard"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature,
		fx.When(fx.AnotherOwnCyberSummoned, func(card *match.Card, ctx *match.Context) {
			fx.MayDraw1(card, ctx)
		}))
}
