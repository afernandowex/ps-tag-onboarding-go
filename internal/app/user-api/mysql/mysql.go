package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitialiseMySQL() (db *gorm.DB) {
	dsn := os.Getenv("MYSQL_CONNECTION_STRING")
	sqlDB := getSQLConnection(dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect to mysql")
	}

	if err = db.AutoMigrate(&model.User{}); err != nil {
		panic(fmt.Sprintf("failed to instantiate tables %s", err.Error()))
	}
	// // Delete all recs
	// if err = db.Exec("TRUNCATE TABLE users").Error; err != nil {
	// 	panic(fmt.Sprintf("Unable to delete tables %s", err.Error()))
	// }

	return db
}

func getSQLConnection(dsn string) *sql.DB {
	var db *sql.DB
	var err error
	retries := 0
	for retries < 30 {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			retries++
			time.Sleep(1 * time.Second)
			continue
		}
		return db
	}
	return nil
}
