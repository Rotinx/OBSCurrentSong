// Derived from https://github.com/sqweek/dialog

package dialogs

import (
	"errors"
	"github.com/TheTitanrain/w32"
	"syscall"
	"unsafe"
)

// ErrCancelled is an error returned when a user cancels/closes a dialog.
var ErrCancelled = errors.New("cancelled")

// Dlg is the common type for dialogs.
type Dlg struct {
	Title string
}

// DirectoryBuilder is used for directory browse dialogs.
type DirectoryBuilder struct {
	Dlg
	StartDir string
}

type dirdlg struct {
	bi *w32.BROWSEINFO
}

const (
	bffm_INITIALIZED   = 1
	bffm_SETSELECTIONW = w32.WM_USER + 103
	bffm_SETSELECTION  = bffm_SETSELECTIONW
)

// Directory initialises a DirectoryBuilder using the default configuration.
func Directory() *DirectoryBuilder {
	return &DirectoryBuilder{}
}

// Browse spawns the directory selection dialog using the configured settings,
// asking the user to select a single folder. Returns ErrCancelled as the error
// if the user cancels or closes the dialog.
func (b *DirectoryBuilder) Browse() (string, error) {
	return b.browse()
}

// Title specifies the title to be used for the dialog.
func (b *DirectoryBuilder) Title(title string) *DirectoryBuilder {
	b.Dlg.Title = title
	return b
}

// SetStartDir specifies the initial directory to be used for the dialog.
func (b *DirectoryBuilder) SetStartDir(dir string) *DirectoryBuilder {
	b.StartDir = dir
	return b
}

func callbackDefaultDir(hwnd w32.HWND, msg uint, lParam, lpData uintptr) int {
	if msg == bffm_INITIALIZED {
		_ = w32.SendMessage(hwnd, bffm_SETSELECTION, w32.TRUE, lpData)
	}
	return 0
}

func selectdir(b *DirectoryBuilder) (d dirdlg) {
	d.bi = &w32.BROWSEINFO{Flags: w32.BIF_RETURNONLYFSDIRS | w32.BIF_NEWDIALOGSTYLE}
	if b.Dlg.Title != "" {
		d.bi.Title, _ = syscall.UTF16PtrFromString(b.Dlg.Title)
	}
	if b.StartDir != "" {
		s16, _ := syscall.UTF16PtrFromString(b.StartDir)
		d.bi.LParam = uintptr(unsafe.Pointer(s16))
		d.bi.CallbackFunc = syscall.NewCallback(callbackDefaultDir)
	}
	return d
}

func (b *DirectoryBuilder) browse() (string, error) {
	d := selectdir(b)
	res := w32.SHBrowseForFolder(d.bi)
	if res == 0 {
		return "", ErrCancelled
	}
	return w32.SHGetPathFromIDList(res), nil
}
