package db

import (
	"fmt"
	"log"
	"os"

	"github.com/merzouka/analytics.go/api/product/data/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
    conn *gorm.DB
}

var db *DB

func Get() *DB {
    if db != nil {
        return db
    }

    dsn := os.Getenv("DB_URL")
    var err error
    db = &DB{}
    db.conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) 
    if err != nil {
        log.Println(err)
        return nil
    }

    return db
}

func (db *DB) Close() {
    if db == nil {
        return
    }

    sqlDB, err := db.conn.DB()
    if err != nil {
        return
    }

    sqlDB.Close()
}

func (db DB) GetProducts(ids []uint) []models.Product {
    conn := db.conn
    if conn == nil {
        return nil
    }

    var products []models.Product

    if conn.Where("id in (?)", ids).Find(&products).Error != nil {
        log.Println(fmt.Sprintf("product query failed, ids: %v", ids))
        return nil
    }

    return products
}

func (db DB) Conn() *gorm.DB {
    return db.conn
}
