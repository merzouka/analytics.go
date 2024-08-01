package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func configLogger() *os.File {
    path := os.Getenv("LOGS_PATH")
    if path == "" {
        path = "./logs"
    }
    prefix := os.Getenv("SERVICE_NAME")
    if prefix != "" {
        prefix = fmt.Sprintf("[%s] ", prefix)
    }

    f, err := os.OpenFile(path, os.O_APPEND | os.O_RDWR | os.O_CREATE, 0644)
    if err != nil {
        log.Fatal("failed to set logger")
    }
    log.SetOutput(f)
    log.SetPrefix(prefix)
    return f
}

func main() {
    defer configLogger().Close()

    router := gin.Default()

    router.GET("/customers/:id/transactions", customerTransactions)
    router.GET("/customers/:id/transactions/total", customerTotal)
    router.GET("/customers/sorted", sortedCustomers)

    router.GET("/transactions/:range", transactionsInRange)

    router.GET("/products/sorted", sortedProducts)
    router.GET("/products/:id/transactions", productTransactions)
    router.GET("/products/:id/transactions/total", productTotal)

    log.Fatal(router.Run(":8080"))
}
