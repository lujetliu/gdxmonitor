package model

import (
	"fmt"
	"gdxmonitor/global"
	"gdxmonitor/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		db.LogMode(true)
	}

	// db.SingularTable(true) // 表名称是单数
	// db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns) // 连接池里最大的连接数, TODO: Mysql 参数
	// db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns) // 最大的连接数, TODO: Mysql 参数

	return db, nil
}
