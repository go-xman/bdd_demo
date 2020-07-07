package application

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func ResolveDB() *gorm.DB {
	if db != nil {
		return db
	}

	db, err := gorm.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return db
}

func CloseDB() {
	if db == nil {
		return
	}

	_ = db.Close()
}
