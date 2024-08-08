package main

import (
	"fmt"
	"log"
	"os"
)

func define(dest string) {
    ptr, err := os.Create(fmt.Sprintf("%s/01-create-tables.sql", dest))
    if err != nil {
        log.Fatal(err)
    }

    _, err = ptr.Write([]byte("CREATE TABLE IF NOT EXISTS products ( id SERIAL PRIMARY KEY, name VARCHAR(255), price BIGINT);\nTRUNCATE TABLE products RESTART IDENTITY CASCADE;\n"))
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
}
