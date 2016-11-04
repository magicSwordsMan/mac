package mac

import (
	"testing"

	"github.com/murlokswarm/app"
)

func TestDriverRun(t *testing.T) {
	app.OnLaunch = func() {
		app.NewWindow(app.Window{
			Width:          1340,
			Height:         720,
			TitlebarHidden: true,
		})
	}

	app.Run()
}
