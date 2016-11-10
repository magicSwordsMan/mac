package mac

/*
#include "window.h"
*/
import "C"
import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

type Window struct {
	id   uid.ID
	ptr  unsafe.Pointer
	root markup.Componer
}

// NewWindow creates a window.
func NewWindow(w app.Window) *Window {
	id := uid.Context()

	htmlCtx := app.HTMLContext{
		ID:       id,
		Title:    w.Title,
		Lang:     w.Lang,
		MurlokJS: app.MurlokJS(),
		JS:       app.Resources().JS(),
		CSS:      app.Resources().CSS(),
	}

	log.Info(htmlCtx.HTML())

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
		HTML:            C.CString(htmlCtx.HTML()),
		ResourcePath:    C.CString(app.Resources().Path()),
	}

	defer free(unsafe.Pointer(cwin.ID))
	defer free(unsafe.Pointer(cwin.Title))
	defer free(unsafe.Pointer(cwin.BackgroundColor))
	defer free(unsafe.Pointer(cwin.HTML))
	defer free(unsafe.Pointer(cwin.ResourcePath))

	win := &Window{
		id:  id,
		ptr: C.Window_New(cwin),
	}

	app.RegisterContext(win)
	return win
}

// ID the identifier of the window.
func (w *Window) ID() uid.ID {
	return w.id
}

// Mount mounts a component.
func (w *Window) Mount(c markup.Componer) {
	var html string
	var err error

	if err = markup.Mount(c, w.ID()); err != nil {
		log.Panic(err)
		return
	}

	if html, err = markup.ComponentToHTML(c); err != nil {
		log.Panic(err)
		return
	}

	w.root = c

	html = strconv.Quote(html)
	call := fmt.Sprintf(`Mount("%v", %v)`, w.ID(), html)

	ccall := C.CString(call)
	defer free(unsafe.Pointer(ccall))

	C.Window_CallJS(w.ptr, ccall)
}

func (w *Window) Render(elem *markup.Element) {
	html := strconv.Quote(elem.HTML())
	call := fmt.Sprintf(`Render("%v", %v)`, elem.ID, html)

	ccall := C.CString(call)
	defer free(unsafe.Pointer(ccall))

	C.Window_CallJS(w.ptr, ccall)
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

func (w *Window) Close() {
	markup.Dismount(w.root)
	app.UnregisterContext(w)
}
