package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"fresco/internal/pkg/model"
	"fresco/internal/pkg/model/importer"

	"fresco/internal/pkg/ui"

	"github.com/alecthomas/kong"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const logo = `
███████╗██████╗ ███████╗███████╗ ██████╗ ██████╗ 
██╔════╝██╔══██╗██╔════╝██╔════╝██╔════╝██╔═══██╗
█████╗  ██████╔╝█████╗  ███████╗██║     ██║   ██║
██╔══╝  ██╔══██╗██╔══╝  ╚════██║██║     ██║   ██║
██║     ██║  ██║███████╗███████║╚██████╗╚██████╔╝
╚═╝     ╚═╝  ╚═╝╚══════╝╚══════╝ ╚═════╝ ╚═════╝ 
`
const logoTitle = "A radio frequency database"

type AppMode int

const (
	LOGOMODE AppMode = iota
	NORMALMODE
	FILTERMODE
)

type App struct {
	app     *tview.Application
	appmode AppMode

	pages     *ui.UPages
	menubar   *ui.MenuBar
	statusbar *tview.TextView

	// Bands pannel
	bands *ui.UBands

	// Filter
	commandbar *tview.InputField
	filter     string

	data_bands model.Bands
}

func (a *App) LogoMode() {
	a.commandbar.SetBorder(false)
	a.commandbar.SetText("")
	a.statusbar.SetText("")
	a.GoTo("Logo")
}

func (a *App) NormalMode() {
	a.appmode = NORMALMODE

	if len(a.commandbar.GetText()) > 0 {
		a.commandbar.SetLabel("Filter: ")
	} else {
		a.commandbar.SetLabel("")
	}

	a.commandbar.SetPlaceholder("")
	a.commandbar.SetDisabled(true)
	a.app.SetFocus(a.pages)

	a.statusbar.SetText("/: [darkcyan]Filter[white]  Ctrl-S: [darkcyan]Add to selection[white]  Ctrl-I: [darkcyan]Import[white]  Ctrl-Q: [darkcyan]Quit[white]")
}

func (a *App) UpdateFilter(filter string) {
	if filter != a.filter {
		a.filter = filter
		a.bands.Refresh(filter)
	}

	a.commandbar.SetText(filter)
	a.NormalMode()
}

func (a *App) FilterMode() {
	a.appmode = FILTERMODE

	a.commandbar.SetLabel("Filter: ")
	a.commandbar.SetPlaceholder("E.g. FIXE|MOBILE|!NAVIGATION")
	a.commandbar.SetDisabled(false)
	a.commandbar.SetChangedFunc(func(content string) {
		a.bands.Refresh(content)
	})
	a.app.SetFocus(a.commandbar)
	a.statusbar.SetText(`[green]Filter Pattern[white]: |: [darkcyan]Or[white]  !: [darkcyan]Not[white]  |  [green]Keys[white]: Enter: [darkcyan]Validate filter[white]  Esc: [darkcyan]Clear filter[white]`)
}

func (a *App) InitLogoPage() {
	// Compute logo width
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)
	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}
	logoBox := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorGreen).
		SetDoneFunc(func(key tcell.Key) {
			a.NormalMode()
		})
	fmt.Fprint(logoBox, logo)

	frame := tview.NewFrame(tview.NewBox()).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(logoTitle, true, tview.AlignCenter, tcell.ColorWhite)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 7, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logoBox, logoWidth, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true).
		AddItem(frame, 0, 10, false)

	a.pages.AddPage("Logo", flex, true, true)
}

// Init Bands page
func (a *App) InitBandsPage() {
	a.bands = ui.NewBandsTable()
	a.bands.Refresh("")

	a.pages.AddPage("Bands", a.bands, true, true)
}

func (a *App) InitChannelsPage() {
	a.pages.AddPage("Channels", nil, true, true)
}

func (a *App) InitLogsPage() {
	a.pages.AddPage("Logs", nil, true, true)
}

func (a *App) initPanels() {

	a.InitLogoPage()
	a.InitBandsPage()

	a.menubar.SetHighlightedFunc(func(added, removed, remaining []string) {
		idx, _ := strconv.Atoi(added[0])
		pagename := a.pages.List[idx+1]

		a.statusbar.SetText(pagename)
		a.GoTo(pagename)
	})

	layout := tview.NewFlex()
	layout.SetDirection(tview.FlexRow)
	layout.AddItem(a.menubar, 1, 1, false)
	layout.AddItem(a.pages, 0, 16, true)
	layout.AddItem(a.commandbar, 1, 1, false)
	layout.AddItem(a.statusbar, 1, 1, false)

	for idx, menuname := range a.pages.List {
		if idx == 0 {
			continue
		}
		a.menubar.AddEntry(menuname)
	}
	a.menubar.Refresh()

	a.LogoMode()

	a.app.SetRoot(layout, true)
}

func (a *App) GoTo(page string) {
	a.pages.SwitchToPage(page)
}

