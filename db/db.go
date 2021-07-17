package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open("postgres", fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v ",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("SSLMODE"),
	),
	)
	if err != nil {
		panic(err)
	}
	DB = db
	fmt.Println("DATABASE::CONNECTED")
}

func Migrate(tables ...interface{}) *gorm.DB {
	return DB.AutoMigrate(tables...)
}
