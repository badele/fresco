package ui

import (
	"fresco/internal/pkg/model"

	"github.com/rivo/tview"
)

type UPageMenu struct {
	Title   string
	Visible bool
}

type UPages struct {
	*tview.Pages
	MenuNames []UPageMenu

	datas model.Bands

	// func AddMenuPage(p *UPages) (name string, item tview.Primitive)
}

func NewPages() *UPages {
	pages := &UPages{
		Pages:     tview.NewPages(),
		MenuNames: []UPageMenu{},
		datas:     []model.Band{},
	}

	return pages
}

func (p *UPages) AddMenuPage(name string, menuvisible bool, item tview.Primitive) {
	p.AddPage(name, item, true, true)
	p.MenuNames = append(p.MenuNames, UPageMenu{
		Title:   name,
		Visible: menuvisible,
	})
}
