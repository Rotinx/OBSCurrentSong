package main

import (
	"errors"
	"fmt"
	"time"
)

/**
 *  @project OBSCurrentSong
 *  @author github.com/Rotinx
 */

func main() {
	fmt.Println("OBSCurrentSong")

	var lastError error = nil

	ticker := time.NewTicker(2 * time.Second)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("While fetching the current song we encountered a error (Please resolve the issue and restart the application).\n", r)
			}
		}()

		for _ = range ticker.C {
			song, err := currentSong()
			if err != nil {
				if errors.Is(err, NoTitleFound) && lastError != NoTitleFound {
					fmt.Println("No title found, please make sure spotify is currently open.")
					lastError = NoTitleFound
				}

				if errors.Is(err, NoSongPlaying) && lastError != NoSongPlaying {
					fmt.Println("No song is currently playing.")
					lastError = NoSongPlaying
				}

				continue
			}

			if lastSaved != *song {
				err = song.Save()
				if err != nil {
					panic(err)
				}
				lastError = nil
			}
		}
	}()

	select {}
}
