package main

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Game struct {
	Deck          *Deck
	Dealer        *Player
	Players       []*Player
	numPlayers    int
	playerNames   [7]*widget.Entry
	startingChips int
	state         string
}

type Player struct {
	Name        string
	Chips       int
	Hand        []*Card
	Bet         int
	Score       int
	Chip        *canvas.Image
	PlayersTurn bool
	HiddenCard  bool
	DoubleDown  bool
	Quit        bool
}

func NewGame() *Game {
	game := &Game{}
	game.Deck = buildDeck()
	game.Deck.shuffle()
	game.Dealer = &Player{Name: "Dealer", HiddenCard: true}
	game.Players = []*Player{}
	game.state = "Get Players"
	game.startingChips = 100

	return game
}

func setNextPlayer(t *tableRender, i int) {
	t.game.Players[i].PlayersTurn = false
	t.playerActions[i].Hide()
	if i == t.game.numPlayers-1 {
		t.game.Players[0].PlayersTurn = true
		t.game.state = "Dealer Actions"
		t.playerActions[i].Hide()
	} else if t.game.Players[i+1].Score == 21 {
		if i+2 == t.game.numPlayers {
			t.game.Players[0].PlayersTurn = true
			t.game.state = "Dealer Actions"
			t.playerActions[i].Hide()
		} else {
			t.game.Players[i+2].PlayersTurn = true
			t.playerActions[i+2].Show()
		}
	} else {
		t.game.Players[i+1].PlayersTurn = true
		t.playerActions[i+1].Show()
	}
}

func action(t *tableRender, s string) {
	for i := 0; i < t.game.numPlayers; i++ {
		if t.game.Players[i].PlayersTurn {
			switch s {
			case "Hit":
				card := t.game.Deck.DealCard()
				t.game.Players[i].Hand = append(t.game.Players[i].Hand, card)

				t.game.updateScores(t)
				if t.game.Players[i].Score >= 21 {
					setNextPlayer(t, i)
				}
				t.Refresh()

			case "Stand":
				setNextPlayer(t, i)
				t.Refresh()

			case "Split":

			case "Double Down":
				t.game.Players[i].Chips = t.game.Players[i].Chips - t.game.Players[i].Bet
				t.game.Players[i].Bet = t.game.Players[i].Bet * 2
				card := t.game.Deck.DealCard()
				t.game.Players[i].Hand = append(t.game.Players[i].Hand, card)

				setNextPlayer(t, i)
				t.Refresh()
			}

			break
		}
	}
}

func (g *Game) dealHands(t *tableRender) {
	for i := 0; i < 2; i++ {
		for _, p := range g.Players {
			p.Hand = append(p.Hand, g.Deck.DealCard())
		}
		g.Dealer.Hand = append(g.Dealer.Hand, g.Deck.DealCard())
	}
	g.updateScores(t)
	g.state = "Player Actions"
	if g.Players[0].Score != 21 {
		t.playerActions[0].Show()
	} else {
		t.playerActions[1].Show()
	}
	t.Refresh()
}

func (g *Game) initPlayers(num int) {
	for i := 0; i < num; i++ {
		g.Players = append(g.Players, &Player{Chips: g.startingChips})
	}
	g.Players[0].PlayersTurn = true
}

