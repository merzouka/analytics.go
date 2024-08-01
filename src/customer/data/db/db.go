package db

import (
	"log"

	"github.com/merzouka/analytics.go/api/customer/data"
	"github.com/merzouka/analytics.go/api/customer/models"
	"gorm.io/gorm"
)

const (
    DATABASE_CONNECTION_ERROR = "failed to connect to database"
)

type DB struct{
    conn *gorm.DB
}

func (db DB) GetTransactionIds(id uint) []uint {
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

func (db DB) GetSortedCustomers(pageSize int, page int) []models.Customer {
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
        Scopes(data.Paginate(pageSize, page))

    result := conn.Where("customer_id in (?)", ids).Find(&customers)
    if result.Error != nil {
        log.Println("failed to retrieve customers")
        return nil
    }
    return customers
}

func (db DB) GetCustomersForTransactions(ids []uint) []models.Customer {
    return nil
}

func (db DB) GetCustomersByName(name string) []models.Customer {
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

