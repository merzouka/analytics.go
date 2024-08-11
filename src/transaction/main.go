package main

import (
	"github.com/gin-gonic/gin"
	"github.com/merzouka/analytics.go/transaction/data"
)

func main() {
    router := gin.Default()
    (*data.GetRetriver()).Close()

    router.GET("/transactions/:id", getTransaction)
    router.GET("/transactions", getTransactions)
    router.GET("/customers/:id/transactions", getCustomerTransactions)
    router.GET("/customers/:id/transactions/total", getTransactionsTotal)

    router.Run(":8080")
}

