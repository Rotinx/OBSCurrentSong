package main

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"time"
)

/**
 *  @project OBSCurrentSong
 *  @author github.com/Rotinx
 */

/**
 *  addSong sets the current playing song, and adds it to the previously played list.
 */
func addSong(song songSingleton) {
	app.QueueUpdateDraw(func() {
		t := time.Now()
		previousPlayedChildList.SetSelectedFocusOnly(true)
		previousPlayedChildList.InsertItem(0, fmt.Sprintf("%s - %s", song.Name, song.Artist), fmt.Sprintf(t.Local().Format(time.ANSIC)), 0, nil)
		mainChildCurrentSong.SetText(fmt.Sprintf("Current song: %s - %s", song.Name, song.Artist))
	})
}

/**
 * setStatus sets the status of the application.
 */
func setStatus(err error) {
	app.QueueUpdateDraw(func() {
		if err != nil {
			mainChildStatus.SetText(fmt.Sprintf("Status: %s", getErrorReason(err))).SetTextColor(tcell.ColorRed)
		} else {
			mainChildStatus.SetText("Status: OK").SetTextColor(tcell.ColorGreen)
		}

	})
}

/**
 *  switchToPage switches the current page to the given page (only main child page).
 */
func switchToPage(page string) {
	mainChildPages.SwitchToPage(page)
}

/**
 *  getErrorReason Breaks down errors into a more readable format.
 */
func getErrorReason(err error) string {
	if errors.Is(err, NoSongPlaying) {
		return "No song is currently playing."
	}

	if errors.Is(err, NoTitleFound) {
		return "Please make sure spotify is currently open."
	}

	return "Unknown error, please restart the application."
}

func getExecutionPath() string {
	return os.Args[0]
}
