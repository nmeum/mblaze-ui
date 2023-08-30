package main

import (
	"github.com/gdamore/tcell/v2"
)

type UserInterface struct {
	Mails  []Mail
	Screen tcell.Screen
	index  int
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

func NewUI(mails []Mail, screen tcell.Screen) *UserInterface {
	return &UserInterface{
		Mails:  mails,
		Screen: screen,
		index:  0,
	}
}

func (ui *UserInterface) visible() int {
	_, ymax := ui.Screen.Size()
	return min(ymax, len(ui.Mails))
}

func (ui *UserInterface) SelectedMail() Mail {
	return ui.Mails[ui.index]
}

func (ui *UserInterface) IsSelected(m Mail) bool {
	selected := ui.SelectedMail()
	return selected.ID == m.ID
}

func (ui *UserInterface) NextMail() {
	if ui.index < len(ui.Mails)-1 {
		ui.index++
		ui.Draw()
	}
}

func (ui *UserInterface) PrevMail() {
	if ui.index > 0 {
		ui.index--
		ui.Draw()
	}
}

func (ui *UserInterface) Draw() {
	xmax, _ := ui.Screen.Size()
	if xmax <= 1 {
		panic("terminal is too small")
	}

	mails := ui.Mails
	visible := ui.visible()
	if ui.index >= visible {
		start := ui.index - visible
		mails = ui.Mails[start+1:]
	}

	y := 0
	for _, mail := range mails {
		text := mail.String()

		var style tcell.Style
		if ui.IsSelected(mail) {
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
		if y >= visible {
			break
		}
	}
}

func (ui *UserInterface) withoutScreen(fn func() error) error {
	ui.Screen.Fini()
	err := fn()
	if err != nil {
		return err
	}

	ui.Screen, err = initScreen()
	if err != nil {
		return err
	}
	ui.Draw()

	return nil
}
