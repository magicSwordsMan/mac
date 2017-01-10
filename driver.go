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
	"strings"
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/errors"
	"github.com/murlokswarm/log"
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
	share     app.Sharer
}

// NewDriver creates a new MacOS driver.
// It initializes the Cocoa app.
func NewDriver() *Driver {
	// runtime.LockOSThread()

	resources := app.ResourcePath("resources")
	if isAppPackaged() {
		cresources := C.Driver_Resources()
		resources = app.ResourcePath(C.GoString(cresources))
	}

	return &Driver{
		ptr:       C.Driver_Init(),
		resources: resources,
		appMenu:   newMenuBar(),
		dock:      newDock(),
		share:     &share{},
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

// Resources return the resources directory path.
func (d *Driver) Resources() app.ResourcePath {
	return d.resources
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

func isAppPackaged() (packaged bool) {
	execName := os.Args[0]

	path, err := filepath.Abs(filepath.Dir(execName))
	if err != nil {
		log.Error(errors.Newf("can't determine if app is packaged: %v", err))
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
