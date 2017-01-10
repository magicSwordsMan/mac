package mac

/*
#include "driver.h"
*/
import "C"
import "github.com/murlokswarm/app"

type menuBar struct {
	*menu
}

func newMenuBar() *menuBar {
	return &menuBar{
		menu: newMenu(app.Menu{}),
	}
}

func (m *menuBar) Mount(c app.Componer) {
	ensureLaunched()
	m.menu.Mount(c)
	C.Driver_SetMenuBar(m.ptr)
}
