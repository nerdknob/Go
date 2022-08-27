package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

const minCardWidth = float32(100)
const minPadding = float32(8)
const cardRatio = 150 / minCardWidth

var (
	cardSize = fyne.Size{Width: minCardWidth, Height: minCardWidth * cardRatio}

	smallPad = minPadding
	overlap  = smallPad * 4
	bigPad   = smallPad + overlap
	textSize = smallPad * 2.3

	minWidth  = cardSize.Width*15 + smallPad*8
	minHeight = cardSize.Height*6 + bigPad + smallPad*2
)

type tableRender struct {
	table   *Table
	game    *Game
	objects []fyne.CanvasObject
	deck    []*Card
	stacks  []*stackRender
	dealer
	playerInput
	playerDetails
}

type dealer struct {
	name    *canvas.Text
	score   *canvas.Text
	scoreBG *canvas.Circle
	stack   *stackRender
}
type stackRender struct {
	cards [9]*canvas.Image
}

type playerInput struct {
	numPlayers    *fyne.Container
	playerNames   *fyne.Container
	playerActions [7]*fyne.Container
	playerBets    *fyne.Container
	resetTable    *fyne.Container
}
type playerDetails struct {
	names          [7]*canvas.Text
	scores         [7]*canvas.Text
	scoresBG       [7]*canvas.Circle
	scoreTest      [7]*fyne.Container
	bets           [7]*canvas.Text
	chipImage      [7]*canvas.Image
	playersChips   [7]*canvas.Text
	playerOutcomes [7]*canvas.Image
}

func (t *tableRender) MinSize() fyne.Size {
	return fyne.NewSize(minWidth, minHeight)
}

func (t *tableRender) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *tableRender) Destroy() {
}

func updateSizes(pad float32) {
	smallPad = pad
	overlap = smallPad * 4
	bigPad = smallPad + overlap
	textSize = smallPad * 2.3
}

func updateCardPosition(c *canvas.Image, x, y float32, s fyne.Size) {
	c.Resize(s)
	c.Move(fyne.NewPos(x, y))
}

func (t *tableRender) updateNumPlayersPosition(s fyne.Size) {
	t.numPlayers.Resize(fyne.NewSize(300, 125))
	t.numPlayers.Move(fyne.NewPos(s.Width/2-150, s.Height/2-62.5))
}

func (t *tableRender) updatePlayerNamesPosition(s fyne.Size) {
	if t.game.numPlayers > 0 {
		w := []float32{135, 175, 215, 255, 295, 335, 375}
		i := t.game.numPlayers - 1
		t.playerNames.Resize(fyne.NewSize(200, w[i]))
		t.playerNames.Move(fyne.NewPos(s.Width/2-100, s.Height/2-w[i]/2))
	}
}

func (t *tableRender) updateBetPosition() {
	f := []float32{0.06, 2.06, 4, 6, 8, 10, 12}
	if t.game.numPlayers > 0 {
		for i := 0; i < t.game.numPlayers; i++ {
			if t.game.Players[i].PlayersTurn == true {
				pos := fyne.NewPos(cardSize.Width*f[i], cardSize.Height*4)
				t.playerInput.playerBets.Move(pos)
				t.playerInput.playerBets.Resize(fyne.NewSize(300, 150))
				break
			}
		}
	}
}

func (t *tableRender) updateNameAndChipsPosition() {
	var m float32 = 1
	for i := 0; i < len(t.names); i++ {
		t.names[i].Move(fyne.NewPos(cardSize.Width*m, cardSize.Height*3.5))
		t.names[i].TextSize = textSize
		t.playersChips[i].Move(fyne.NewPos(cardSize.Width*m, cardSize.Height*3.7))
		t.playersChips[i].TextSize = textSize

		m = m + 2
	}
	m = 6
	t.dealer.name.Move(fyne.NewPos(cardSize.Width*m, cardSize.Height*0.1))
	t.dealer.name.TextSize = textSize
}

