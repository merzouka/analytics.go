package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
    if db.IsInvalid() {
        return
    }

    sqlDB, err := db.conn.DB()
    if err != nil {
        log.Println("closing database connection failed")
        return
    }

    sqlDB.Close()
}

func StringifyArray(arr []uint) []string {
    result := []string{}
    for _, elt := range arr {
        result = append(result, strconv.FormatUint(uint64(elt), 10))
    }

    return result
}

func (db DB) GetCustomersInOrder(ids []uint) []models.Customer {
    if db.IsInvalid() {
        log.Println("failed to connect to database")
        return nil
    }
    conn := db.conn

    var customers []models.Customer
    conn.
        Where("id in (?)", ids).
        Order(fmt.Sprintf("array_position(array[%s], id) DESC", strings.Join(StringifyArray(ids), ", "))).
        Find(&customers)

    return customers
}

var db *DB = &DB{}

func (db *DB) IsInvalid() bool {
    return db.conn == nil
}

func GetInstance() *DB {
    if !db.IsInvalid() {
        return db
    }

    db = &DB{
        conn: models.GetConn(),
    } 

    return db
}

