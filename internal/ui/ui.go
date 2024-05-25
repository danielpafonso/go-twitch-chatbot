package ui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	bannerView      bool
	bannerText      string
	writeFromBottom bool
	gui             *gocui.Gui
	mainView        *gocui.View
	cmdView         *gocui.View
	sideView        *gocui.View
}

func NewUI(showView bool, banner string) *UI {
	return &UI{
		bannerView:      showView,
		bannerText:      banner,
		writeFromBottom: true,
	}
}

func (ui *UI) Write(data []byte) (n int, err error) {
	ui.gui.UpdateAsync(func(g *gocui.Gui) error {
		//ui.mainView.Write(data)
		fmt.Fprint(ui.mainView, string(data))
		return nil
	})
	return len(data), nil
}

func (ui *UI) WriteMain(text string) {
	// func (g *Gui) Update(f func(*Gui) error)
	ui.gui.UpdateAsync(func(g *gocui.Gui) error {
		//fmt.Fprintf(ui.gui.CurrentView(), ">%s\n", text)
		fmt.Fprintf(ui.mainView, " %s", text)
		return nil
	})
}

func (ui *UI) WriteCmd(text string) {
	// func (g *Gui) Update(f func(*Gui) error)
	ui.gui.UpdateAsync(func(g *gocui.Gui) error {
		//fmt.Fprintf(ui.gui.CurrentView(), ">%s\n", text)
		fmt.Fprint(ui.cmdView, text)
		return nil
	})
}

func (ui *UI) WriteSide(text string, clearView ...bool) {
	clear := false
	if len(clearView) > 0 {
		clear = clearView[0]
	}
	// func (g *Gui) Update(f func(*Gui) error)
	ui.gui.UpdateAsync(func(g *gocui.Gui) error {
		//fmt.Fprintf(ui.gui.CurrentView(), ">%s\n", text)
		if clear {
			ui.sideView.Clear()
		}
		fmt.Fprint(ui.sideView, text)
		return nil
	})
}

func (ui *UI) Start(startedUI chan struct{}) error {
	var err error
	// create terminal gui
	ui.gui, err = gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return err
	}
	defer ui.gui.Close()
	// set graphical manager
	ui.gui.SetManagerFunc(ui.Layout)

	// set keybindings
	ui.SetKeybindings()

	// notify that UI is started
	startedUI <- struct{}{}

	// enter UI mainloop
	if err := ui.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}
