package initdatabase

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var defaultDb *gorm.DB
var onceDb sync.Once

func GetDB() *gorm.DB {
	onceDb.Do(func() {
		defaultDb, _ = createDBConnection()
	})
	return defaultDb
}

func createDBConnection() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	)
	dsn := fmt.Sprintf("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}
	// 设置连接池大小
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(100) // 设置最大连接数
	sqlDB.SetConnMaxLifetime(100 * time.Second)
	sqlDB.SetMaxIdleConns(20) // 设置最大空闲连接数
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	return db, nil
}
