package main

import (
	"OBSCurrentSong/ps"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

/**
 *  @project OBSCurrentSong
 *  @author github.com/Rotinx
 */

var user32 = syscall.NewLazyDLL("user32.dll")
var procGetWindowW = user32.NewProc("GetWindowTextW")
var procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
var procEnumWindows = user32.NewProc("EnumWindows")
var procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
var procIsWindowVisible = user32.NewProc("IsWindowVisible")

var callbackFunc uintptr
var spotifyTitle string

var NoSongPlaying = errors.New("no song playing")
var NoTitleFound = errors.New("title couldn't be found")
var UnknownError = errors.New("unknown error")

type songSingleton struct {
	Name   string
	Artist string
}

var lastSaved songSingleton

// Cache of all songs that have played (in this session).
var songCache []songSingleton

func currentSong() (*songSingleton, error) {
	pid, err := _fetchPid()
	if err != nil {
		return nil, err
	}

	if callbackFunc == 0 {
		callbackFunc = syscall.NewCallback(_fetchTitle)
	}

	success, _, _ := procEnumWindows.Call(callbackFunc, pid)
	if success != 0 {
		return nil, NoTitleFound
	}

	split := strings.SplitN(spotifyTitle, " - ", 2)

	if len(split) < 2 {
		return nil, NoSongPlaying
	}

	artist, name := split[0], split[1]

	song := &songSingleton{
		Name:   name,
		Artist: artist,
	}

	return song, nil
}

// _fetchTitle obtains title from provided pid.
func _fetchTitle(h syscall.Handle, pid uintptr) uintptr {
	handle := uintptr(h)

	/* 	Return types

	1 continues loop.
	0 ends loop.
	*/

	var _loopPid uintptr
	procGetWindowThreadProcessId.Call(handle, uintptr(unsafe.Pointer(&_loopPid)))

	if _loopPid == pid {

		// TODO: improve filtering process.
		_isVisible, _, _ := procIsWindowVisible.Call(handle)
		if _isVisible == 0 {
			return 1
		}

		_textLength, _, _ := procGetWindowTextLengthW.Call(handle)
		if _textLength == 0 {
			return 1
		}

		_textBuff := make([]uint16, _textLength+1)
		_titleRes, _, _ := procGetWindowW.Call(handle, uintptr(unsafe.Pointer(&_textBuff[0])), uintptr(int32(len(_textBuff))))
		if _titleRes == 0 {
			return 1
		}

		spotifyTitle = syscall.UTF16ToString(_textBuff)
		return 0
	}

	return 1
}

// _fetchPid attempts to fetch the current Spotify pid.
func _fetchPid() (uintptr, error) {
	var pid = 0

	processes, err := ps.Processes()
	if err != nil {
		return 0, err
	}

	for _, process := range processes {

		// TODO: improve filtering process.
		if process.Executable() == "Spotify.exe" {
			pid = process.Pid()
			break
		}
	}

	/*	References
		https://go.dev/play/p/YfGDtIuuBw
		https://stackoverflow.com/a/14029835
	*/

	if callbackFunc == 0 {
		callbackFunc = syscall.NewCallback(_fetchTitle)
	}

	return uintptr(pid), nil
}

func _writeSongFile(filename string, content string) error {
	bytes := []byte(content)

	err := os.WriteFile(fmt.Sprintf("%v/%v.txt", config.SavePath, filename), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *songSingleton) save() error {
	if lastSaved == *s {
		return nil
	}

	err := _writeSongFile("artist", s.Artist)
	if err != nil {
		return fmt.Errorf("couldn't save artist file. %v", err)
	}

	err = _writeSongFile("song", s.Name)
	if err != nil {
		return fmt.Errorf("couldn't save songSingleton file. %v", err)
	}

	err = _writeSongFile("entire", fmt.Sprintf("%v - %v", s.Name, s.Artist))
	if err != nil {
		return fmt.Errorf("couldn't save artist & songSingleton file. %v", err)
	}

	err = _writeSongFile("entire-descending", fmt.Sprintf("%v - %v", s.Artist, s.Name))
	if err != nil {
		return fmt.Errorf("couldn't save artist & songSingleton file. %v", err)
	}

	lastSaved = *s

	// Adds song to cache list.
	songCache = append(songCache, *s)

	return nil
}
