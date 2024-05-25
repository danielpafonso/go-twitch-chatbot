package ui

import (
	"errors"

	"github.com/awesome-gocui/gocui"
)

func (ui *UI) Layout(*gocui.Gui) error {
	maxX, maxY := ui.gui.Size()
	leftColumn := int(0.2 * float32(maxX))
	mainTop := -1
	mainOverlaps := 0
	if view, err := ui.gui.SetView("side", -1, -1, leftColumn, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		ui.sideView = view
	}

	if ui.bannerView {
		mainTop = 5
		mainOverlaps = gocui.TOP
		if view, err := ui.gui.SetView("title", leftColumn, -1, maxX, 5, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			view.WriteString(ui.bannerText)
		}
	}
	if view, err := ui.gui.SetView("main", leftColumn, mainTop, maxX, maxY-5, byte(mainOverlaps)); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Autoscroll = true
		view.Wrap = true
		//ui.gui.SetCurrentView("main")
		ui.mainView = view
		if !ui.bannerView {
			view.WriteString(ui.bannerText)
		}
		if ui.writeFromBottom {
			_, viewY := view.Size()
			view.SetWritePos(0, viewY-1)
		}
	}
	if view, err := ui.gui.SetView("cmdline", leftColumn, maxY-5, maxX, maxY, gocui.TOP); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Autoscroll = true
		view.Wrap = true
		ui.cmdView = view
	}
	return nil
}
