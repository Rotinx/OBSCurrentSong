package main

import (
	"fmt"
	"github.com/rivo/tview"
	"os/exec"
	"time"
)

/**
 *  @project OBSCurrentSong
 *  @author github.com/Rotinx
 */

func main() {
	config.loadConfig()

	err := exec.Command("cmd", "/C", "title", "OBSCurrentSong").Run()
	if err != nil {
		panic(err.Error())
	}

	app = tview.NewApplication()

	go fetcher()

	grid := initGrid()
	if err = app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func fetcher() {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				setStatus(fmt.Errorf("error: %v", r))
			}
		}()

		for _ = range ticker.C {
			song, err := currentSong()
			if err != nil {
				setStatus(err)
				continue
			}

			if lastSaved != *song {
				err = song.save()
				if err != nil {
					panic(err)
				}

				addSong(*song)
			}

			setStatus(nil)
		}
	}()

	select {}
}
