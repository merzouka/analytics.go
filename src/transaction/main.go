package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/merzouka/analytics.go/transaction/data"
)

func main() {
	router := gin.Default()
	defer data.GetRetriever().Close()
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	router.GET("/transactions/:id", getTransaction)
	router.GET("/transactions", getTransactions)
	router.GET("/transactions/customers/:id", getCustomerTransactions)
	router.GET("/transactions/customers/:id/total", getTransactionsTotal)
	router.GET("/transactions/customers/sorted", getSortedCustomerIds)

	router.Run(":8080")
}
