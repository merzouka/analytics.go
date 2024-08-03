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
        Scopes(helpers.Paginate(pageSize, page))

    result := conn.Where("id in (?)", ids).Find(&customers)
    if result.Error != nil {
        log.Println("failed to retrieve customers")
        return nil
    }
    return customers
}

func (db DB) GetCustomersForTransactions(ids []uint) []models.Customer {
    conn := db.conn
    if conn == nil {
        log.Println(DATABASE_CONNECTION_ERROR)
        return nil
    }

    customerIds := conn.
        Table("transactions").
        Select("distinct(customer_id)").
        Where("transaction_id in (?)", ids)

    var customers []models.Customer
    result := conn.Where("id in (?)", customerIds).Find(&customers)
    if result.Error != nil {
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

