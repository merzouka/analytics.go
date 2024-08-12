package db

import (
	"log"
	"os"
	"time"

	"github.com/merzouka/analytics.go/transaction/data/helpers"
	"github.com/merzouka/analytics.go/transaction/data/models"
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
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
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

func (db *DB) Close() {
	if db == nil {
		return
	}
	sqlDB, err := db.conn.DB()
	if err == nil {
		log.Println("failed to close connection to database")
		return
	}

	sqlDB.Close()
}

func (db DB) GetTransactions(ids []uint) []models.Transaction {
	if db.conn == nil {
		log.Println("failed to retrieve transactions: not connected to database")
		return nil
	}
	conn := db.conn
	var transactions []models.Transaction
	if conn.Preload("Products").Find(&transactions, ids).Error != nil {
		log.Println("failed to retrieve trasactions: query failed")
		return nil
	}

	return transactions
}

func (db DB) GetTransaction(id uint) *models.TransactionProductIDs {
	if db.conn == nil {
		log.Println("failed to retrieve transaction, not connected to database")
		return nil
	}

	var tmp models.Transaction
	if db.conn.Preload("Products").First(&tmp, id).Error != nil {
		log.Println("failed to retrieve transaction")
		return nil
	}

	ids := helpers.ExtractProductIds(tmp.Products)
	tmp.Products = nil

	return &models.TransactionProductIDs{
		Transaction: tmp,
		ProductIds:  ids,
	}
}
