package database

import (
	"database/sql"
	"fmt"
	"go-restapi/app"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDB(conf app.Config) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Database.Username, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Database)
	sqlDB, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	gormConfig := &gorm.Config{}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), gormConfig)
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
