// Package mac implements the macOS driver.
// Usage:
// import _ "github.com/murlokswarm/mac"
// During initialization, the package calls yui.RegisterDriver() with its
// Driver implementation.
package mac

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa
#cgo LDFLAGS: -framework WebKit
#cgo LDFLAGS: -framework CoreImage
#cgo LDFLAGS: -framework Security
#include "driver.h"
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/errors"
	"github.com/murlokswarm/log"
)

var (
	driver   *Driver
	launched = false
)

func init() {
	runtime.LockOSThread()
	driver = NewDriver()
	app.RegisterDriver(driver)
}

// Driver is the implementation of the MacOS driver.
type Driver struct {
	ptr     unsafe.Pointer
	storage storage
	appMenu app.Contexter
	dock    app.Docker
	share   app.Sharer
}

// NewDriver creates a new MacOS driver.
// It initializes the Cocoa app.
func NewDriver() *Driver {
	return &Driver{
		ptr:     C.Driver_Init(),
		appMenu: newMenuBar(),
		dock:    newDock(),
		share:   &share{},
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
		return newContextMenu(c)

	default:
		log.Panicf("ctx for %T is not implemented", ctx)
		return nil
	}
}

// MenuBar returns the menu bar.
func (d *Driver) MenuBar() app.Contexter {
	return d.appMenu
}

// Dock returns the dock.
func (d *Driver) Dock() app.Docker {
	return d.dock
}

// Storage returns the directories location to use during app lifecycle.
func (d *Driver) Storage() app.Storer {
	return d.storage
}

// JavascriptBridge returns the javascript statement to allow javascript to
// call go component methods.
func (d *Driver) JavascriptBridge() string {
	return "window.webkit.messageHandlers.Call.postMessage(msg);"
}

// Share returns a sharing service.
func (d *Driver) Share() app.Sharer {
	return d.share
}

func (d *Driver) terminate() {
	C.Driver_Terminate()
}

func ensureLaunched() {
	if !launched {
		log.Panic(errors.New(`creating and interacting with contexts requires the app to be launched. set app.OnLaunch handler and launch the app by calling app.Run()`))
	}
}

//export onLaunch
func onLaunch() {
	launched = true

	app.UIChan <- func() {
		if app.OnLaunch != nil {
			app.OnLaunch()
		}
	}
}

//export onFocus
func onFocus() {
	app.UIChan <- func() {
		if app.OnFocus != nil {
			app.OnFocus()
		}
	}
}

//export onBlur
func onBlur() {
	app.UIChan <- func() {
		if app.OnBlur != nil {
			app.OnBlur()
		}
	}
}

//export onReopen
func onReopen(hasVisibleWindow bool) {
	app.UIChan <- func() {
		if app.OnReopen != nil {
			app.OnReopen(hasVisibleWindow)
		}
	}
}

//export onFileOpen
func onFileOpen(cfilename *C.char) {
	filename := C.GoString(cfilename)

	app.UIChan <- func() {
		if app.OnFileOpen != nil {
			app.OnFileOpen(filename)
		}
	}
}

//export onTerminate
func onTerminate() bool {
	termChan := make(chan bool)

	app.UIChan <- func() {
		if app.OnTerminate != nil {
			termChan <- app.OnTerminate()
			return
		}

		termChan <- true
	}
	return <-termChan
}

//export onFinalize
func onFinalize() {
	if app.OnFinalize != nil {
		app.OnFinalize()
	}
}
