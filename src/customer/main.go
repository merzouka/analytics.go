package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
    if err := godotenv.Load(); err != nil {
        log.Println(err)
    }
	defer getSource().Close()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "PONG\n")
	})

	router.GET("/customers/:id/transactions", customerTransactions)
	router.GET("/customers/:id/transactions/total", customerTotal)
	router.GET("/customers/sorted", sortedCustomers)

	log.Println("successfully established connection")
	router.Run(":8080")
}