func (g *Game) getNumberOfPlayers(t *tableRender) {

	b1 := widget.NewButton("1", func() { g.numPlayers = 1; g.initPlayers(1); g.state = "Get Names"; t.Refresh() })
	b2 := widget.NewButton("2", func() { g.numPlayers = 2; g.initPlayers(2); g.state = "Get Names"; t.Refresh() })
	b3 := widget.NewButton("3", func() { g.numPlayers = 3; g.initPlayers(3); g.state = "Get Names"; t.Refresh() })
	b4 := widget.NewButton("4", func() { g.numPlayers = 4; g.initPlayers(4); g.state = "Get Names"; t.Refresh() })
	b5 := widget.NewButton("5", func() { g.numPlayers = 5; g.initPlayers(5); g.state = "Get Names"; t.Refresh() })
	b6 := widget.NewButton("6", func() { g.numPlayers = 6; g.initPlayers(6); g.state = "Get Names"; t.Refresh() })
	b7 := widget.NewButton("7", func() { g.numPlayers = 7; g.initPlayers(7); g.state = "Get Names"; t.Refresh() })

	msg := widget.NewLabel("Welcome to the game of Blackjack!\nChose how many players:")
	msg.Alignment = fyne.TextAlignCenter
	buttons := container.NewHBox(b1, b2, b3, b4, b5, b6, b7)

	c := container.NewVBox(
		msg,
		container.NewCenter(buttons),
	)

	bg := canvas.NewRectangle(color.NRGBA{R: 1, G: 50, B: 32, A: 200})

	t.numPlayers.Add(bg)
	t.numPlayers.Add(c)
	t.numPlayers.Layout = layout.NewMaxLayout()
	t.numPlayers.Show()
}

func (g *Game) getPlayerNames(t *tableRender) {
	count := 0
	names := [7]*widget.Entry{}
	b := widget.NewButton("Submit", func() { g.state = "Bet"; t.Refresh() })
	b.Resize(fyne.NewSize(50, 150))
	b.Disable()
	for i := 0; i < 7; i++ {
		if i < g.numPlayers {
			e := widget.NewEntry()
			e.SetPlaceHolder(fmt.Sprintf("Player %v name...", i+1))
			e.Validator = func(s string) error {
				if s == "" {
					return errors.New("Required")
				}
				return nil
			}
			e.SetOnValidationChanged(func(err error) {
				count++
				if count/2 == g.numPlayers {
					b.Enable()
				}
			})
			names[i] = e

		} else {
			e := widget.NewEntry()
			e.Hide()
			names[i] = e
		}
	}
	msg := widget.NewLabel("Enter player names:")
	msg.Alignment = fyne.TextAlignCenter
	c := container.NewHBox(
		layout.NewSpacer(),
		container.NewVBox(
			msg,
			names[0],
			names[1],
			names[2],
			names[3],
			names[4],
			names[5],
			names[6],
			b,
		),
		layout.NewSpacer(),
	)
	g.playerNames = names
	bg := canvas.NewRectangle(color.NRGBA{R: 1, G: 50, B: 32, A: 240})

	t.playerNames.Add(bg)
	t.playerNames.Add(c)
	t.playerNames.Layout = layout.NewMaxLayout()
	t.playerNames.Show()
	t.numPlayers.Hide()
}

func (g *Game) setPlayerNames() {
	for i := 0; i < g.numPlayers; i++ {
		g.Players[i].Name = g.playerNames[i].Text
	}
}

func (g *Game) playerBets(t *tableRender) {
	bet := 0
	for i := 0; i < g.numPlayers; i++ {
		if g.Players[i].PlayersTurn == true {
			chips := g.Players[i].Chips
			b := widget.NewButton("Bet", func() {
				g.Players[i].PlayersTurn = false
				t.playerBets.Objects = []fyne.CanvasObject{}
				if i == g.numPlayers-1 {
					g.Players[0].PlayersTurn = true
					t.playerBets.Hide()
					g.state = "Deal Hands"
				} else {
					g.Players[i+1].PlayersTurn = true
				}
				g.Players[i].Bet = bet
				g.Players[i].Chips = chips - bet

				t.Refresh()
			})
			b.Disable()

			e := widget.NewEntry()
			e.SetPlaceHolder("Place your bet...")
			e.Validator = func(s string) error {
				v, _ := strconv.Atoi(s)
				r := regexp.MustCompile(`^\d+$`)
				if !r.MatchString(s) || v == 0 || v > g.Players[i].Chips || v%2 != 0 {
					b.Disable()
					return errors.New("Required")
				} else {
					b.Enable()
				}
				return nil
			}

			e.OnChanged = func(string) {
				bet, _ = strconv.Atoi(e.Text)
			}
			s := fmt.Sprintf("%v\nYou have %v chips. Place your bet:", g.Players[i].Name, g.Players[i].Chips)
			msg := widget.NewLabel(s)
			msg.Alignment = fyne.TextAlignCenter
			c := container.NewHBox(
				layout.NewSpacer(),
				container.NewVBox(
					msg,
					e,
					b,
				),
				layout.NewSpacer(),
			)
			bg := canvas.NewRectangle(color.NRGBA{R: 1, G: 50, B: 32, A: 240})
			t.playerBets.Add(bg)
			t.playerBets.Add(c)
			t.playerBets.Layout = layout.NewMaxLayout()
			t.playerBets.Show()
			t.playerNames.Hide()
			t.resetTable.Hide()
			break
		}
	}
}

