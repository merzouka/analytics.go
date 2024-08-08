package models

type Customer struct {
	ID           uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string        `json:"name"`
	Age          int           `json:"age"`
	Country      *string       `json:"country"`
	Language     *string       `json:"language"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	CustomerID    uint `json:"customerId"`
	TransactionID uint `json:"transactionId"`
}
