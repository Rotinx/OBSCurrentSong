package main

import (
	"OBSCurrentSong/dialogs"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/**
 *  @project OBSCurrentSong
 *  @author github.com/Rotinx
 */

var app *tview.Application

var mainChildPages *tview.Pages
var mainChildCurrentSong *tview.TextView
var mainChildStatus *tview.TextView

var previouslyPlayedChild *tview.Frame
var previousPlayedChildList *tview.List

// TODO: improve window/state management.

// initGrid initializes parent elements of the grid.
func initGrid() *tview.Grid {
	mainChildPages = tview.NewPages().AddPage("main", mainPageFrame(), true, true).AddPage("settings", settingsPageFrame(), true, false)

	previousPlayedChildList = tview.NewList()
	previouslyPlayedChild = tview.NewFrame(previousPlayedChildList).AddText("Previously Played", true, tview.AlignLeft, tcell.ColorWhite)

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(0, 0, 0, 0).
		SetBorders(true)

	grid.AddItem(mainChildPages, 0, 0, 1, 2, 0, 0, false).
		AddItem(previouslyPlayedChild, 0, 2, 1, 2, 0, 0, false)

	return grid
}

// mainPageFrame initializes the main page frame (left side of grid).
func mainPageFrame() *tview.Frame {
	mainChildCurrentSong = tview.NewTextView().SetText("Current song: ...")
	mainChildStatus = tview.NewTextView().SetText("Status: Initializing").SetTextColor(tcell.ColorGreen)

	flex := tview.NewFlex().
		AddItem(mainChildCurrentSong, 1, 1, false).
		AddItem(mainChildStatus, 0, 1, false).SetDirection(tview.FlexRow).
		AddItem(tview.NewButton("Settings").SetSelectedFunc(func() {
			switchToPage("settings")
		}), 3, 1, false)

	return tview.NewFrame(flex)
}

// settingsPageFrame initializes the settings page frame.
func settingsPageFrame() *tview.Frame {
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(settingsOutputDir(), 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(tview.NewButton("Back").SetSelectedFunc(func() {
			switchToPage("main")
		}), 3, 1, false)
	return tview.NewFrame(flex)
}

// settingsOutputDir initializes the output directory setting.
func settingsOutputDir() *tview.Flex {
	flex := tview.NewFlex()

	settingsSavePath := tview.NewTextView().SetText(config.SavePath).SetTextColor(tcell.ColorAqua)
	button := tview.NewButton("Set").SetSelectedFunc(func() {
		path, err := dialogs.Directory().SetStartDir(getExecutionPath()).Title("Select output location").Browse()
		if err != nil {
			// TODO: handle error, currently silently fails.
			return
		}

		settingsSavePath.SetText(path)

		config.SavePath = path
		config.save()
	})

	flex.AddItem(tview.NewTextView().SetText("Output Dir"), 0, 1, false).AddItem(button, 20, 1, false)

	wrapper := tview.NewFlex().SetDirection(tview.FlexRow)
	return wrapper.AddItem(flex, 1, 1, false).AddItem(settingsSavePath, 1, 1, false)
}
