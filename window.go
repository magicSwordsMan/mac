package mac

/*
#include "window.h"
*/
import "C"
import (
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

type Window struct {
	id  uid.ID
	ptr unsafe.Pointer
}

func NewWindow(w app.Window) *Window {
	id := uid.Context()

	cwin := C.Window__{
		ID:              C.CString(id.String()),
		Title:           C.CString(w.Title),
		X:               C.CGFloat(w.X),
		Y:               C.CGFloat(w.Y),
		Width:           C.CGFloat(w.Width),
		Height:          C.CGFloat(w.Height),
		BackgroundColor: C.CString(w.BackgroundColor),
		Vibrancy:        C.NSVisualEffectMaterial(w.Vibrancy),
		Borderless:      boolToBOOL(w.Borderless),
		FixedSize:       boolToBOOL(w.FixedSize),
		CloseHidden:     boolToBOOL(w.CloseHidden),
		MinimizeHidden:  boolToBOOL(w.MinimizeHidden),
		TitlebarHidden:  boolToBOOL(w.TitlebarHidden),
	}

	defer free(unsafe.Pointer(cwin.ID))
	defer free(unsafe.Pointer(cwin.Title))
	defer free(unsafe.Pointer(cwin.BackgroundColor))

	return &Window{
		id:  id,
		ptr: C.Window_New(cwin),
	}
}

func (w *Window) ID() uid.ID {
	return w.id
}

func (w *Window) Mount(c markup.Componer) {
	// call objc for inserting html
}

func (w *Window) Move(x float64, y float64) {
	// call method to move window
}

func (w *Window) Resize(width float64, height float64) {
	// call method to Resize window
}

func (w *Window) SetIcon(path string) {
	return
}
