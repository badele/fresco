package ui

import (
	"fresco/internal/pkg/model"

	"fresco/internal/pkg/tools"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UBands struct {
	*tview.Table

	filter string
	Datas  *model.Bands
}

func NewBandsTable(bands *model.Bands) *UBands {
	table := tview.NewTable()

	table.SetSelectable(true, false)
	table.Select(0, 0)
	table.SetFixed(1, 1)

	uibands := &UBands{
		Table:  table,
		filter: "",
		Datas:  bands,
	}
	uibands.SetBorder(true)

	return uibands
}

func (b *UBands) Refresh(filter string) {

	if filter != "" && filter == b.filter {
		return
	}

	b.Clear()

	headers := []string{
		"Name",
		"Location",
		"Lower Frequency",
		"Upper Frequency",
		"Source",
		"Type",
	}

	for i, header := range headers {
		b.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tview.Styles.PrimaryTextColor,
			BackgroundColor: tview.Styles.PrimitiveBackgroundColor,
			Attributes:      tcell.AttrBold,
		})
	}

	var cell *tview.TableCell
	idx := 0
	for i, band := range *b.Datas {
		if filter != "" && !(band.IsFiltered(filter)) {
			continue
		}

		idx += 1
		col := 0
		cell = tview.NewTableCell("[darkcyan]" + band.Name)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		cell.SetReference(i)
		b.SetCell(idx+1, col, cell)

		col += 1
		cell = tview.NewTableCell(band.Location)
		cell.SetTextColor(tview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		cell.SetReference(i)
		b.SetCell(idx+1, col, cell)

		col += 1
		cell = tview.NewTableCell("[green]" + tools.ConvertFrequencyToString(float64(band.LowerFrequency)))
		cell.SetTextColor(tview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		cell.SetReference(i)
		b.SetCell(idx+1, col, cell)

		col += 1
		cell = tview.NewTableCell("[green]" + tools.ConvertFrequencyToString(float64(band.UpperFrequency)))
		cell.SetTextColor(tview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		cell.SetReference(i)
		b.SetCell(idx+1, col, cell)

		col += 1
		cell = tview.NewTableCell(band.Source)
		cell.SetTextColor(tview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		cell.SetReference(i)
		b.SetCell(idx+1, col, cell)

		col += 1
		cell = tview.NewTableCell("[orange]" + band.Type)
		cell.SetTextColor(tview.Styles.PrimaryTextColor)
		cell.SetMaxWidth(30)
		cell.SetExpansion(1)
		cell.SetReference(i)
		b.SetCell(idx+1, col, cell)
	}
	b.Select(0, 1)
}
