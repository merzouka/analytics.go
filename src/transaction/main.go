package main

import "github.com/gin-gonic/gin"

func main() {
    router := gin.Default()

    router.GET("/transactions/:id", getTransaction)
    router.GET("/transactions", getTransactions)
    router.GET("/transactions/total", getTransactionsTotal)

    router.Run(":8080")
}

