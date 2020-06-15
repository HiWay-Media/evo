package rdb

import (
	"fmt"
	"github.com/getevo/evo"
	"github.com/getevo/evo/lib/concurrent"
	"github.com/getevo/evo/lib/fontawesome"
	"github.com/getevo/evo/menu"
	"github.com/getevo/evo/user"
	"github.com/jinzhu/gorm"
)

var connections = concurrent.Map{}
var db *gorm.DB

// App settings app struct
type App struct{}

var registered = true

// Register register the rdb in io apps
func Register() {
	if !registered {
		return
	}
	registered = false
	evo.Register(App{})
}

// Register settings app
func (App) Register() {
	fmt.Println("RDB Registered")
	connections.Init()
	db = evo.GetDBO()
	PushDB("default", db.DB())
}

// Router setup routers
func (App) Router() {

}

// Permissions setup permissions of app
func (App) Permissions() []user.Permission {
	return []user.Permission{
		{Title: "Access Settings", CodeName: "view", Description: "Access list to view list of settings"},
		{Title: "Modify Settings", CodeName: "modify", Description: "Modify Settings"},
	}
}

// Menus setup menus
func (App) Menus() []menu.Menu {
	return []menu.Menu{
		{Title: "Settings", Url: "admin/settings", Permission: "settings.view", Icon: fontawesome.Cog},
	}
}

// WhenReady called after setup all apps
func (App) WhenReady() {}

func (App) Pack() {}
