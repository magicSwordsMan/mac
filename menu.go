package mac

/*
#include "driver.h"
#include "menu.h"
*/
import "C"
import (
	"strconv"
	"time"
	"unsafe"

	"fmt"

	"os"

	"github.com/murlokswarm/app"
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
	"github.com/murlokswarm/uid"
)

type appMenu struct {
	*menu
}

func newAppMenu() *appMenu {
	return &appMenu{
		menu: newMenu(),
	}
}

func (m *appMenu) Mount(c markup.Componer) {
	ensureLaunched()

	m.menu.Mount(c)
	C.Driver_SetAppMenu(m.ptr)
}

type contextMenu struct {
	*menu
}

func newContextMenu() *contextMenu {
	m := &contextMenu{
		menu: newMenu(),
	}

	return m
}

func (m *contextMenu) Mount(c markup.Componer) {
	m.menu.Mount(c)
	C.Driver_ShowContextMenu(m.ptr)
}

type menu struct {
	id        uid.ID
	ptr       unsafe.Pointer
	component markup.Componer
}

func newMenu() *menu {
	id := uid.Context()

	cmenu := C.Menu__{
		ID: C.CString(id.String()),
	}

	defer free(unsafe.Pointer(cmenu.ID))

	menu := &menu{
		id:  id,
		ptr: C.Menu_New(cmenu),
	}

	app.RegisterContext(menu)
	return menu
}

func (m *menu) ID() uid.ID {
	return m.id
}

func (m *menu) Mount(c markup.Componer) {
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

func (m *menu) mount(elem *markup.Element) (err error) {
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

func (m *menu) mountContainer(elem *markup.Element) error {
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

func (m *menu) mountItem(elem *markup.Element) (err error) {
	var iconPath string

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

	if len(icon.Value) != 0 {
		iconPath = app.Resources().Join(icon.Value)

		if !app.IsSupportedImageExtension(iconPath) {
			err = fmt.Errorf("extension of %v is not supported", iconPath)
			return
		}

		if _, err = os.Stat(iconPath); err != nil {
			return
		}
	}

	item := C.MenuItem__{
		ID:        C.CString(elem.ID.String()),
		Label:     C.CString(label.Value),
		Icon:      C.CString(iconPath),
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
	return
}

func (m *menu) associate(parent *markup.Element, child *markup.Element) {
	parentID := C.CString(parent.ID.String())
	childID := C.CString(child.ID.String())

	defer free(unsafe.Pointer(parentID))
	defer free(unsafe.Pointer(childID))

	C.Menu_Associate(m.ptr, parentID, childID)
}

func (m *menu) Render(elem *markup.Element) {
	if err := m.mount(elem); err != nil {
		log.Error(err)
	}
}

func (m *menu) Position() (x float64, y float64) {
	return
}

func (m *menu) Move(x float64, y float64) {
}

func (m *menu) Size() (width float64, height float64) {
	return
}

func (m *menu) Resize(width float64, height float64) {
}

func (m *menu) SetIcon(path string) {
}

func (m *menu) SetBadge(v interface{}) {
}

func (m *menu) Close() {
}

//export onMenuItemClick
func onMenuItemClick(id *C.char, method *C.char) {
	if err := markup.Call(uid.ID(C.GoString(id)), C.GoString(method), ""); err != nil {
		log.Error(err)
	}
}

//export onMenuCloseFinal
func onMenuCloseFinal(cid *C.char) {
	ctx, err := app.ContextByID(uid.ID(C.GoString(cid)))
	if err != nil {
		return
	}

	menu := ctx.(*menu)

	go func() {
		time.Sleep(time.Millisecond * 42)
		markup.Dismount(menu.component)
		app.UnregisterContext(menu)
	}()
}
