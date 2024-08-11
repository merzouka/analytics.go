package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
    // TRANSACTION_ROWS_DEFAULT = 1_000_000
    TRANSACTION_ROWS_DEFAULT = 1_0
    PRODUCT_ROWS_DEFAULT = 1_000
)

func define(ptr *os.File) {
    def := `CREATE TABLE IF NOT EXISTS transactions (id SERIAL PRIMARY KEY, client_id BIGINT, total BIGINT, created_at TIMESTAMP WITH TIME ZONE);
CREATE TABLE IF NOT EXISTS transaction_products (transaction_id BIGINT, product_id BIGINT);`

    _, err := ptr.WriteString(def)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("defined tables successfully")
}

func seed(ptr *os.File, pivotPtr *os.File) {
    query := new(strings.Builder)
    pivotQuery := new(strings.Builder)
    customerMaxId := getRows("customers")
    productMaxId := getRows("products")

    maxTransactionId := getRows("transactions")
    for i := uint(0); i < maxTransactionId; i++ {
        transaction, pivot := generateTransaction(customerMaxId, productMaxId)
        if i == 0 {
            query.WriteString(`TRUNCATE TABLE transactions RESTART IDENTITY CASCADE;
INSERT INTO transactions (id, client_id, total, created_at) VALUES `)
            query.WriteString(transaction)
            pivotQuery.WriteString(`TRUNCATE TABLE transaction_products RESTART IDENTITY CASCADE;
INSERT INTO transaction_products (transaction_id, product_id) VALUES `)
            pivotQuery.WriteString(pivot)
            continue
        }
        query.WriteString(",\n")
        query.WriteString(transaction)
        pivotQuery.WriteString(",\n")
        pivotQuery.WriteString(pivot)
    }
    query.WriteString(";\n")
    pivotQuery.WriteString(";\n")

    _, err := ptr.WriteString(query.String())
    if err != nil {
        log.Fatal(err)
    }
    _, err = pivotPtr.WriteString(pivotQuery.String())
    if err != nil {
        log.Fatal(err)
    }

    log.Println("generated transactions and transaction_products successfully")
}

func prodDefine(ptr *os.File) {
    def := `CREATE TABLE IF NOT EXISTS products ( id SERIAL PRIMARY KEY, name VARCHAR(255), price BIGINT);
TRUNCATE TABLE products RESTART IDENTITY CASCADE;\n`
    _, err := ptr.WriteString(def)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("defined products table successfully")
}

var products []Product
func prodSeed(ptr *os.File) {
    result := new(strings.Builder)
    for i := uint(0); i < getRows("products"); i++ {
        product := generateProduct()
        products = append(products, product)
        if i == 0 {
            result.WriteString("INSERT INTO products (id, name, price) VALUES ")
            result.WriteString(product.String())
            continue
        }
        result.WriteString(",\n")
        result.WriteString(product.String())
    }
    result.WriteString(";")

    _, err := ptr.WriteString(result.String())
    if err != nil {
        log.Fatal("failed to write data")
    }
    log.Println("generated products successfully")
}

func main() {
    prodDir := os.Getenv("PRODUCT_OUTPUT_DIR")
    if prodDir == "" {
        prodDir = "./sql/product"
    }
    def, err := os.Create(fmt.Sprintf("%s/01-create-table.sql", prodDir))
    if err != nil {
        log.Fatal(err)
    }
    ptr, err := os.Create(fmt.Sprintf("%s/02-populate-tables.sql", prodDir))
    if err != nil {
        log.Fatal(err)
    }

    prodDefine(def)
    prodSeed(ptr)
    id = 0

    dir := os.Getenv("OUTPUT_DIR")
    if dir == "" {
        dir = "./sql/"
    }
    
    def, err = os.Create(fmt.Sprintf("%s/01-create-tables.sql", dir))
    if err != nil {
        log.Fatal(err)
    }
    ptr, err = os.Create(fmt.Sprintf("%s/02-populate-main.sql", dir))
    if err != nil {
        log.Fatal(err)
    }
    pivot, err := os.Create(fmt.Sprintf("%s/03-populate-pivot.sql", dir))
    if err != nil {
        log.Fatal(err)
    }

    define(def)
    seed(ptr, pivot)
}

