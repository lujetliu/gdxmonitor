package model

import "gdxmonitor/global"

func Migrate() {
	global.DBEngine.AutoMigrate(&Block{})
}
