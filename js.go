package mac

import "C"
import (
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
)

//export onJSCall
func onJSCall(msg *C.char) {
	app.CallComponentMethod(C.GoString(msg))
}

//export onJSAlert
func onJSAlert(alert *C.char) {
	log.Warn(C.GoString(alert))
}
