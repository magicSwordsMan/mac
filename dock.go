package mac

/*
#include "driver.h"
*/
import "C"
import (
	"unsafe"

	"os"

	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
)

type Dock struct {
	*Menu
}

func NewDock() *Dock {
	return &Dock{
		Menu: NewMenu(),
	}
}

func (d *Dock) Mount(c markup.Componer) {
	d.Menu.Mount(c)
	C.Driver_SetDockMenu(d.ptr)
}

func (d *Dock) SetIcon(path string) {
	if _, err := os.Stat(path); len(path) != 0 && err != nil {
		log.Error(err)
		return
	}

	cpath := C.CString(path)
	defer free(unsafe.Pointer(cpath))

	C.Driver_SetDockIcon(cpath)
}

func (d *Dock) SetBadge(v string) {
	cv := C.CString(v)
	defer free(unsafe.Pointer(cv))

	C.Driver_SetDockBadge(cv)
}

func (d *Dock) Close() {
	log.Error("dock can't be closed")
}
