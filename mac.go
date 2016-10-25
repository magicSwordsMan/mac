// Package mac implements the macOS driver.
// Usage:
// import _ "github.com/murlokswarm/mac"
// During initialization, the package calls yui.RegisterDriver() with its
// Driver implementation.
package mac

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
*/
import "C"
import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

var (
	driver = NewDriver()
)

func init() {
	app.RegisterDriver(driver)
}

//export onLaunch
func onLaunch() {
	if app.OnLaunch != nil {
		app.OnLaunch()
	}
}

//export onFocus
func onFocus() {
	if app.OnFocus != nil {
		app.OnFocus()
	}
}

//export onBlur
func onBlur() {
	if app.OnBlur != nil {
		app.OnBlur()
	}
}

//export onReopen
func onReopen(hasVisibleWindow bool) {
	if app.OnReopen != nil {
		app.OnReopen(hasVisibleWindow)
	}
}

//export onFileOpen
func onFileOpen(filename *C.char) {
	if app.OnFileOpen != nil {
		app.OnFileOpen(C.GoString(filename))
	}
}

//export onTerminate
func onTerminate() bool {
	if app.OnTerminate != nil {
		return app.OnTerminate()
	}

	return true
}

//export onFinalize
func onFinalize() {
	if app.OnFinalize != nil {
		app.OnFinalize()
	}
}

//export onEvent
func onEvent(id *C.char, name *C.char, jsonArg *C.char) {
	err := markup.Call(uid.ID(
		string(C.GoString(id))),
		C.GoString(name),
		C.GoString(jsonArg))

	if err != nil {
		log.Error(err)
	}
}
