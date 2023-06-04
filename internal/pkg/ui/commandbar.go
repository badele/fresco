package ui

import "github.com/rivo/tview"

func NewCommandBar() *tview.InputField {
	commandbar := tview.NewInputField()
	commandbar.SetDisabled(true)
	commandbar.SetBorder((false))
	commandbar.SetLabel("")

	return commandbar
}
