package main

import (
	"time"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	_ "github.com/murlokswarm/mac"
	"github.com/murlokswarm/markup"
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
		win := app.NewWindow(app.Window{
			Width:          1340,
			Height:         720,
			Vibrancy:       app.VibeMediumLight,
			TitlebarHidden: true,
		})

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

			win.Move(300, 300)
			// win.Resize(42, 42)
			w, h := win.Size()
			log.Infof("win size: %vx%v", w, h)

			x, y := win.Position()
			log.Infof("win pos: (%v, %v)", x, y)

		}()
	}

	app.Run()
}
