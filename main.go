package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
)

const (
	// Rune used to indicate that the row has been abbreviated.
	Abbreviated = '…'
)

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
		ctx.Draw()
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
		ctx.Draw()
	} else if ev.Key() == tcell.KeyDown {
		ctx.Index.Inc()
		ctx.Draw()
	} else if ev.Key() == tcell.KeyUp {
		ctx.Index.Dec()
		ctx.Draw()
	} else if ev.Key() == tcell.KeyCtrlL {
		ctx.Screen.Sync()
	}
}

func eventLoop(ctx *Context) {
	for {
		ctx.Screen.Show()
		ev := ctx.Screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			ctx.Draw()
			ctx.Screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}

			handleEventKey(ctx, ev)
		}
	}
}

func cleanup(ctx *Context) {
	// You have to catch panics in a defer, clean up, and
	// re-raise them - otherwise your application can
	// die without leaving any diagnostic trace.
	maybePanic := recover()

	ctx.Screen.Fini()
	if maybePanic != nil {
		panic(maybePanic)
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

	ctx := &Context{mails, idx, s}
	defer cleanup(ctx)

	ctx.Draw()
	eventLoop(ctx)
}
