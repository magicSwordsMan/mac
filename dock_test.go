package mac

import "testing"

func TestNewDock(t *testing.T) {
	newDock()
}

func TestDockMount(t *testing.T) {
	launched = true
	defer func() { launched = false }()

	d := newDock()
	c := &MenuComponent{}

	d.Mount(c)
}

func TestDockSetIcon(t *testing.T) {
	launched = true
	defer func() { launched = false }()

	d := newDock()

	// Set.
	d.SetIcon("resources/logo.png")

	// Unset.
	d.SetIcon("")

	// Bad extension.
	d.SetIcon("resources/logo.bmp")

	// Nonexistent.
	d.SetIcon("resources/logosh.png")
}

func TestDockSetBadge(t *testing.T) {
	launched = true
	defer func() { launched = false }()

	d := newDock()
	d.SetBadge(42)
}