func (t *tableRender) updateChipsAndBetsPosition() {
	var m float32 = 1.25
	for i := 0; i < t.game.numPlayers; i++ {
		t.chipImage[i].Move(fyne.NewPos(cardSize.Width*m, cardSize.Height*3))
		t.chipImage[i].Resize(fyne.NewSize(bigPad+smallPad, bigPad+smallPad))
		t.chipImage[i].Show()
		t.bets[i].Move(fyne.NewPos(cardSize.Width*m+smallPad*3, cardSize.Height*3.08))
		t.bets[i].TextSize = textSize
		t.bets[i].Alignment = fyne.TextAlignCenter
		t.bets[i].Show()

		m = m + 2
	}
}

func (t *tableRender) updateScorePosition() {
	var m float32 = 1.8
	for i := 0; i < len(t.names); i++ {
		if t.game.state == "Bet" {
			t.scoreTest[i].Hide()
		} else {
			t.scoreTest[i].Move(fyne.NewPos(cardSize.Width*m, cardSize.Height*3.89))
			t.scoreTest[i].Show()
		}
		m = m + 2
	}
	m = 6
	t.dealer.score.Move(fyne.NewPos(cardSize.Width*m, cardSize.Height*0.3))
	t.dealer.score.TextSize = textSize
}

func (t *tableRender) updateStacksPosition() {
	// set layout of player's cards
	var m float32 = 1
	for i := 0; i < 7; i++ {
		pos := fyne.NewPos(cardSize.Width*m, cardSize.Height*4)
		top := pos.Y
		side := pos.X
		for j := range t.stacks[i].cards {
			updateCardPosition(t.stacks[i].cards[j], side, top, cardSize)

			top += overlap
			side += smallPad
		}
		m = m + 2
	}
	// set layout of dealer's cards
	m = 7
	pos := fyne.NewPos(cardSize.Width*m, cardSize.Height*0.1)
	side := pos.X
	for i := range t.dealer.stack.cards {
		updateCardPosition(t.dealer.stack.cards[i], side, pos.Y, cardSize)

		side += overlap
	}
}

func (t *tableRender) updatePlayerActionsPosition() {
	f := []float32{0.8, 2.8, 4.8, 6.8, 8.8, 10.8, 12.8}
	if t.game.numPlayers > 0 {
		for i := 0; i < t.game.numPlayers; i++ {
			pos := fyne.NewPos(cardSize.Width*f[i], cardSize.Height*1.6)
			t.playerInput.playerActions[i].Move(pos)
			t.playerInput.playerActions[i].Resize(fyne.NewSize(150, 200))
		}
	}
}

func (t *tableRender) updatePlayerOutcomesPosition() {
	var m float32 = 1
	for i := 0; i < t.game.numPlayers; i++ {
		pos := fyne.NewPos(cardSize.Width*m, cardSize.Height*4.05)
		t.playerOutcomes[i].Move(pos)
		t.playerOutcomes[i].Resize(fyne.NewSize(cardSize.Width, cardSize.Height))
		m = m + 2
	}
}

func (t *tableRender) updatePlayAgainPromptPosition(s fyne.Size) {
	if t.game.numPlayers > 0 {
		w := []float32{155, 155, 195, 195, 235, 235, 275}
		i := t.game.numPlayers - 1
		t.resetTable.Resize(fyne.NewSize(300, w[i]))
		t.resetTable.Move(fyne.NewPos(s.Width/2-150, cardSize.Height*1.2))
	}
}

// Layout() defines where the canvas objects in the widget renderer objects slice are positioned on screen
func (t *tableRender) Layout(size fyne.Size) {
	padding := size.Width * .006
	updateSizes(padding)

	newWidth := size.Width / 15.0
	cardSize = fyne.NewSize(newWidth, newWidth*cardRatio)

	t.table.tabletop.Resize(fyne.NewSize(size.Width, cardSize.Height*4))
	t.table.tabletop.Move(fyne.NewPos(0, cardSize.Height*.2))

	t.updateNumPlayersPosition(size)
	t.updatePlayerNamesPosition(size)
	t.updateBetPosition()
	t.updateNameAndChipsPosition()
	t.updateChipsAndBetsPosition()
	t.updateScorePosition()
	t.updateStacksPosition()
	t.updatePlayerActionsPosition()
	t.updatePlayerOutcomesPosition()
	t.updatePlayAgainPromptPosition(size)
}

