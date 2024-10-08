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

var db *DB = &DB{}

func (db *DB) IsInvalid() bool {
    return db.conn == nil
}

func Get() *DB {
    if !db.IsInvalid() {
        return db
    }

    dsn := os.Getenv("DB_URL")
    var err error
    db.conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) 
    if err != nil {
        log.Println(err)
        db.conn = nil
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
    if db.IsInvalid() {
        return nil
    }
    conn := db.conn

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
