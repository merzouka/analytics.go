package db

import (
	"log"
	"os"

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
    conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Println("failed to connect to database")
        return nil
    }

    db = &DB{
        conn: conn,
    }
    return db
}

func (db *DB) Conn() *gorm.DB {
    if db == nil {
        return nil
    }

    return db.conn
}
