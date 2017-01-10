package mac

/*
#include "share.h"
*/
import "C"
import (
	"net/url"
	"unsafe"
)

type share struct {
}

func (s *share) Text(v string) {
	cv := cString(v)
	defer free(unsafe.Pointer(cv))

	C.Share_Text(cv)
}

func (s *share) URL(v *url.URL) {
	cv := cString(v.String())
	defer free(unsafe.Pointer(cv))

	C.Share_URL(cv)
}
