package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func SteamStar(c *match.Card) {

	c.Name = "Steam Star"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature)
}

func RippleLotusQ(c *match.Card) {

	c.Name = "Ripple Lotus Q"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus, family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Survivor, fx.When(fx.MySurvivorSummoned, fx.MayTapOpCreature))

}
