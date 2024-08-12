package main

import (
	"math/rand"
	"strconv"
	"time"
    "strings"
    "fmt"
    "os"
    "log"
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

        rowMaps[parts[0]] = TRANSACTION_ROWS_DEFAULT
    }

    return rowMaps[name]
}

var id uint = 0
func randId(maxId uint) uint {
    return 1 + uint(rand.Int63n(int64(maxId)))
}

func generateTransaction(customerMax, productMax uint) (string, string) {
    id++
    total, ids := getTotal(productMax)
    rows := []string{}
    for _, productId := range ids {
        rows = append(rows, TransactionProduct{
            TransactionID: id,
            ProductID: productId,
        }.String())
    }

    return Transaction{
        ID: id,
        CustomerID: randId(customerMax),
        Total: total,
        CreatedAt: time.Now().UTC(),
    }.String(), strings.Join(rows, ",\n")
}

func getTotal(maxId uint) (uint, []uint) {
    // generate random product ids to use
    ids := []uint{}
    for i := 0; i < 1 + rand.Intn(15); i++ {
        ids = append(ids, 1+ uint(rand.Int63n(int64(maxId))))
    }

    total := uint(0)
    for _, id := range ids {
        total += products[id - 1].Price
    }

    return total, ids
}

func generateProduct() Product {
    id++
    return Product{
        ID: id,
        Name: fmt.Sprintf("Product %d", uint(id)),
        Price: uint(rand.Intn(4_000_000)),
    }
}

