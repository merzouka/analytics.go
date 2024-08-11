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
    sqlDB, err := models.GetConn().DB()
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
    conn := db.conn
    if conn == nil {
        log.Println("failed to connect to database")
        return nil
    }

    var customers []models.Customer
    conn.
        Where("id in (?)", ids).
        Order(fmt.Sprintf("array_position(array[%s], id) DESC", strings.Join(StringifyArray(ids), ", "))).
        Find(&customers)

    return customers
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