func (t *tableRender) refreshNameAndChips() {
	if t.game.numPlayers > 0 && len(t.game.Players) > 0 {
		for i := 0; i < len(t.game.Players); i++ {
			t.names[i].Text = t.game.Players[i].Name
			if t.game.state != "Get Names" {
				t.playersChips[i].Text = fmt.Sprintf("Chips: %v", t.game.Players[i].Chips)
			}
		}
		t.dealer.name.Text = t.game.Dealer.Name
	}
}

func (t *tableRender) refreshScores() {
	t.game.updateScores(t)
	t.dealer.score.Text = fmt.Sprintf("%v", t.game.Dealer.Score)
	t.game.playerScores(t.scoreTest)
	// }
}

func (t *tableRender) refreshChipsAndBets() {
	for i := 0; i < t.game.numPlayers; i++ {
		if t.game.Players[i].Bet > 0 {
			t.chipImage[i].Resource = imgChipPng
			t.chipImage[i].Show()
			t.bets[i].Text = fmt.Sprint(t.game.Players[i].Bet)
			t.bets[i].Show()
		} else {
			t.chipImage[i].Resource = nil
			t.bets[i].Text = ""
		}
	}
}

func (t *tableRender) refreshStacks() {
	// adds player's card face images to the widget resources
	if t.game.numPlayers > 0 {
		for i := 0; i < t.game.numPlayers; i++ {
			for j, card := range t.game.Players[i].Hand {
				t.stacks[i].cards[j].Resource = cardFace(card.Rank, card.Suit)
			}

		}
	} else {
		for i := 0; i < 7; i++ {
			t.stacks[i].cards[0].Resource = imgCardLocationSvg
			t.stacks[i].cards[0].Show()
		}
	}

	// adds dealers's card face images to the widget resources
	if len(t.game.Dealer.Hand) > 0 {
		if t.game.Dealer.HiddenCard {
			card := t.game.Dealer.Hand[0]
			t.dealer.stack.cards[0].Resource = cardFace(card.Rank, card.Suit)
			t.dealer.stack.cards[0].Show()
			t.dealer.stack.cards[1].Resource = imgBackSvg
			t.dealer.stack.cards[1].Show()
		} else {
			for i, card := range t.game.Dealer.Hand {
				t.dealer.stack.cards[i].Resource = cardFace(card.Rank, card.Suit)
			}
		}
	}
}

// Refresh() is called anytime there is a state change and the widget needs to be rerendered
func (t *tableRender) Refresh() {

	switch t.game.state {
	case "Get Players":
		t.game.getNumberOfPlayers(t)

	case "Get Names":
		t.game.getPlayerNames(t)

	case "Bet":
		t.game.setPlayerNames()
		t.game.playerBets(t)

	case "Deal Hands":
		t.game.dealHands(t)

	case "Player Actions":
		t.refreshScores()
		t.game.playerActions(t)

	case "Dealer Actions":
		t.game.dealerActions(t)
		t.refreshScores()
		t.Refresh()
	case "Settle Bets":
		t.game.settleBets(t)
		t.Refresh()
	case "Restart":
		t.game.resetTable(t)
		canvas.Refresh(t.table)
	}

	t.refreshChipsAndBets()
	t.refreshStacks()
	t.refreshNameAndChips()
	t.table.tabletop.Resource = imgBlackjackTablePng
	t.table.tabletop.Show()

	canvas.Refresh(t.table)
}

func (t *tableRender) newNameAndScore() {
	for i := 0; i < 7; i++ {
		t.names[i] = &canvas.Text{}
		t.names[i].Text = ""
		t.names[i].TextSize = textSize
		t.scores[i] = &canvas.Text{}
		t.scores[i].Text = ""
		t.scores[i].TextSize = textSize
		t.scoresBG[i] = &canvas.Circle{}
		t.scoreTest[i] = &fyne.Container{}
	}
	t.dealer.name = &canvas.Text{}
	t.dealer.name.Text = ""
	t.dealer.name.TextSize = textSize
	t.dealer.score = &canvas.Text{}
	t.dealer.score.Text = ""
	t.dealer.score.TextSize = textSize
}

