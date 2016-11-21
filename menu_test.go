package mac

import (
	"testing"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/markup"
)

type MenuComponent struct {
	Greeting                  string
	ErrorInvalidTag           bool
	ErrorCompositionContainer bool
	ErrorCompositionItem      bool
}

func (m *MenuComponent) Render() string {
	return `
<menu>
	<menuitem label="hello" separator="true" />
	<menuitem label="{{if .Greeting}}{{.Greeting}}{{else}}world{{end}}" />
	<SubMenuComponent />

	{{if .ErrorInvalidTag}}
		<div>Pouette</div>
	{{end}}

	{{if .ErrorCompositionContainer}}
		<menuitem>
			<menu></menu>
		</menuitem>
	{{end}}

	{{if .ErrorCompositionItem}}
		<menuitem>
			<menuitem></menuitem>
		</menuitem>
	{{end}}
</menu>
	`
}

type SubMenuComponent struct {
	Placeholder bool
}

func (m *SubMenuComponent) Render() string {
	return `
<menu>
	<menuitem label="foo" />
	<menuitem label="bar" />
</menu>
	`
}

func init() {
	markup.RegisterComponent("MenuComponent", func() markup.Componer { return &MenuComponent{} })
	markup.RegisterComponent("SubMenuComponent", func() markup.Componer { return &SubMenuComponent{} })
}

func TestMenu(t *testing.T) {
	m := NewMenu()
	defer m.Close()

	m.Position()
	m.Move(42, 42)
	m.Size()
	m.Resize(42, 42)
	m.SetIcon("")
}

func TestMenuMount(t *testing.T) {
	m := NewMenu()
	defer m.Close()

	c := &MenuComponent{}
	m.Mount(c)

	c2 := &MenuComponent{Greeting: "Maxoo"}
	m.Mount(c2)
}

func TestMenuMountInvalidTag(t *testing.T) {
	defer func() { recover() }()

	m := NewMenu()
	defer m.Close()

	c := &MenuComponent{ErrorInvalidTag: true}
	m.Mount(c)
	t.Error("should panic")
}

func TestMenuMountErrorCompositionContainer(t *testing.T) {
	defer func() { recover() }()

	m := NewMenu()
	defer m.Close()

	c := &MenuComponent{ErrorCompositionContainer: true}
	m.Mount(c)
	t.Error("should panic")
}

func TestMenuMountErrorCompositionItem(t *testing.T) {
	defer func() { recover() }()

	m := NewMenu()
	defer m.Close()

	c := &MenuComponent{ErrorCompositionItem: true}
	m.Mount(c)
	t.Error("should panic")
}

func TestMenuRender(t *testing.T) {
	m := NewMenu()
	// defer m.Close()

	c := &MenuComponent{}
	m.Mount(c)

	c.Greeting = "Maxence"
	app.Render(c)
}
