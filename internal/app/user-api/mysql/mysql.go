package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/afernandowex/ps-tag-onboarding-go/internal/app/user-api/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitialiseMySQL() (db *gorm.DB) {
	dsn := "wex_rpc_user:wex_rpc_password@tcp(mysqlDB)/users"
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
	// Delete all recs
	if err = db.Exec("TRUNCATE TABLE users").Error; err != nil {
		panic(fmt.Sprintf("Unable to delete tables %s", err.Error()))
	}
	// Seed some users
	InsertRandomUsers(db)

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

func InsertRandomUsers(db *gorm.DB) {
	user1 := model.User{ID: 1, FirstName: "WexFirst", LastName: "WexLast", Email: "wexfirst.wexlast@wexinc.com", Age: 18}
	user2 := model.User{ID: 2, FirstName: "WexFirst2", LastName: "WexLast2", Email: "wexfirst.wexlast2@wexinc.com", Age: 20}
	user3 := model.User{ID: 3, FirstName: "WexFirst3", LastName: "WexLast3", Email: "wexfirst.wexlast3@wexinc.com", Age: 25}
	users := []model.User{user1, user2, user3}
	tx := db.Create(&users)
	if tx.Error != nil {
		panic(fmt.Sprintf("Unable to delete tables %s", tx.Error))
	}
}
