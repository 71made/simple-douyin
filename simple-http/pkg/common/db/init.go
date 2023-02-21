package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"log"
	"os"
	"simple-main/simple-http/pkg/configs"
	"time"
)

/*
 @Author: 71made
 @Date: 2023/01/24 21:29
 @ProductName: init.go
 @Description: gorm 创建数据库连接
*/

var db *gorm.DB

func Init() {
	// 构建 MySQL 数据库连接
	var err error
	db, err = gorm.Open(
		mysql.Open(configs.MySQLDataBaseDSN),
		&gorm.Config{
			PrepareStmt: true,
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags|log.Lmicroseconds), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
				LogLevel:                  logger.Info,
			}),
		},
	)
	if err != nil {
		panic(err)
	}

	// 使用 tracing
	if err = db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}

func GetInstance() *gorm.DB {
	if db == nil {
		panic("date sources is missing")
	}

	return db
}
