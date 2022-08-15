package main

import (
	"fmt"
	"log"
	"os"
)

/**
 *  @project OBSCurrentSong
 *  @author github.com/Rotinx
 */

var lastSaved Song

func _write(filename string, content string) error {
	bytes := []byte(content)

	err := os.WriteFile(fmt.Sprintf("./%v.txt", filename), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Song) Save() error {
	if lastSaved == *s {
		return nil
	}

	err := _write("artist", s.Artist)
	if err != nil {
		return fmt.Errorf("couldn't save artist file. %v", err)
	}

	err = _write("song", s.Name)
	if err != nil {
		return fmt.Errorf("couldn't save song file. %v", err)
	}

	err = _write("entire", fmt.Sprintf("%v - %v", s.Artist, s.Name))
	if err != nil {
		return fmt.Errorf("couldn't save artist & song file. %v", err)
	}

	lastSaved = *s
	log.Printf("Artist: %v, Song: %v", s.Artist, s.Name)

	return nil
}
