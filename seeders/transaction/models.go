package main

import (
	"fmt"
	"time"
)

type Transaction struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CustomerID  uint      `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Total     uint
	Products  []TransactionProduct
}

type TransactionProduct struct {
	TransactionID uint `json:"transaction_id"`
	ProductID     uint `json:"product_id"`
}

type Product struct {
	ID    uint   `gorm:"primaryKey,autoIncrement" json:"id"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

func (t Transaction) String() string {
    return fmt.Sprintf("(%d, %d, %d, '%s')", uint(t.ID), uint(t.CustomerID), uint(t.Total), t.CreatedAt.Format("2006-01-02 15:04:05-07"))
}

func (tp TransactionProduct) String() string {
    return fmt.Sprintf("(%d, %d)", uint(tp.TransactionID), uint(tp.ProductID))
}

func (p Product) String() string {
    return fmt.Sprintf("(%d, '%s', %d)", uint(p.ID), p.Name, uint(p.Price))
}
