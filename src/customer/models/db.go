package models

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetConn() *gorm.DB {
    if db != nil {
        return db
    }
    var err error

    dsn := os.Getenv("DB_URL")
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Println(err)
        return nil
    }

    return db
}

func CloseConn() {
    if db == nil {
        return
    }
    sqlDB, err := db.DB()
    if err != nil {
        log.Println("failed to close database connection")
        return
    }
    sqlDB.Close()
}