func (g *Game) playerActions(t *tableRender) {
	for i := 0; i < g.numPlayers; i++ {
		b1 := widget.NewButton("Stand", func() { action(t, "Stand") })
		b2 := widget.NewButton("Hit", func() { action(t, "Hit") })
		b3 := widget.NewButton("Double Down", func() { action(t, "Double Down") })
		b4 := widget.NewButton("Split", func() { action(t, "Split") })
		b4.Disable()

		if len(g.Players[i].Hand) == 2 {
			if g.Players[i].Hand[0].Rank == g.Players[i].Hand[1].Rank {
				b4.Enable()
			}
		} else {
			b3.Disable()
		}

		msg := canvas.NewText(fmt.Sprintf("%v's turn", g.Players[i].Name), color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		msg.Alignment = fyne.TextAlignCenter
		buttons := container.NewVBox(b1, b2, b3, b4)

		c := container.NewHBox(
			layout.NewSpacer(),
			container.NewVBox(
				msg,
				container.NewCenter(buttons),
			),
			layout.NewSpacer(),
		)
		bg := canvas.NewRectangle(color.NRGBA{R: 1, G: 50, B: 32, A: 240})
		t.playerActions[i].Add(bg)
		t.playerActions[i].Add(c)
		t.playerActions[i].Layout = layout.NewMaxLayout()
	}
}

func (g *Game) dealerActions(t *tableRender) {
	g.Dealer.HiddenCard = false
	g.updateScores(t)

	if g.Dealer.Score < 17 {
		card := g.Deck.DealCard()
		g.Dealer.Hand = append(g.Dealer.Hand, card)
	} else {
		g.state = "Settle Bets"
	}

}

func (g *Game) playerScores(c [7]*fyne.Container) {
	for i := 0; i < g.numPlayers; i++ {
		if g.numPlayers > 0 {
			bg := canvas.NewCircle(color.NRGBA{R: 0, G: 0, B: 0, A: 255})
			t := canvas.NewText(fmt.Sprint(g.Players[i].Score), color.NRGBA{R: 255, G: 255, B: 255, A: 255})
			t.Alignment = fyne.TextAlignCenter
			t.TextSize = textSize

			// temp solution
			z := canvas.NewText("00", color.NRGBA{R: 0, G: 0, B: 0, A: 255})
			z.Alignment = fyne.TextAlignCenter
			z.TextSize = textSize
			nb := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), z)

			ml := container.New(layout.NewMaxLayout(), bg, nb, t)
			c[i].Add(ml)
		}
	}
}

func (g *Game) settleBets(t *tableRender) {
	for i := 0; i < g.numPlayers; i++ {
		var payout int
		if g.Players[i].Score > 21 {
			// Player Busts
			continue
		} else if g.Players[i].Score == g.Dealer.Score { // Draw
			payout = g.Players[i].Bet
		} else if g.Players[i].Score == 21 { // Blackjack
			payout = int(float32(g.Players[i].Bet) * 2.5)
		} else if g.Dealer.Score > 21 || g.Players[i].Score > g.Dealer.Score { // Dealer Busts
			payout = g.Players[i].Bet * 2
		}

		if g.Players[i].DoubleDown {
			g.Players[i].Chips = g.Players[i].Chips + payout*2
		} else {
			g.Players[i].Chips = g.Players[i].Chips + payout
		}
	}
	g.state = "Restart"
	t.resetTable.Show()
	t.Refresh()
}

