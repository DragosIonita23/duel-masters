package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// CosmicWing ...
func CosmicWing(c *match.Card) {

	c.Name = "Cosmic Wing"
	c.Civ = civ.Light
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose 1 of your creatures in the battlezone. It can't be blocked this turn.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
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

				x.AddUniqueSourceCondition(cnd.CantBeBlocked, true, card.ID)

			})
		})
	}))
}

// NexusCharger ...
func NexusCharger(c *match.Card) {

	c.Name = "Nexus Charger"
	c.Civ = civ.Light
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, fx.HandCardToShield))
}

// AbductionCharger ...
func AbductionCharger(c *match.Card) {

	c.Name = "Abduction Charger"
	c.Civ = civ.Water
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.Charger,
		fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

			cards := make(map[string][]*match.Card, 0)

			myCreatures := fx.Find(card.Player, match.BATTLEZONE)
			oppCreatures := fx.Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

			cards["Your creatures"] = myCreatures
			cards["Your opponent's creatures"] = oppCreatures

			fx.SelectMultipart(
				card.Player,
				ctx.Match,
				cards,
				fmt.Sprintf("%s's effect: Choose up to 2 creatures in the battlezone and return them to their owner's hands.", card.Name),
				1,
				2,
				true,
			).Map(func(x *match.Card) {
				_, err := x.Player.MoveCard(x.ID, match.BATTLEZONE, match.HAND, card.ID)

				if err != nil {
					return
				}

				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was sent to its owner's hand from the battlezone by %s's effect.", x.Name, card.Name))
			})

		}))

}

// GrinningHunger ...
func GrinningHunger(c *match.Card) {

	c.Name = "Grinning Hunger"
	c.Civ = civ.Darkness
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		indexChoice := fx.MultipleChoiceQuestion(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			fmt.Sprintf("%s's effect: Choose one of your creatures in the battlezone or one of your shields and put it into your graveyard.\r\nChoose 'Battlezone' OR 'Shields' to continue.", card.Name),
			[]string{"Battlezone", "Shields"},
		)

		if indexChoice == 0 {
			fx.Select(
				ctx.Match.Opponent(card.Player),
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: Choose one of your creatures in the battlezone and destroy it.", card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
			})
		} else if indexChoice == 1 {
			fx.SelectBackside(
				ctx.Match.Opponent(card.Player),
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.SHIELDZONE,
				fmt.Sprintf("%s's effect: Choose one of your shields and destroy it.", card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.SHIELDZONE, match.GRAVEYARD, card.ID)
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was put into the graveyard from %s's shieldzone.", x.Name, x.Player.Username()))
			})
		}

	}))

}

// UnifiedResistance ...
func UnifiedResistance(c *match.Card) {

	c.Name = "Unified Resistance"
	c.Civ = civ.Light
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		family := fx.ChooseAFamily(
			card,
			ctx,
			fmt.Sprintf("%s's effect: Choose a race. Until the start of your next turn, each of your creatures in the battlezone of that race gets 'Blocker'", card.Name),
		)

		if family != "" {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.HasFamily(family) && !x.HasCondition(card.Name+"-special")
					},
				).Map(func(x *match.Card) {
					ctx2.Match.ApplyPersistentEffect(func(ctx3 *match.Context, exit2 func()) {
						x.AddUniqueSourceCondition(card.Name+"-special", true, card.ID)

						_, ok := ctx3.Event.(*match.StartOfTurnStep)
						if ok && ctx3.Match.IsPlayerTurn(card.Player) {
							x.RemoveConditionBySource(card.ID)
							exit2()
							return
						}

						fx.ForceBlocker(x, ctx3, card.ID)
					})
				})

				_, ok := ctx2.Event.(*match.StartOfTurnStep)
				if ok && ctx2.Match.IsPlayerTurn(card.Player) {
					exit()
					return
				}

			})
		}

	}))

}

