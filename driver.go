package mac

/*
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
)

// Driver is the implementation of the MacOS driver.
type Driver struct {
	ptr       unsafe.Pointer
	resources app.ResourcePath
	appMenu   app.Contexter
	dock      app.Contexter
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
		appMenu:   NewAppMenu(),
		dock:      NewDock(),
	}
}

// Run launches the Cocoa app.
func (d *Driver) Run() {
	C.Driver_Run()
}

func (d *Driver) NewContext(ctx interface{}) app.Contexter {
	switch c := ctx.(type) {
	case app.Window:
		return NewWindow(c)

	case app.ContextMenu:
		return NewContextMenu()

	default:
		return app.NewZeroContext(reflect.TypeOf(c).String())
	}
}

func (d *Driver) AppMenu() app.Contexter {
	return d.appMenu
}

func (d *Driver) Dock() app.Contexter {
	return d.dock
}

func (d *Driver) Resources() app.ResourcePath {
	return d.resources
}

func (d *Driver) JavascriptBridge() string {
	return "window.webkit.messageHandlers.Call.postMessage(msg);"
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
