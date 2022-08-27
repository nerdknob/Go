package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Table struct {
	widget.BaseWidget

	game     *Game
	tabletop *canvas.Image
}

func (t *Table) CreateRenderer() fyne.WidgetRenderer {
	return newTableRender(t)
}

func (t *Table) Tapped(event *fyne.PointEvent) {
	return
}

func NewTable(g *Game) *Table {
	table := &Table{game: g}
	table.ExtendBaseWidget(table)
	return table
}
