package main

import (
	"github.com/gdamore/tcell/v2"
	"time"
)

var (
	curDate = time.Now()
)

// Adaptive time printing depending on the distance to the current date.
//
// Inspired by https://github.com/leahneukirchen/mblaze/blob/v1.2/mscan.c#L179-L184
func adaptiveTime(t time.Time) string {
	if t.Year() != curDate.Year() {
		return t.Format("2006-01-02")
	} else if t.After(curDate) || curDate.Sub(t) > 24*time.Hour {
		return t.Format("Mon Jan 02")
	} else {
		return t.Format("Mon 15:04")
	}
}

func drawText(s tcell.Screen, row, col int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}
