package models

import "time"

type Transaction struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  uint      `json:"clientId"`
	CreatedAt time.Time `json:"created_at"`
	Products  []Product
}

type Product struct {
	TransactionID uint `json:"transaction_id"`
	ProductID     uint `json:"product_id"`
}
