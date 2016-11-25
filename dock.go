package mac

/*
#include "driver.h"
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
)

type dock struct {
	*menu
}

func newDock() *dock {
	return &dock{
		menu: newMenu(),
	}
}

func (d *dock) Mount(c markup.Componer) {
	ensureLaunched()
	d.menu.Mount(c)
	C.Driver_SetDockMenu(d.ptr)
}

func (d *dock) SetIcon(path string) {
	ensureLaunched()

	cpath := C.CString(path)
	defer free(unsafe.Pointer(cpath))

	if len(path) == 0 {
		C.Driver_SetDockIcon(cpath)
		return
	}

	if !app.IsSupportedImageExtension(path) {
		log.Errorf("extension of %v is not supported", path)
		return
	}

	if _, err := os.Stat(path); err != nil {
		log.Error(err)
		return
	}

	C.Driver_SetDockIcon(cpath)
}

func (d *dock) SetBadge(v interface{}) {
	ensureLaunched()

	cv := C.CString(fmt.Sprint(v))
	defer free(unsafe.Pointer(cv))

	C.Driver_SetDockBadge(cv)
}
