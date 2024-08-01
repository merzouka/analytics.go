package main

import (
	"github.com/gin-gonic/gin"
)



func main() {
    defer setLogger().Close()

    router := gin.Default()

    router.GET("/customers/:id/transactions", customerTransactions)
    router.GET("/customers/:id/transactions/total", customerTotal)
    router.GET("/customers/sorted", sortedCustomers)


    router.Run(":8080")
}
