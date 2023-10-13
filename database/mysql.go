package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDB(username, password, host, port, database string) (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
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
