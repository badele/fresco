package ui

import (
	"fmt"
	"fresco/internal/pkg/model"
	"strconv"
	"strings"

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

	pages     *UPages
	menubar   *MenuBar
	statusbar *tview.TextView

	// Bands pannel
	bands *UBands

	// Filter
	commandbar *tview.InputField
	filter     string

	// Datas
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

	a.statusbar.SetText("Tab: [darkcyan]Next[white],[darkcyan]Prev[white]  /: [darkcyan]Filter[white]  Ctrl-Q: [darkcyan]Quit[white]")
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

func (a *App) InitLogoPage() *tview.Flex {
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

	return flex
}

// Init Bands page
func (a *App) InitBandsPage() *UBands {
	bands := NewBandsTable(&a.data_bands)
	bands.Refresh("")

	return bands
}

func (a *App) initPanels() {

	// Init pages
	plogo := a.InitLogoPage()
	pbands := a.InitBandsPage()

	a.bands = pbands

	a.pages.AddMenuPage("Logo", false, plogo)
	a.pages.AddMenuPage("Bands", true, pbands)

	// Define menu selected function(hightlight)
	a.menubar.SetHighlightedFunc(func(added, removed, remaining []string) {
		idx, _ := strconv.Atoi(added[0])
		menupage := a.pages.MenuNames[idx]

		a.statusbar.SetText(menupage.Title)
		a.GoTo(menupage.Title)
	})

	// Define flex layout
	layout := tview.NewFlex()
	layout.SetDirection(tview.FlexRow)
	layout.AddItem(a.menubar, 1, 1, false)
	layout.AddItem(a.pages, 0, 16, true)
	layout.AddItem(a.commandbar, 1, 1, false)
	layout.AddItem(a.statusbar, 1, 1, false)

	// Generate Menu from pages list
	for _, menuname := range a.pages.MenuNames {
		if !menuname.Visible {
			continue
		}
		a.menubar.AddEntry(menuname.Title)
	}
	a.menubar.Refresh()

	// Show logo
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
				if event.Key() == tcell.KeyTab {
					a.menubar.Next()
					return nil
				} else if event.Key() == tcell.KeyBacktab {
					a.menubar.Prev()
					return nil
				} else if event.Key() == tcell.KeyCtrlC {
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
					return nil
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
		pages:      NewPages(),
		menubar:    NewMenuBar(),
		statusbar:  NewStatusBar(),
		bands:      &UBands{},
		commandbar: NewCommandBar(),
		filter:     "",
		data_bands: model.NewBands("datasets/bands.fwf"),
	}

	return &app
}