// ImpossibleTunnel ...
func ImpossibleTunnel(c *match.Card) {

	c.Name = "Impossible Tunnel"
	c.Civ = civ.Water
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		family := fx.ChooseAFamily(
			card,
			ctx,
			fmt.Sprintf("%s's effect: Choose a race. Creatures of that race can't be blocked this turn.", card.Name),
		)

		if family != "" {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.HasFamily(family) && !x.HasCondition(card.Name+"-special")
					},
				).Map(func(x *match.Card) {
					ctx2.Match.ApplyPersistentEffect(func(ctx3 *match.Context, exit2 func()) {
						x.AddUniqueSourceCondition(card.Name+"-special", true, card.ID)
						x.AddUniqueSourceCondition(cnd.CantBeBlocked, true, card.ID)

						_, ok := ctx3.Event.(*match.EndOfTurnStep)
						if ok && ctx3.Match.IsPlayerTurn(card.Player) {
							x.RemoveConditionBySource(card.ID)
							exit2()
						}
					})
				})

				_, ok := ctx2.Event.(*match.EndOfTurnStep)
				if ok && ctx2.Match.IsPlayerTurn(card.Player) {
					exit()
				}

			})
		}

	}))

}

// SubmarineProject ...
func SubmarineProject(c *match.Card) {

	c.Name = "Submarine Project"
	c.Civ = civ.Water
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.LookTop4Put1IntoHandReorderRestOnBottomDeck)

}

// SlashCharger ...
func SlashCharger(c *match.Card) {

	c.Name = "Slash Charger"
	c.Civ = civ.Darkness
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.Charger, func(card *match.Card, ctx *match.Context) {

		choices := []string{"Yourself", "Your opponent"}

		choiceIndex := fx.MultipleChoiceQuestion(
			card.Player,
			ctx.Match,
			fmt.Sprintf("Choose between yourself and your opponent for applying %s's effect.", card.Name),
			choices,
		)

		var choicePlayer *match.Player
		var choiceMessageFormat string
		var moveMessageFormat string

		if choiceIndex == 0 {
			choicePlayer = card.Player
			choiceMessageFormat = "You may take a card from your deck and put it into your graveyard."
			moveMessageFormat = "his deck."
		} else if choiceIndex == 1 {
			choicePlayer = ctx.Match.Opponent(card.Player)
			choiceMessageFormat = "You may take a card from your opponent's deck and put it into his graveyard."
			moveMessageFormat = "his opponent's deck."
		} else {
			return
		}

		fx.Select(
			card.Player,
			ctx.Match,
			choicePlayer,
			match.DECK,
			fmt.Sprintf("%s's effect: %s", card.Name, choiceMessageFormat),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.DECK, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(choicePlayer, fmt.Sprintf("%s put %s in graveyard from %s", card.Player.Username(), x.Name, moveMessageFormat))
		})

		if choicePlayer == card.Player {
			fx.ShuffleDeck(card, ctx, false)
		} else {
			fx.ShuffleDeck(card, ctx, true)
		}

	})

}

// ZombieCarnival ...
func ZombieCarnival(c *match.Card) {

	c.Name = "Zombie Carnival"
	c.Civ = civ.Darkness
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Spell, fx.ShieldTrigger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. Return up to 3 creatures of that race from your graveyard to your hand.", card.Name))

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s's effect: Return up to 3 creatures of that race from your graveyard to your hand.", card.Name),
			1,
			3,
			true,
			func(x *match.Card) bool {
				return x.HasFamily(family) && x.HasCondition(cnd.Creature)
			},
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was return to your hand from your graveyard by %s's effect.", x.Name, card.Name))
		})

	}))

}

