package mac

/*
#include "driver.h"
*/
import "C"
import (
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

func (d *Dock) Close() {
	log.Error("dock can't be closed")
}
