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

func handleEventKey(ui *UserInterface, ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEnter {
		ui.Screen.Fini()
		mail := ui.SelectedMail()
		err := mblaze_show(mail)
		if err != nil {
			log.Fatal(err)
		}
		ui.Screen = initScreen()
		ui.Draw()
	} else if ev.Key() == tcell.KeyRune {
		mail := ui.SelectedMail()
		switch ev.Rune() {
		case 's':
			mblaze_flag(mail, Seen)
		case 'f':
			mblaze_flag(mail, Flagged)
		}

		var err error
		ui.Mails, err = mblaze_mscan()
		if err != nil {
			log.Fatal(err)
		}
		ui.Draw()
	} else if ev.Key() == tcell.KeyDown {
		ui.NextMail()
	} else if ev.Key() == tcell.KeyUp {
		ui.PrevMail()
	} else if ev.Key() == tcell.KeyCtrlL {
		ui.Screen.Sync()
	}
}

func eventLoop(ui *UserInterface) {
	for {
		ui.Screen.Show()
		ev := ui.Screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			ui.Draw()
			ui.Screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}

			handleEventKey(ui, ev)
		}
	}
}

func cleanup(ui *UserInterface) {
	// You have to catch panics in a defer, clean up, and
	// re-raise them - otherwise your application can
	// die without leaving any diagnostic trace.
	maybePanic := recover()

	ui.Screen.Fini()
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

	ui := &UserInterface{mails, idx, s}
	defer cleanup(ui)

	ui.Draw()
	eventLoop(ui)
}
