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
    router.GET("/transactions/customers/:id", getCustomerTransactions)
    router.GET("/transactions/customers/:id/total", getTransactionsTotal)
    router.GET("/transactions/customers/sorted", getSortedCustomerIds)

    router.Run(":8080")
}
