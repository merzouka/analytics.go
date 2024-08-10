package db

import (
	"log"

	"github.com/merzouka/analytics.go/customer/data/helpers"
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

func (db DB) GetTransactions(id uint) []uint {
    conn := db.conn
    if conn == nil {
        log.Println(DATABASE_CONNECTION_ERROR)
        return nil
    }

    var ids []uint
    result := conn.Model(&models.Transaction{}).Where("customer_id = ?", id).Pluck("transaction_id", &ids)
    if result.Error != nil {
        log.Println("failed to retrieve transactions")
        return nil
    }

    return ids
}

func (db DB) GetSortedCustomers(pageSize, page string) []models.Customer {
    conn := db.conn
    if conn == nil {
        log.Println(DATABASE_CONNECTION_ERROR)
        return nil
    }

    var customers []models.Customer
    ids := conn.
        Table("transactions").
        Select("customer_id").
        Group("customer_id").
        Order("count(customer_id) DESC").
        Scopes(helpers.Paginate(pageSize, page))

    if conn.Where("id in (?)", ids).Preload("Transactions").Find(&customers).Error != nil {
        log.Println("failed to retrieve customers")
        return nil
    }

    return customers
}

func (db DB) GetCustomersByName(name string) []models.Customer {
    conn := db.conn
    if conn == nil {
        log.Println(DATABASE_CONNECTION_ERROR)
        return nil
    }

    var customers []models.Customer
    if conn.Where("name like %?%", name).Find(&customers).Error != nil {
        return nil
    }

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

