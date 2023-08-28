package main

import (
	"github.com/gdamore/tcell/v2"
)

type UserInterface struct {
	Mails  []Mail
	Index  *Index
	Screen tcell.Screen
}

const (
	// Rune used to indicate that the row has been abbreviated.
	Abbreviated = '…'
)

var (
	// Style used for non-selected rows.
	defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Style used for the currently selected row.
	selStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
)

func drawText(s tcell.Screen, row, col int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}

func (ui *UserInterface) SelectedMail() Mail {
	return ui.Mails[ui.Index.Cur()]
}

func (ui *UserInterface) NextMail() {
	ui.Index.Inc()
	ui.Draw()
}

func (ui *UserInterface) PrevMail() {
	ui.Index.Dec()
	ui.Draw()
}

func (ui *UserInterface) Draw() {
	xmax, _ := ui.Screen.Size()
	if xmax <= 1 {
		panic("terminal is too small")
	}

	y := 0
	for i, row := range ui.Mails {
		text := row.Subject

		var style tcell.Style
		if ui.Index.IsSelected(i) {
			style = selStyle
		} else {
			style = defStyle
		}

		truncated := false
		if len(text) >= xmax {
			text = text[0 : len(text)-1]
			truncated = true
		}
		drawText(ui.Screen, y, 0, style, text)
		if truncated {
			// TODO: Determine cells needed for Abbreviated.
			ui.Screen.SetContent(xmax-1, y, Abbreviated, nil, style)
		} else {
			for x := len(text); x < xmax; x++ {
				ui.Screen.SetContent(x, y, ' ', nil, style)
			}
		}

		y++
	}
}