func (g *Game) resetTable(t *tableRender) {
	b := widget.NewButton("Let's Go!", func() {

		// Reset players
		for i := 0; i < g.numPlayers; i++ {
			g.Players[i].Hand = []*Card{}
			g.Players[i].Bet = 0
			t.scores[i].Text = ""
			t.scoresBG[i] = &canvas.Circle{}
			t.playerOutcomes[i].Resource = nil
			for j, c := range t.stacks[i].cards {
				if j == 0 {
					c.Resource = imgCardLocationSvg
				} else {
					c.Resource = nil
				}
			}
		}
		g.Players[0].PlayersTurn = true

		// Reset Dealer
		t.game.Dealer.Hand = []*Card{}
		t.game.Dealer.HiddenCard = true
		t.dealer.score.Text = ""
		for _, c := range t.dealer.stack.cards {
			c.Resource = nil
		}

		// Build and new deck and shuffle
		t.game.Deck = buildDeck()
		t.game.Deck.shuffle()

		g.state = "Bet"
		t.Refresh()
	})

	s := "Play Again?\nUncheck name if leaving the table"
	msg := widget.NewLabel(s)
	msg.Alignment = fyne.TextAlignCenter
	cb := [7]*widget.Check{}

	for i := 0; i < 7; i++ {
		if i < g.numPlayers {
			n := fmt.Sprint(g.Players[i].Name)
			ch := widget.NewCheck(n, func(b bool) {})
			ch.Checked = true
			cb[i] = ch
		} else {
			ch := widget.NewCheck("", func(b bool) {})
			ch.Hide()
			cb[i] = ch
		}
	}
	c := container.NewHBox(
		layout.NewSpacer(),
		container.NewVBox(
			msg,
			container.NewCenter(
				container.NewGridWithColumns(2,
					cb[0],
					cb[1],
					cb[2],
					cb[3],
					cb[4],
					cb[5],
					cb[6],
				),
			),
			b,
		),
		layout.NewSpacer(),
	)
	bg := canvas.NewRectangle(color.NRGBA{R: 1, G: 50, B: 32, A: 240})
	t.resetTable.Add(bg)
	t.resetTable.Add(c)
	t.resetTable.Layout = layout.NewMaxLayout()
	t.resetTable.Show()
}

func (g *Game) tallyHand(p *Player) int {
	var ace bool
	s := 0
	for _, h := range p.Hand {
		if h.Value == 11 {
			ace = true
		}
		s += h.Value
	}
	if s > 21 && ace {
		s = s - 10
	}
	return s
}

func (g *Game) updateScores(t *tableRender) {
	// tally dealer's hands
	if g.Dealer.HiddenCard {
		g.Dealer.Score = g.Dealer.Hand[0].Value
	} else {
		g.Dealer.Score = g.tallyHand(g.Dealer)
	}

	// tally player's hands
	for i, p := range g.Players {
		p.Score = g.tallyHand(p)
		if p.Score > 21 {
			t.playerOutcomes[i].Resource = imgBustSvg
			t.playerOutcomes[i].Show()
		} else if g.state == "Settle Bets" {
			if p.Score == g.Dealer.Score {
				t.playerOutcomes[i].Resource = imgDrawSvg
				t.playerOutcomes[i].Show()
			} else if p.Score > g.Dealer.Score && p.Score != 21 {
				t.playerOutcomes[i].Resource = imgWinSvg
				t.playerOutcomes[i].Show()
			} else if g.Dealer.Score > 21 {
				t.playerOutcomes[i].Resource = imgWinSvg
				t.playerOutcomes[i].Show()
			} else if p.Score < g.Dealer.Score {
				t.playerOutcomes[i].Resource = imgLoseSvg
				t.playerOutcomes[i].Show()
			}
		} else if p.Score == 21 {
			t.playerOutcomes[i].Resource = imgBlackjackSvg
			t.playerOutcomes[i].Show()
		}
	}

}
