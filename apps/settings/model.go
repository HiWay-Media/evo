package settings

import "github.com/hiway-media/evo"

type Settings struct {
	evo.Model
	Reference string
	Title     string
	Data      string
	Default   string
	Ptr       interface{} `gorm:"-"`
}
