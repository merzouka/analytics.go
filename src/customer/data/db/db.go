package db

import (
	"log"

	"github.com/merzouka/analytics.go/customer/data/models"
	"gorm.io/gorm"
)

const (
    DATABASE_CONNECTION_ERROR = "failed to connect to database"
)

type DB struct{
    conn *gorm.DB
}

func (db DB) Close() {
    sqlDB, err := models.GetConn().DB()
    if err != nil {
        log.Println("closing database connection failed")
        return
    }

    sqlDB.Close()
}

func (db DB) GetCustomers(ids []uint, pageSize, page string) []models.Customer {
    return nil
}

var db *DB

func GetInstance() *DB {
    if db != nil {
        return db
    }

    conn := models.GetConn()
    if conn != nil {
        db = &DB{
            conn: conn,
        } 
    }
    return db
}