// BlizzardOfSpears ...
func BlizzardOfSpears(c *match.Card) {

	c.Name = "Blizzard of Spears"
	c.Civ = civ.Fire
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		creatures := fx.FindFilter(
			card.Player,
			match.BATTLEZONE,
			func(x *match.Card) bool {
				return ctx.Match.GetPower(x, false) <= 4000
			},
		)

		creatures = append(creatures, fx.FindFilter(
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			func(x *match.Card) bool {
				return ctx.Match.GetPower(x, false) <= 4000
			},
		)...)

		creatures.Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})

	}))

}

// FistsOfForever ...
func FistsOfForever(c *match.Card) {

	c.Name = "Fists of Forever"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose one of your creatures in the battlezone. Whenever that creature wins a battle this turn, untap it.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

				if x.Zone != match.BATTLEZONE {
					exit()
					return
				}

				if event, ok := ctx2.Event.(*match.Battle); ok {
					if (event.Attacker == x && event.AttackerPower > event.DefenderPower) ||
						(event.Blocked && event.Defender == x && event.AttackerPower < event.DefenderPower) {
						ctx2.ScheduleAfter(func() {
							x.Tapped = false
						})
					}
				}

				if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
					exit()
					return
				}

			})
		})

	}))

}

// DanceOfTheSproutlings ...
func DanceOfTheSproutlings(c *match.Card) {

	c.Name = "Dance of the Sproutlings"
	c.Civ = civ.Nature
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. You may put any number of creatures of that race from your hand into your manazone.", card.Name))

		max := len(fx.FindFilter(
			card.Player,
			match.HAND,
			func(x *match.Card) bool {
				return x.HasFamily(family) && x.HasCondition(cnd.Creature)
			},
		))

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: You may put any number of '%s' creatures from your hand into your manazone.", card.Name, family),
			1,
			max,
			true,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.HAND, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was put into %s's manazone from his hand due to %s's effect.", x.Name, card.Player.Username(), card.Name))
		})
	}))

}

// ManaBonanza ...
func ManaBonanza(c *match.Card) {

	c.Name = "Mana Bonanza"
	c.Civ = civ.Nature
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		maxDeck := len(fx.Find(
			card.Player,
			match.DECK,
		))

		manaLen := len(fx.Find(
			card.Player,
			match.MANAZONE,
		))

		if manaLen > maxDeck {
			maxDeck = manaLen
		}

		for _, deckCard := range card.Player.PeekDeck(maxDeck) {
			card.Player.MoveCard(deckCard.ID, match.DECK, match.MANAZONE, card.ID)
			deckCard.Tapped = true
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was put into %s's manazone from the top of this deck due to %s's effect.", deckCard.Name, card.Player.Username(), card.Name))
		}

		ctx.Match.BroadcastState()
	}))

}

// VineCharger ...
func VineCharger(c *match.Card) {

	c.Name = "Vine Charger"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, fx.Charger, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose one of your creatures in the battlezone and put it into your manazone.", card.Name),
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Opponent(card.Player).MoveCard(x.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was put into %s's manazone from his battlezone due to %s's effect.", x.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
		})
	}))

}

// RelentlessBlitz ...
func RelentlessBlitz(c *match.Card) {

	c.Name = "Relentless Blitz"
	c.Civ = civ.Fire
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {
		family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. This turn, each creature of that race can attack untapped creatures and can't be blocked while attacking a creature.", card.Name))

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: This turn, each '%s' creature can attack untapped creatures and can't be blocked while attacking a creature.", card.Name, family),
			1,
			1,
			false,
			func(x *match.Card) bool {
				return x.HasFamily(family) && x.HasCondition(cnd.Creature)
			},
			false,
		).Map(func(x *match.Card) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if x.Zone != match.BATTLEZONE {
					x.RemoveConditionBySource(card.ID)
					exit()
					return
				}

				if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
					x.RemoveConditionBySource(card.ID)
					exit()
					return
				}

				x.AddUniqueSourceCondition(cnd.AttackUntapped, true, card.ID)
				fx.CantBeBlockedWhileAttackingACreature(x, ctx2)
			})
		})
	}))

}
