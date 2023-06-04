package ui

import "github.com/rivo/tview"

func NewStatusBar() *tview.TextView {
	statusbar := tview.NewTextView()
	statusbar.SetDynamicColors(true)

	return statusbar
}
