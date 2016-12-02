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

type window struct {
	id        uid.ID
	ptr       unsafe.Pointer
	component app.Componer
	config    app.Window
}

func newWindow(w app.Window) *window {
	id := uid.Context()

	htmlCtx := app.HTMLContext{
		ID:       id,
		Title:    w.Title,
		Lang:     w.Lang,
		MurlokJS: app.MurlokJS(),
		JS:       app.Resources().JS(),
		CSS:      app.Resources().CSS(),
	}

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

	win := &window{
		id:     id,
		ptr:    C.Window_New(cwin),
		config: w,
	}

	app.RegisterContext(win)
	return win
}

func (w *window) ID() uid.ID {
	return w.id
}

func (w *window) Mount(c app.Componer) {
	var html string
	var err error

	if w.component != nil {
		markup.Dismount(w.component)
	}

	w.component = c

	if _, err = markup.Mount(c, w.ID()); err != nil {
		log.Panic(err)
	}

	if html, err = markup.ComponentToHTML(c); err != nil {
		log.Panic(err)
	}

	html = strconv.Quote(html)
	call := fmt.Sprintf(`Mount("%v", %v)`, w.ID(), html)

	ccall := C.CString(call)
	defer free(unsafe.Pointer(ccall))

	C.Window_CallJS(w.ptr, ccall)
}

func (w *window) Render(elem *markup.Element) {
	html := strconv.Quote(elem.HTML())
	call := fmt.Sprintf(`Render("%v", %v)`, elem.ID, html)

	ccall := C.CString(call)
	defer free(unsafe.Pointer(ccall))

	C.Window_CallJS(w.ptr, ccall)
}

func (w *window) Position() (x float64, y float64) {
	frame := C.Window_Frame(w.ptr)
	x = float64(frame.origin.x)
	y = float64(frame.origin.y)
	return
}

func (w *window) Move(x float64, y float64) {
	C.Window_Move(w.ptr, C.CGFloat(x), C.CGFloat(y))
}

func (w *window) Size() (width float64, height float64) {
	frame := C.Window_Frame(w.ptr)
	width = float64(frame.size.width)
	height = float64(frame.size.height)
	return
}

func (w *window) Resize(width float64, height float64) {
	C.Window_Resize(w.ptr, C.CGFloat(width), C.CGFloat(height))
}

func (w *window) SetIcon(path string) {
	return
}

func (w *window) SetBadge(v interface{}) {
	return
}

func (w *window) Close() {
	C.Window_Close(w.ptr)
}

//export onWindowMinimize
func onWindowMinimize(cid *C.char) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)

	if win.config.OnMinimize != nil {
		win.config.OnMinimize()
	}
}

//export onWindowDeminimize
func onWindowDeminimize(cid *C.char) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)

	if win.config.OnDeminimize != nil {
		win.config.OnDeminimize()
	}
}

//export onWindowFullScreen
func onWindowFullScreen(cid *C.char) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)

	if win.config.OnFullScreen != nil {
		win.config.OnFullScreen()
	}
}

//export onWindowExitFullScreen
func onWindowExitFullScreen(cid *C.char) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)

	if win.config.OnExitFullScreen != nil {
		win.config.OnExitFullScreen()
	}
}

//export onWindowMove
func onWindowMove(cid *C.char, x C.CGFloat, y C.CGFloat) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)

	if win.config.OnMove != nil {
		win.config.OnMove(float64(x), float64(y))
	}
}

//export onWindowResize
func onWindowResize(cid *C.char, width C.CGFloat, height C.CGFloat) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)

	if win.config.OnResize != nil {
		win.config.OnResize(float64(width), float64(height))
	}
}

//export onWindowFocus
func onWindowFocus(cid *C.char) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)

	if win.config.OnFocus != nil {
		win.config.OnFocus()
	}
}

//export onWindowBlur
func onWindowBlur(cid *C.char) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)

	if win.config.OnBlur != nil {
		win.config.OnBlur()
	}
}

//export onWindowClose
func onWindowClose(cid *C.char) bool {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return true
	}

	win := ctx.(*window)

	if win.config.OnClose != nil {
		return win.config.OnClose()
	}
	return true
}

//export onWindowCloseFinal
func onWindowCloseFinal(cid *C.char) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	win := ctx.(*window)
	markup.Dismount(win.component)
	app.UnregisterContext(win)
}
