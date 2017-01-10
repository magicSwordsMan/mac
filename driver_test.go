package mac

import (
	"testing"

	"github.com/murlokswarm/app"
)

func TestDriver(t *testing.T) {
	t.Log(driver.MenuBar())
	t.Log(driver.Dock())
	t.Log(driver.Resources())
	t.Log(driver.JavascriptBridge())
}

func TestDriverNewContext(t *testing.T) {
	launched = true
	defer func() { launched = false }()

	// Window => would bloc.
	// driver.NewContext(app.Window{})

	// Menu.
	driver.NewContext(app.ContextMenu{})
}

func TestDriverNewContextNotImplemented(t *testing.T) {
	defer func() { recover() }()

	launched = true
	defer func() { launched = false }()

	driver.NewContext("not implement")
}

func TestDriverNewContextPanic(t *testing.T) {
	defer func() { recover() }()

	driver.NewContext(app.Window{})
	t.Error("should panic")

}

func TestOnLaunch(t *testing.T) {
	app.OnLaunch = func() {
		t.Log("MacOS driver onLaunch")
	}

	onLaunch()
}

func TestFocused(t *testing.T) {
	app.OnFocus = func() {
		t.Log("MacOS driver onFocus")
	}

	onFocus()
}

func TestOnBlur(t *testing.T) {
	app.OnBlur = func() {
		t.Log("MacOS driver onBlur")
	}

	onBlur()
}

func TestOnReopen(t *testing.T) {
	app.OnReopen = func(v bool) {
		t.Log("MacOS driver onReopen:", v)

		if !v {
			t.Error("v should be true")
		}
	}

	onReopen(true)
}

func TestFileOpened(t *testing.T) {
	app.OnFileOpen = func(n string) {
		t.Log("MacOS driver fileOpened:", n)
		if n != "zune" {
			t.Error("n should be zune")
		}
	}

	onFileOpen(cString("zune"))
}

func TestOnTerminate(t *testing.T) {
	app.OnTerminate = func() bool {
		t.Log("MacOS driver onTerminate")
		return false
	}

	if ret := onTerminate(); ret {
		t.Error("ret should be false")
	}

	app.OnTerminate = nil

	if ret := onTerminate(); !ret {
		t.Error("ret should be true")
	}
}

func TestOnFinalize(t *testing.T) {
	app.OnFinalize = func() {
		t.Log("MacOS driver onFinalize")
	}

	onFinalize()
}
