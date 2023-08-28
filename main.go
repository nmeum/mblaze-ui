package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
)

func initScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	s.SetStyle(defStyle)
	s.EnablePaste()
	s.Clear()
	return s, nil
}

func handleEventKey(ui *UserInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEnter:
		ui.Screen.Fini()
		mail := ui.SelectedMail()
		err := mail.Show()
		if err != nil {
			log.Fatal(err)
		}
		ui.Screen, err = initScreen()
		if err != nil {
			log.Fatal(err)
		}
		ui.Draw()
	case tcell.KeyRune:
		mail := ui.SelectedMail()
		switch ev.Rune() {
		case 'd':
			path, err := mail.Path()
			if err != nil {
				log.Fatal(err)
			}

			err = os.Remove(path)
			if err != nil {
				log.Fatal(err)
			}
		case 's':
			mail.Flag(Unseen)
		case 'S':
			mail.Flag(Seen)
		case 'f':
			mail.Flag(Flagged)
		case 'F':
			mail.Flag(Unflagged)
		default:
			return
		}

		var err error
		ui.Mails, err = mscan()
		if err != nil {
			log.Fatal(err)
		}
		ui.Draw()
	case tcell.KeyDown:
		ui.NextMail()
	case tcell.KeyUp:
		ui.PrevMail()
	case tcell.KeyCtrlL:
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
	mails, err := mscan()
	if err != nil {
		log.Fatal(err)
	}

	screen, err := initScreen()
	if err != nil {
		log.Fatal(err)
	}

	ui := NewUI(mails, screen)
	defer cleanup(ui)

	ui.Draw()
	eventLoop(ui)
}
