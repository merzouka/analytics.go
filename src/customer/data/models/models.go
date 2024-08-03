package models

type Customer struct {
    ID uint `gorm:"primaryKey;autoIncrement"`
    Name string
    Age int
    Country *string
    Language *string
    Transactions []Transaction
}

type Transaction struct {
    CustomerID uint
    TransactionID uint
}