func (t *tableRender) appendNameAndScore() {
	for i := 0; i < len(t.names); i++ {
		t.objects = append(t.objects, t.names[i], t.scoreTest[i], t.scoresBG[i], t.scores[i])
	}
	t.objects = append(t.objects, t.dealer.name, t.dealer.score)
}

// instantiate player and dealer card resources
func (t *tableRender) newStacks() {
	// instantiate player's cards
	for i := 0; i < 7; i++ {
		t.stacks = append(t.stacks, &stackRender{})

		for j := 0; j < 9; j++ {
			t.stacks[i].cards[j] = &canvas.Image{}
		}
	}
	// instantiate dealer's cards
	t.dealer.stack = &stackRender{}
	for i := 0; i < len(t.dealer.stack.cards); i++ {
		t.dealer.stack.cards[i] = &canvas.Image{}
	}
}

// appends card resources to the widget renderer objects slice
func (t *tableRender) appendStack() {
	// add player's card resources to the widget renderer objects slice
	for _, s := range t.stacks {
		for _, card := range s.cards {
			t.objects = append(t.objects, card)
		}
	}
	// add dealers's card resources to the widget renderer objects slice
	for _, card := range t.dealer.stack.cards {
		t.objects = append(t.objects, card)
	}
}

func (t *tableRender) newChipsAndBets() {
	for i := 0; i < 7; i++ {
		t.playersChips[i] = &canvas.Text{}
		t.chipImage[i] = &canvas.Image{}
		t.bets[i] = &canvas.Text{}
	}
}

func (t *tableRender) appendChipsAndBets() {
	for i := 0; i < 7; i++ {
		t.objects = append(t.objects, t.playersChips[i])
		t.objects = append(t.objects, t.chipImage[i])
		t.objects = append(t.objects, t.bets[i])
	}
}

func (t *tableRender) newPlayerActions() {
	for i := 0; i < 7; i++ {
		t.playerActions[i] = fyne.NewContainer()
		t.playerActions[i].Hide()
	}
}

func (t *tableRender) appendPlayerActions() {
	for i := 0; i < 7; i++ {
		t.objects = append(t.objects, t.playerActions[i])
	}
}

func (t *tableRender) newPlayerOutcomes() {
	for i := 0; i < 7; i++ {
		t.playerOutcomes[i] = &canvas.Image{}
		t.playerOutcomes[i].Hide()
	}
}

func (t *tableRender) appendPlayerOutcomes() {
	for i := 0; i < 7; i++ {
		t.objects = append(t.objects, t.playerOutcomes[i])
	}
}

// Instantiates new canvas objects and places them into the objects slice to be displayed on screen
func newTableRender(table *Table) *tableRender {

	render := &tableRender{}
	render.table = table
	render.game = table.game
	render.numPlayers = fyne.NewContainer()
	render.numPlayers.Hide()
	render.playerNames = fyne.NewContainer()
	render.newStacks()
	render.numPlayers = fyne.NewContainer()
	render.playerNames = fyne.NewContainer()
	render.playerNames.Hide()
	render.playerBets = fyne.NewContainer()
	render.playerBets.Hide()
	render.newPlayerActions()
	render.newNameAndScore()
	render.newChipsAndBets()
	render.newPlayerOutcomes()
	render.table.tabletop = &canvas.Image{Resource: imgBlackjackTableSvg}
	render.table.tabletop.Hide()
	render.resetTable = fyne.NewContainer()
	render.resetTable.Hide()

	render.objects = []fyne.CanvasObject{
		render.table.tabletop,
		render.numPlayers,
		render.resetTable,
	}
	render.appendStack()
	render.appendNameAndScore()
	render.appendChipsAndBets()
	render.objects = append(render.objects, render.playerNames, render.playerBets)
	render.appendPlayerActions()
	render.appendPlayerOutcomes()

	render.Refresh()
	return render
}

func newCardPos(card *Card) *canvas.Image {
	if card == nil {
		return &canvas.Image{}
	}

	image := &card.Image
	image.Resize(cardSize)

	return image
}
