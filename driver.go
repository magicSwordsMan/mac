// Package mac implements the macOS driver.
// Usage:
// import _ "github.com/murlokswarm/mac"
// During initialization, the package calls yui.RegisterDriver() with its
// Driver implementation.
package mac

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework WebKit -framework CoreImage
#include "driver.h"
*/
import "C"
import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

var (
	driver   = NewDriver()
	launched = false
)

func init() {
	app.RegisterDriver(driver)
}

// Driver is the implementation of the MacOS driver.
type Driver struct {
	ptr       unsafe.Pointer
	resources app.ResourcePath
	appMenu   app.Contexter
	dock      app.Docker
}

// NewDriver creates a new MacOS driver.
// It initializes the Cocoa app.
func NewDriver() *Driver {
	resources := app.ResourcePath("resources")
	if isAppPackaged() {
		cresources := C.Driver_Resources()
		resources = app.ResourcePath(C.GoString(cresources))
	}

	return &Driver{
		ptr:       C.Driver_Init(),
		resources: resources,
		appMenu:   newAppMenu(),
		dock:      newDock(),
	}
}

// Run launches the Cocoa app.
func (d *Driver) Run() {
	C.Driver_Run()
}

// NewContext creates a new context.
func (d *Driver) NewContext(ctx interface{}) app.Contexter {
	ensureLaunched()

	switch c := ctx.(type) {
	case app.Window:
		return newWindow(c)

	case app.ContextMenu:
		return newContextMenu()

	default:
		return app.NewZeroContext(reflect.TypeOf(c).String())
	}
}

// AppMenu returns the application menu.
func (d *Driver) AppMenu() app.Contexter {
	return d.appMenu
}

// Dock returns the dock.
func (d *Driver) Dock() app.Docker {
	return d.dock
}

// Resources return the resources directory path.
func (d *Driver) Resources() app.ResourcePath {
	return d.resources
}

// JavascriptBridge returns the javascript statement to allow javascript to
// call go component methods.
func (d *Driver) JavascriptBridge() string {
	return "window.webkit.messageHandlers.Call.postMessage(msg);"
}

func (d *Driver) terminate() {
	C.Driver_Terminate()
}

func isAppPackaged() (packaged bool) {
	execName := os.Args[0]

	path, err := filepath.Abs(filepath.Dir(execName))
	if err != nil {
		log.Errorf("can't determine if app is packaged: %v", err)
		return
	}

	for _, dir := range strings.Split(path, "/") {
		if strings.HasSuffix(dir, ".app") {
			return true
		}
	}
	return
}

func ensureLaunched() {
	if !launched {
		log.Panic(`creating and interacting with contexts requires the app to be launched. set app.OnLaunch handler and launch the app by calling app.Run()`)
	}
}

//export onLaunch
func onLaunch() {
	if app.OnLaunch != nil {
		launched = true
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
	markup.Call(
		uid.ID(string(C.GoString(id))),
		C.GoString(name),
		C.GoString(jsonArg))
}
