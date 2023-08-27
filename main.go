package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
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

const (
	// Rune used to indicate that the row has been abbreviated.
	Abbreviated = '…'
)

func drawText(s tcell.Screen, row, col int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}

func drawRows(s tcell.Screen, idx *Index, rows []Mail) {
	xmax, _ := s.Size()
	if xmax <= 1 {
		panic("terminal is too small")
	}

	y := 0
	for i, row := range rows {
		text := row.Subject

		var style tcell.Style
		if idx.IsSelected(i) {
			style = selStyle
		} else {
			style = defStyle
		}

		truncated := false
		if len(text) >= xmax {
			text = text[0 : len(text)-1]
			truncated = true
		}
		drawText(s, y, 0, style, text)
		if truncated {
			// TODO: Determine cells needed for Abbreviated.
			s.SetContent(xmax-1, y, Abbreviated, nil, style)
		} else {
			for x := len(text); x < xmax; x++ {
				s.SetContent(x, y, ' ', nil, style)
			}
		}

		y++
	}
}

func initScreen() tcell.Screen {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}
	s.SetStyle(defStyle)
	s.EnablePaste()
	s.Clear()
	return s
}

func handleEventKey(ctx *Context, ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEnter {
		ctx.Screen.Fini()
		mail := ctx.Mails[ctx.Index.Cur()]
		err := mblaze_show(mail)
		if err != nil {
			log.Fatal(err)
		}
		ctx.Screen = initScreen()
		drawRows(ctx.Screen, ctx.Index, ctx.Mails)
	} else if ev.Key() == tcell.KeyRune {
		mail := ctx.Mails[ctx.Index.Cur()]
		switch ev.Rune() {
		case 's':
			mblaze_flag(mail, Seen)
		case 'f':
			mblaze_flag(mail, Flagged)
		}

		var err error
		ctx.Mails, err = mblaze_mscan()
		if err != nil {
			log.Fatal(err)
		}
		drawRows(ctx.Screen, ctx.Index, ctx.Mails)
	} else if ev.Key() == tcell.KeyDown {
		ctx.Index.Inc()
		drawRows(ctx.Screen, ctx.Index, ctx.Mails)
	} else if ev.Key() == tcell.KeyUp {
		ctx.Index.Dec()
		drawRows(ctx.Screen, ctx.Index, ctx.Mails)
	} else if ev.Key() == tcell.KeyCtrlL {
		ctx.Screen.Sync()
	}
}

func main() {
	mails, err := mblaze_mscan()
	if err != nil {
		log.Fatal(err)
	}

	s := initScreen()
	idx := NewIndex(func() int {
		_, ymax := s.Size()
		return min(ymax, len(mails))
	})
	drawRows(s, idx, mails)

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	ctx := &Context{mails, idx, s}
	for {
		ctx.Screen.Show()
		ev := ctx.Screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			drawRows(s, idx, mails)
			ctx.Screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}

			handleEventKey(ctx, ev)
		}
	}
}
