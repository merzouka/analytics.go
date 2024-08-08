package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func define(dest string) {
    ptr, err := os.Create(fmt.Sprintf("%s/01-create-tables.sql", dest))
    if err != nil {
        log.Fatal(err)
    }

    _, err = ptr.WriteString("CREATE TABLE IF NOT EXISTS products ( id SERIAL PRIMARY KEY, name VARCHAR(255), price BIGINT);\nTRUNCATE TABLE products RESTART IDENTITY CASCADE;\n")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("generated table schema successfully")
}

var id uint = 0
func getProduct() Product {
    id++
    return Product{
        ID: id,
        Name: fmt.Sprintf("Product %d", uint(id)),
        Price: uint(rand.Intn(4_000_000)),
    }
}

const (
    ROWS_DEFAULT = 1_000_000
)

func getRows() uint {
    setting := os.Getenv("ROWS_NUMBER")
    if setting == "" {
        return ROWS_DEFAULT
    }

    parts := strings.Split(setting, ":")
    if len(parts) < 2 {
        return ROWS_DEFAULT
    }

    rows, err := strconv.ParseUint(parts[1], 10, 64)
    if err != nil {
        return ROWS_DEFAULT
    }
    return uint(rows)
}

func (p Product) String() string {
    return fmt.Sprintf("(%d, '%s', %d)", uint(p.ID), p.Name, uint(p.Price))
}

func seed(dest string) {
    ptr, err := os.Create(fmt.Sprintf("%s/02-populate-tables.sql", dest))
    if err != nil {
        log.Fatal(err)
    }

    result := new(strings.Builder)
    for i := uint(0); i < getRows(); i++ {
        if i == 0 {
            result.WriteString("INSERT INTO products (id, name, price) VALUES ")
            result.WriteString(getProduct().String())
            continue
        }
        result.WriteString(",\n")
        result.WriteString(getProduct().String())
    }
    result.WriteString(";")
    _, err = ptr.WriteString(result.String())
    if err != nil {
        log.Fatal("failed to write data")
    }
    log.Println("generated seeder data successfully")
}

func main() {
    outDir := os.Getenv("OUTPUT_DIR")
    if outDir == "" {
        outDir = "./sql/"
    }
    define(outDir)
    seed(outDir)
}
