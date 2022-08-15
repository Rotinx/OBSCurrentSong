package main

import (
	"OBSCurrentSong/ps"
	"errors"
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

type Song struct {
	Name   string
	Artist string
}

var NoSongPlaying = errors.New("no song playing")
var NoTitleFound = errors.New("title couldn't be found")

func currentSong() (*Song, error) {
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

	song := &Song{
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
