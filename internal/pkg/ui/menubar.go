package ui

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

type MenuBar struct {
	*tview.TextView
	menus    []string
	selected int
}

func (m *MenuBar) AddEntry(entry string) {
	m.menus = append(m.menus, entry)

	labels := ""
	for idx, label := range m.menus {
		labels += fmt.Sprintf(`%d ["%d"][darkcyan]%s[white][""]  `, idx+1, idx, label)
	}

	m.SetText(labels)
	m.Refresh()
}

func (m *MenuBar) Next() {
	m.selected = (m.selected + 1) % len(m.menus)
	m.Refresh()
}

func (m *MenuBar) Prev() {
	m.selected = (m.selected - 1 + len(m.menus)) % len(m.menus)
	m.Refresh()
}

func (m *MenuBar) Refresh() {
	m.Highlight(strconv.Itoa(m.selected))
	m.ScrollToHighlight()
}

func NewMenuBar() *MenuBar {
	menubar := MenuBar{
		TextView: tview.NewTextView(),
		menus:    []string{},
		selected: 0,
	}

	menubar.SetDynamicColors(true)
	menubar.SetRegions(true)
	menubar.SetWrap(false)

	return &menubar
}
