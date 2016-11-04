package mac

/*
#include "driver.h"
*/
import "C"
import (
	"unsafe"

	"reflect"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/uid"
)

// Driver is the implementation of the MacOS driver.
type Driver struct {
	delegatePtr unsafe.Pointer
}

// NewDriver creates a new MacOS driver.
// It initializes the Cocoa app.
func NewDriver() *Driver {
	return &Driver{
		delegatePtr: C.Driver_Init(),
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

	default:
		return app.NewZeroContext(reflect.TypeOf(c).String())
	}
}

func (d *Driver) Render(target uid.ID, HTML string) (err error) {
	return
}

func (d *Driver) AppMenu() app.Contexter {
	return nil
}

func (d *Driver) Dock() app.Contexter {
	return nil
}
