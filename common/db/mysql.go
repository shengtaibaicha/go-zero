package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func Init(datasource string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(datasource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 日志模式
	})
	if err != nil {
		log.Printf("mysql init err:%v", err)
	}
	return db
}
