package main

import (
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
)

func init() {
	markup.RegisterComponent("AppMainMenu", func() markup.Componer { return &AppMainMenu{} })
	markup.RegisterComponent("WindowMenu", func() markup.Componer { return &WindowMenu{} })
}

type AppMainMenu struct {
	CustomTitle string
	Sep         bool
	Disabled    bool
}

func (m *AppMainMenu) Render() string {
	return `
<menu>
    <menu label="app">
        <menuitem label="About" selector="orderFrontStandardAboutPanel:"/>
        <menuitem label="Quit" shortcut="meta+q" selector="terminate:" separator="{{.Sep}}" />
        <menuitem label="{{if .CustomTitle}}{{.CustomTitle}}{{else}}Booooooo{{end}}" 
                  shortcut="ctrl+c" 
                  _onclick="OnCustomMenuClick" 
                  icon="contexticon.png"
                  disabled="{{.Disabled}}" />
    </menu>
    <WindowMenu />
</menu>
    `
}

func (m *AppMainMenu) OnCustomMenuClick() {
	log.Info("OnCustomMenuClick")
}

type WindowMenu struct {
	Placeholder bool
}

func (m *WindowMenu) Render() string {
	return `
<menu label="Window">
    <menuitem label="Close" selector="performClose:" shortcut="meta+w" />
</menu>
    `
}
