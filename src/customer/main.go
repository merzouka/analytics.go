package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



func main() {
    router := gin.Default()
    defer getSource().Close()

    router.GET("/ping", func(ctx *gin.Context) {
        ctx.String(http.StatusOK, "PONG\n")
    })
    router.GET("/customers/:id/transactions", customerTransactions)
    router.GET("/customers/:id/transactions/total", customerTotal)
    router.GET("/customers/sorted", sortedCustomers)

    router.Run(":8080")
}
