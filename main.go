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

func fatal(ui *UserInterface, v ...any) {
	cleanup(ui)
	log.Fatal(v)
}

func handleEventKey(ui *UserInterface, ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEnter:
		err := ui.withoutScreen(func() error {
			mail := ui.SelectedMail()
			return mail.Show()
		})
		if err != nil {
			fatal(ui, err)
		}
	case tcell.KeyRune:
		mail := ui.SelectedMail()
		switch ev.Rune() {
		case 'r':
			err := ui.withoutScreen(func() error {
				mail := ui.SelectedMail()
				return mail.Reply()
			})
			if err != nil {
				fatal(ui, err)
			}
		case 'd':
			path, err := mail.Path()
			if err != nil {
				fatal(ui, err)
			}

			err = os.Remove(path)
			if err != nil {
				fatal(ui, err)
			}
		case 'f':
			mail.Flag(Flagged)
		case 'F':
			mail.Flag(Unflagged)
		case 's':
			mail.Flag(Unseen)
		case 'S':
			mail.Flag(Seen)
		case 't':
			mail.Flag(Untrashed)
		case 'T':
			mail.Flag(Trashed)
		default:
			return
		}

		// TODO: Consider using `mflag -v` above to determine the
		// new file names and thereby keep the sequence in tact.

		err := ui.Refresh()
		if err != nil {
			fatal(ui, err)
		}

		ui.Screen.Clear()
		ui.Draw()
	case tcell.KeyDown:
		ui.NextMail()
	case tcell.KeyUp:
		ui.PrevMail()
	case tcell.KeyPgDn:
		ui.NextPage()
	case tcell.KeyPgUp:
		ui.PrevPage()
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
			if ev.Key() == tcell.KeyEscape ||
				ev.Key() == tcell.KeyCtrlC ||
				(ev.Key() == tcell.KeyRune && ev.Rune() == 'q') {
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
