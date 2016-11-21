package mac

/*
#include "driver.h"
#include "menu.h"
*/
import "C"
import (
	"strconv"
	"unsafe"

	"fmt"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

type AppMenu struct {
	*Menu
}

func NewAppMenu() *AppMenu {
	return &AppMenu{
		Menu: NewMenu(),
	}
}

func (m *AppMenu) Mount(c markup.Componer) {
	m.Menu.Mount(c)
	C.Driver_SetAppMenu(m.ptr)
}

func (m *AppMenu) Close() {
	log.Error("application menu can't be closed")
}

type Menu struct {
	id        uid.ID
	ptr       unsafe.Pointer
	component markup.Componer
}

func NewMenu() *Menu {
	id := uid.Context()

	cmenu := C.Menu__{
		ID: C.CString(id.String()),
	}

	defer free(unsafe.Pointer(cmenu.ID))

	menu := &Menu{
		id:  id,
		ptr: C.Menu_New(cmenu),
	}

	app.RegisterContext(menu)
	return menu
}

func (m *Menu) ID() uid.ID {
	return m.id
}

func (m *Menu) Mount(c markup.Componer) {
	if m.component != nil {
		C.Menu_Clear(m.ptr)
		markup.Dismount(m.component)

	}

	m.component = c

	root, err := markup.Mount(c, m.ID())
	if err != nil {
		log.Panic(err)
	}

	if err := m.mount(root); err != nil {
		log.Panic(err)
	}

	rootID := C.CString(root.ID.String())
	defer free(unsafe.Pointer(rootID))

	C.Menu_Mount(m.ptr, rootID)
}

func (m *Menu) mount(elem *markup.Element) (err error) {
	switch elem.Name {
	case "menu":
		if err = m.mountContainer(elem); err != nil {
			return
		}

	case "menuitem":
		if err = m.mountItem(elem); err != nil {
			return
		}

	default:
		return fmt.Errorf("%v markup is not supported in a menu context. valid tags are menu and menuitem", elem)
	}

	for _, child := range elem.Children {
		if markup.IsComponentName(child.Name) {
			child, _ = markup.ComponentRoot(child.Component)
		}

		if err = m.mount(child); err != nil {
			return
		}

		m.associate(elem, child)
	}

	return
}

func (m *Menu) mountContainer(elem *markup.Element) error {
	if elem.Parent != nil && elem.Parent.Name != "menu" {
		return fmt.Errorf("%v can only have another menu as parent: %v", elem, elem.Parent)
	}

	label, _ := elem.Attributes.Attr("label")

	container := C.MenuContainer__{
		ID:    C.CString(elem.ID.String()),
		Label: C.CString(label.Value),
	}

	defer free(unsafe.Pointer(container.ID))
	defer free(unsafe.Pointer(container.Label))

	C.Menu_MountContainer(m.ptr, container)
	return nil
}

func (m *Menu) mountItem(elem *markup.Element) error {
	if elem.Parent == nil || elem.Parent.Name != "menu" {
		return fmt.Errorf("%v should have a menu as parent: %v", elem, elem.Parent)
	}

	label, _ := elem.Attributes.Attr("label")
	icon, _ := elem.Attributes.Attr("icon")
	shortcut, _ := elem.Attributes.Attr("shortcut")
	selector, _ := elem.Attributes.Attr("selector")
	onclick, _ := elem.Attributes.Attr("_onclick")
	disabled, _ := elem.Attributes.Attr("disabled")
	separator, _ := elem.Attributes.Attr("separator")

	isDisabled, _ := strconv.ParseBool(disabled.Value)
	isSeparator, _ := strconv.ParseBool(separator.Value)

	item := C.MenuItem__{
		ID:        C.CString(elem.ID.String()),
		Label:     C.CString(label.Value),
		Icon:      C.CString(icon.Value),
		Shortcut:  C.CString(shortcut.Value),
		Selector:  C.CString(selector.Value),
		OnClick:   C.CString(onclick.Value),
		Disabled:  boolToBOOL(isDisabled),
		Separator: boolToBOOL(isSeparator),
	}

	defer free(unsafe.Pointer(item.ID))
	defer free(unsafe.Pointer(item.Label))
	defer free(unsafe.Pointer(item.Icon))
	defer free(unsafe.Pointer(item.Shortcut))
	defer free(unsafe.Pointer(item.Selector))
	defer free(unsafe.Pointer(item.OnClick))

	C.Menu_MountItem(m.ptr, item)
	return nil
}

func (m *Menu) associate(parent *markup.Element, child *markup.Element) {
	parentID := C.CString(parent.ID.String())
	childID := C.CString(child.ID.String())

	defer free(unsafe.Pointer(parentID))
	defer free(unsafe.Pointer(childID))

	C.Menu_Associate(m.ptr, parentID, childID)
}

func (m *Menu) Render(elem *markup.Element) {
	if err := m.mount(elem); err != nil {
		log.Error(err)
	}
}

func (m *Menu) Position() (x float64, y float64) {
	return
}

func (m *Menu) Move(x float64, y float64) {
}

func (m *Menu) Size() (width float64, height float64) {
	return
}

func (m *Menu) Resize(width float64, height float64) {
}

func (m *Menu) SetIcon(path string) {
	return
}

func (m *Menu) Close() {
	C.Menu_Close(m.ptr)
	app.UnregisterContext(m)
}

//export onMenuItemClick
func onMenuItemClick(id *C.char, method *C.char) {
	markup.Call(uid.ID(C.GoString(id)), C.GoString(method), "")
}
