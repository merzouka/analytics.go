package models

import (
	"encoding/json"
	"log"
	"time"
)

type Transaction struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  uint      `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Total     uint
	Products  []TransactionProduct
}

type TransactionProduct struct {
	TransactionID uint `json:"transaction_id"`
	ProductID     uint `json:"product_id"`
}

type TransactionProductIDs struct {
	Transaction
	ProductIds []uint
}

func (transaction Transaction) Map() map[string]interface{} {
	str, err := json.Marshal(transaction)
	if err != nil {
		log.Println(err)
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal(str, &result)
	if err != nil {
		log.Println(err)
		return nil
	}

	return result
}
