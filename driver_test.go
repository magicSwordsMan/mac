package mac

import (
	"net/url"
	"sync"
	"testing"

	"github.com/murlokswarm/app"
)

func TestDriver(t *testing.T) {
	t.Log(driver.MenuBar())
	t.Log(driver.Dock())
	t.Log(driver.Storage())
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

func TestOnFileOpen(t *testing.T) {
	app.OnFileOpen = func(n string) {
		t.Log("MacOS driver onFileOpen:", n)
		if n != "zune" {
			t.Error("n should be zune")
		}
	}
	onFileOpen(cString("zune"))
}

func TestOnFilesOpen(t *testing.T) {
	app.OnFilesOpen = func(filenames []string) {
		t.Log("MacOS driver onFilesOpen:", filenames)
		if filenames[0] != "zune" {
			t.Error("filenames[0] should be zune")
		}
		if filenames[1] != "mune" {
			t.Error("filenames[1] should be mune")
		}
	}
	onFilesOpen(cString(`["zune", "mune"]`))
}

func TestOnURLOpen(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	app.OnURLOpen = func(URL url.URL) {
		t.Log("MacOS driver onURLOpen:", URL)
		wg.Done()
	}
	onURLOpen(cString("github-mac://openRepo/https://github.com/murlokswarm/app"))

	wg.Wait()

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
