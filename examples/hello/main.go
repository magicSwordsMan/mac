package main

import (
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	_ "github.com/murlokswarm/mac"
)

var (
	win app.Contexter
)

func init() {
	app.RegisterComponent(&Hello{})
}

type Hello struct {
	Greeting string
}

func (h *Hello) Render() string {
	return `
<div>
    Hello,
	<span>{{if .Greeting}}{{html .Greeting}}{{else}}World{{end}}</span>
    <input type="text" 
		   placeholder="What is your name?" 
		   _onchange="OnInputChange" 
		   _onkeyup="OnKeyUp"
		   _onmousewheel="OnWheel"
		   _onclick="OnClick"
		   _oncontextmenu="OnContextMenu"
		   value="{{html .Greeting}}"  
		   autofocus="true" />
</div>
    `
}

func (h *Hello) OnInputChange(e app.ChangeArg) {
	h.Greeting = e.Value
	app.Render(h)
}

func (h *Hello) OnKeyUp(e app.KeyboardArg) {
	log.Infof("%+v", e)
}

func (h *Hello) OnWheel(e app.WheelArg) {
	log.Infof("%+v", e)
}

func (h *Hello) OnClick(e app.MouseArg) {
	log.Infof("%+v", e)
}

func (h *Hello) OnContextMenu() {
	m := app.NewContextMenu()
	m.Mount(&AppMainMenu{Sep: true})
}

func main() {
	app.OnLaunch = func() {
		menu := &AppMainMenu{Sep: true}
		app.Menu().Mount(menu)

		dock := &AppMainMenu{}
		app.Dock().Mount(dock)
		app.Dock().SetBadge(42)
		app.Dock().SetIcon(app.Resources().Join("contexticon.png"))

		win = newWindow()

		hello := &Hello{}
		win.Mount(hello)

		go func() {
			// name := []string{
			// 	"m",
			// 	"ma",
			// 	"max",
			// 	"maxe",
			// 	"maxen",
			// 	"maxenc",
			// 	"maxence",
			// }

			// time.Sleep(time.Second * 3)
			// menu.CustomTitle = "La vie est belle"
			// menu.Sep = false
			// menu.Disabled = true
			app.Render(menu)
			// app.Dock().SetBadge("Hello Ach")
			// app.Dock().SetIcon(app.Resources().Join("dsffa.png"))

			ctxm := app.NewContextMenu()
			ctxm.Mount(&AppMainMenu{})

			time.Sleep(time.Second * 3)
			ctxm.Close()
			// app.Dock().SetIcon("")

			// for _, s := range name {
			// 	time.Sleep(time.Millisecond * 15)
			// 	hello.Greeting = s
			// 	app.Render(hello)
			// }

			// win.Move(300, 300)
			// win.Resize(42, 42)
			// w, h := win.Size()
			// log.Infof("win size: %vx%v", w, h)

			// x, y := win.Position()
			// log.Infof("win pos: (%v, %v)", x, y)

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
		Width:    1340,
		Height:   720,
		Vibrancy: app.VibeDark,
		// TitlebarHidden: true,
		Title: "main",

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
