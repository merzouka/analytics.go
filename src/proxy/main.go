package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
    if err := godotenv.Load(); err != nil {
        log.Println(err)
    }

    router.Use(func(ctx *gin.Context) {
        ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        ctx.Next()
    })

	router.GET("/customers/:id/transactions", customerTransactions)
	router.GET("/customers/:id/transactions/total", customerTotal)
	router.GET("/customers/sorted", sortedCustomers)

	router.GET("bulk", bulk)

	log.Fatal(router.Run(":8080"))
}
