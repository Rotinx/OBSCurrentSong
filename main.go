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
	fmt.Printf("OBSCurrentSong")

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
				if errors.Is(err, NoTitleFound) {
					fmt.Println("No title found, please make sure spotify is currently open.")
				}

				if errors.Is(err, NoSongPlaying) {
					fmt.Println("No song is currently playing.")
				}

				continue
			}

			if lastSaved != *song {
				err = song.Save()
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	select {}
}
