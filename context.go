package main

import (
	"github.com/gdamore/tcell/v2"
)

type Context struct {
	Mails  []Mail
	Index  *Index
	Screen tcell.Screen
}

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

func (ctx Context) Draw() {
	xmax, _ := ctx.Screen.Size()
	if xmax <= 1 {
		panic("terminal is too small")
	}

	y := 0
	for i, row := range ctx.Mails {
		text := row.Subject

		var style tcell.Style
		if ctx.Index.IsSelected(i) {
			style = selStyle
		} else {
			style = defStyle
		}

		truncated := false
		if len(text) >= xmax {
			text = text[0 : len(text)-1]
			truncated = true
		}
		drawText(ctx.Screen, y, 0, style, text)
		if truncated {
			// TODO: Determine cells needed for Abbreviated.
			ctx.Screen.SetContent(xmax-1, y, Abbreviated, nil, style)
		} else {
			for x := len(text); x < xmax; x++ {
				ctx.Screen.SetContent(x, y, ' ', nil, style)
			}
		}

		y++
	}
}
