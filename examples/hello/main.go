package main

import (
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	_ "github.com/murlokswarm/mac"
	"github.com/murlokswarm/markup"
)

var (
	win app.Contexter
)

func init() {
	markup.RegisterComponent("Hello", func() markup.Componer {
		return &Hello{}
	})
}

type Hello struct {
	Greeting string
}

func (h *Hello) Render() string {
	return `
<div>
    Hello, {{if .Greeting}}{{.Greeting}}{{else}}World{{end}}
    <input type="text" placeholder="What is your name?" />
</div>
    `
}

func main() {
	app.OnLaunch = func() {
		win = newWindow()

		hello := &Hello{}
		win.Mount(hello)

		go func() {
			name := []string{
				"m",
				"ma",
				"max",
				"maxe",
				"maxen",
				"maxenc",
				"maxence",
			}

			time.Sleep(time.Second)

			for _, s := range name {
				time.Sleep(time.Millisecond * 15)
				hello.Greeting = s
				app.Render(hello)
			}

			// win.Move(300, 300)
			// win.Resize(42, 42)
			w, h := win.Size()
			log.Infof("win size: %vx%v", w, h)

			x, y := win.Position()
			log.Infof("win pos: (%v, %v)", x, y)

			// win.Close()

		}()
	}

	app.OnReopen = func(hasVisibleWindow bool) {
		if win != nil {
			return
		}

		win = newWindow()
		hello := &Hello{}
		win.Mount(hello)
	}

	app.Run()
}

func newWindow() app.Contexter {
	return app.NewWindow(app.Window{
		Width:          1340,
		Height:         720,
		Vibrancy:       app.VibeMediumLight,
		TitlebarHidden: true,
		Title:          "main",

		OnMinimize:       func() { log.Info("OnMinimize") },
		OnDeminimize:     func() { log.Info("OnDeminimize") },
		OnFullScreen:     func() { log.Info("OnFullScreen") },
		OnExitFullScreen: func() { log.Info("OnExitFullScreen") },
		OnMove:           func(x float64, y float64) { log.Infof("OnMove (%v, %v)", x, y) },
		OnResize:         func(w float64, h float64) { log.Infof("OnResize %vx%v", w, h) },
		OnFocus:          func() { log.Info("OnFocus") },
		OnBlur:           func() { log.Info("OnBlur") },
		OnClose: func() bool {
			log.Info("OnClose")
			win = nil
			return true
		},
	})
}