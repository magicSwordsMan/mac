package main

import (
	"github.com/murlokswarm/app"
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

		win.Mount(&Hello{})
	}

	app.Run()
}
