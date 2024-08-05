package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var id uint = 0

func generateCustomer() Customer {
    id++
    return Customer{
        ID: id,
        Name: fmt.Sprintf("Customer %d", uint(id)),
        Age: 3 + rand.Intn(98),
    }
}

func generateTransaction(finalCustomer uint) Transaction {
    id++
    return Transaction{
        TransactionID: id,
        CustomerID: 1 + uint(rand.Int63n(int64(finalCustomer))),
    }
}

const (
    ROWS_DEFAULT = 1_000_000
)

var part int = -1 
func printProgress(current, total uint64, legend string) {
    if int((current * 30) / total) == part {
        return
    }
    part++
    fmt.Printf("[%s%s] %s\n", strings.Repeat("#", part), strings.Repeat(" ", 30 - part), legend)
}

func getRows(key string) uint64 {
    configStr := os.Getenv("ROWS_NUMBER")
    rowsMap := map[string]uint64{}

    for _, tableRows := range strings.Split(configStr, ",") {
        parts := strings.Split(tableRows, ":")
        var rows uint64
        if len(parts) > 1 {
            if parts[1] == "" {
                rows = ROWS_DEFAULT
            } else {
                var err error
                rows, err = strconv.ParseUint(parts[1], 10, 64)
                if err != nil {
                    rows = ROWS_DEFAULT
                }
            }
        } else {
            rows = ROWS_DEFAULT
        }
        rowsMap[parts[0]] = rows
    }

    return rowsMap[key] 
}

func seed(appender Appender) {
    rows := getRows("customers")
    for i := uint64(0); i < rows; i++ {
        model := generateCustomer()
        appender.AddCustomer(&model)
        printProgress(i, rows, "seeding customers...")
    }
    if rows > 0 {
        appender.Finalize()
        log.Println("seeding customers succeeded")
    }

    rows = getRows("transactions")
    finalId := id
    id = 0
    part = -1
    for i := uint64(0); i < rows; i++ {
        model := generateTransaction(finalId)
        appender.AddTransaction(&model)
        printProgress(i, rows, "seeding transactions...")
    }
    if rows > 0 {
        appender.Finalize()
        log.Println("seeding transactions succeeded")
    }
}

func define() {
}

func main() {
    appender := getAppender()
    appender.Define()
    defer appender.Close()
    seed(appender)
}
