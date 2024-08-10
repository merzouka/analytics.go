package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
    // ROWS_DEFAULT = 1_000_000
    ROWS_DEFAULT = 1_0
)

func getRows(name string) uint {
    defs := os.Getenv("ROWS_NUMBER")

    rowMaps := map[string]uint{}
    for _, def := range strings.Split(defs, ",") {
        parts := strings.Split(def, ":")
        if len(parts) > 1 {
            num, err := strconv.ParseUint(parts[1], 10, 64)
            if err != nil {
                log.Fatal(err)
            }

            rowMaps[parts[0]] = uint(num)
            continue
        }

        rowMaps[parts[0]] = ROWS_DEFAULT
    }

    return rowMaps[name]
}

func define(ptr *os.File) {
    def := `CREATE TABLE IF NOT EXISTS transactions (id SERIAL PRIMARY KEY, client_id BIGINT, created_at TIMESTAMP WITH TIME ZONE);
CREATE TABLE IF NOT EXISTS transaction_products (transaction_id BIGINT, product_id BIGINT);`

    _, err := ptr.WriteString(def)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("defined tables successfully")
}

type Transaction struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  uint      `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Products  []TransactionProduct
}

type TransactionProduct struct {
	TransactionID uint `json:"transaction_id"`
	ProductID     uint `json:"product_id"`
}

var id uint = 0
func randId(maxId uint) uint {
    return 1 + uint(rand.Int63n(int64(maxId)))
}
func generateTransaction(maxId uint) Transaction {
    id++
    return Transaction{
        ID: id,
        ClientID: randId(maxId),
        CreatedAt: time.Now().UTC(),
    }
}

func generateTransactionProducts(maxTransactionId, maxProductId uint) TransactionProduct {
    return TransactionProduct{
        TransactionID: randId(maxTransactionId),
        ProductID: randId(maxProductId),
    }
}

func (t Transaction) String() string {
    return fmt.Sprintf("(%d, %d, '%s')", uint(t.ID), uint(t.ClientID), t.CreatedAt.String())
}

func (tp TransactionProduct) String() string {
    return fmt.Sprintf("(%d, %d)", uint(tp.TransactionID), uint(tp.ProductID))
}

func seed(ptr *os.File) {
    query := new(strings.Builder)
    customerIds := getRows("customers")
    productIds := getRows("products")

    maxTransactionId := getRows("transactions")
    for i := uint(0); i < maxTransactionId; i++ {
        transaction := generateTransaction(customerIds).String()
        if i == 0 {
            query.WriteString("INSERT INTO transactions (id, client_id, created_at) VALUES ")
            query.WriteString(transaction)
            continue
        }
        query.WriteString(",\n")
        query.WriteString(transaction)
    }
    query.WriteString(";\n")
    log.Println("generated transactions successfully")

    for i := uint(0); i < getRows("transactionProducts"); i++ {
        transaction := generateTransactionProducts(maxTransactionId, productIds).String()
        if i == 0 {
            query.WriteString("INSERT INTO transaction_products (transaction_id, product_id) VALUES ")
            query.WriteString(transaction)
            continue
        }
        query.WriteString(",\n")
        query.WriteString(transaction)
    } 
    log.Println("associated transactions to products successfully")

    _, err := ptr.WriteString(query.String())
    if err != nil {
        log.Fatal("failed to write to file")
    }
}

func main() {
    dir := os.Getenv("OUTPUT_DIR")
    if dir == "" {
        dir = "./sql/"
    }
    
    def, err := os.Create(fmt.Sprintf("./%s/01-create-tables.sql", dir))
    if err != nil {
        log.Fatal(err)
    }
    ptr, err := os.Create(fmt.Sprintf("./%s/02-populate-tables.sql", dir))
    if err != nil {
        log.Fatal(err)
    }

    go define(def)
    seed(ptr)
}