func (a *App) Start() error {
	a.initPanels()

	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch a.appmode {
		case LOGOMODE:
			{
				if event.Key() == tcell.KeyCtrlQ {
					a.Stop()
					return nil
				} else {
					a.GoTo("Bands")
					a.NormalMode()
					return nil
				}
			}
		case NORMALMODE:
			{
				if event.Key() == tcell.KeyCtrlN {
					a.menubar.Next()
					return nil
				} else if event.Key() == tcell.KeyCtrlP {
					a.menubar.Prev()
					return nil
				} else if event.Key() == tcell.KeyCtrlQ {
					a.app.Stop()
					return nil
				} else if event.Name() == "Rune[/]" {
					a.FilterMode()
					return nil
				} else if event.Key() == tcell.KeyEnter {
					a.UpdateFilter(a.commandbar.GetText())
					return nil
				} else {
					return event
				}
			}
		case FILTERMODE:
			{
				if event.Key() == tcell.KeyEsc {
					a.UpdateFilter("")
				} else {
					if event.Key() == tcell.KeyEnter {
						a.UpdateFilter(a.commandbar.GetText())
						return nil
					}
				}
			}
		}

		return event
	})

	if err := a.app.Run(); err != nil {
		a.app.Stop()
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.app.Stop()
}

func NewApp() *App {
	app := App{
		app:        tview.NewApplication(),
		appmode:    LOGOMODE,
		pages:      ui.NewPages(),
		menubar:    ui.NewMenuBar(),
		statusbar:  ui.NewStatusBar(),
		commandbar: ui.NewCommandBar(),
		filter:     "",
		data_bands: model.NewBands("datasets/bands.fwf"),
	}

	return &app
}

func ImportBands() {
	oldbands := model.Bands{}
	oldbands.Read("datasets/bands.fwf")

	newbands := model.Bands{
		model.Band{
			Name:           "",
			Location:       "",
			LowerFrequency: 0,
			UpperFrequency: 0,
			Source:         "",
			Type:           "",
		},
		model.Band{Name: "SLF", Location: "Worldwide", LowerFrequency: 3e+1, UpperFrequency: 3e+2, Source: "IEEE", Type: "band"},
		model.Band{Name: "ULF", Location: "Worldwide", LowerFrequency: 3e+2, UpperFrequency: 3e+3, Source: "IEEE", Type: "band"},
		model.Band{Name: "VLF", Location: "Worldwide", LowerFrequency: 3e+3, UpperFrequency: 3e+4, Source: "IEEE", Type: "band"},
		model.Band{Name: "LF", Location: "Worldwide", LowerFrequency: 3e+4, UpperFrequency: 3e+5, Source: "IEEE", Type: "band"},
		model.Band{Name: "MF", Location: "Worldwide", LowerFrequency: 3e+5, UpperFrequency: 3e+6, Source: "IEEE", Type: "band"},
		model.Band{Name: "HF", Location: "Worldwide", LowerFrequency: 3e+6, UpperFrequency: 3e+7, Source: "IEEE", Type: "band"},
		model.Band{Name: "VHF", Location: "Worldwide", LowerFrequency: 3e+7, UpperFrequency: 3e+8, Source: "IEEE", Type: "band"},
		model.Band{Name: "UHF", Location: "Worldwide", LowerFrequency: 3e+8, UpperFrequency: 3e+9, Source: "IEEE", Type: "band"},
		model.Band{Name: "SHF", Location: "Worldwide", LowerFrequency: 3e+9, UpperFrequency: 3e+10, Source: "IEEE", Type: "band"},
		model.Band{Name: "EHF", Location: "Worldwide", LowerFrequency: 3e+10, UpperFrequency: 3e+11, Source: "IEEE", Type: "band"},
		model.Band{Name: "THF", Location: "Worldwide", LowerFrequency: 3e+11, UpperFrequency: 3e+12, Source: "IEEE", Type: "band"},
	}

	var dbimporter importer.Importer

	// IEEE
	dbimporter = &importer.IEEEBands{}
	newbands = importer.Import(newbands, dbimporter)

	// SDRPP
	dbimporter = &importer.SDRPPBands{}
	newbands = importer.Import(newbands, dbimporter)

	// GQRX
	dbimporter = &importer.GQRXBands{}
	newbands = importer.Import(newbands, dbimporter)

	// Merge only news band data
	oldbands.InsertNewValue(&newbands)

	// Save
	sort.Slice(oldbands, func(i, j int) bool {
		return oldbands[i].UUID() < oldbands[j].UUID()
	})

	// Write result file
	oldbands.Write("datasets/bands.fwf")
}

var cli struct {
	Import struct {
		Dataset string `enum:"bands,channels" required:"" short:"d" help:"Dataset name [ bands|channels ]"`
	} `cmd:"" help:"Import dataset"`
}

func main() {
	if len(os.Args) >= 2 {
		ctx := kong.Parse(&cli)
		switch ctx.Command() {
		case "import":
			dataset := strings.ToLower(cli.Import.Dataset)
			switch dataset {
			case "bands":
				fmt.Printf("Importing %s ...\n", dataset)
				ImportBands()
			}
		default:
			fmt.Print(ctx.Command())
		}
	} else {
		app := NewApp()
		app.Start()
	}
}
