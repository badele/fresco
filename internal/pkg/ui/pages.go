package ui

import (
	"fresco/internal/pkg/model"

	"github.com/rivo/tview"
)

type UPages struct {
	*tview.Pages
	List []string

	datas model.Bands
}

func NewPages() *UPages {
	pages := UPages{
		Pages: tview.NewPages(),
		List:  []string{},
		datas: []model.Band{},
	}

	pages.List = []string{
		"Logo",
		"Bands",
		"Channels",
		"Selections",
	}

	return &pages
}
