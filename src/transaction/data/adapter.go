package data

import (
	"log"
	"os"

	"github.com/merzouka/analytics.go/transaction/data/cache"
	"github.com/merzouka/analytics.go/transaction/data/db"
	"github.com/merzouka/analytics.go/transaction/data/helpers"
	"github.com/merzouka/analytics.go/transaction/data/models"
)

const (
    PRODUCTS_SERVICE_ENV = "PRODUCT_SERVICE"
)

type Retriever interface {
    Close()
    GetTransaction(id uint) *models.TransactionProductIDs
    GetTransactions(ids []uint) []models.Transaction
}

var retriever Retriever

func GetTransaction(retriever Retriever, id uint) map[string]interface{} {
    if retriever == nil {
        log.Println("retriever uninitialized")
        return nil
    }

    transaction := retriever.GetTransaction(id)
    result := transaction.Transaction.Map()
    if result == nil {
        return nil
    }

    url := os.Getenv(PRODUCTS_SERVICE_ENV)
    result["products"] = helpers.GetProducts(url, transaction.ProductIds)
    return result
}

func getCustomerTransactionIds(id uint) []uint {
    var ids []uint
    db := db.Get()
    if db == nil {
        log.Println("failed to retrieve ids")
        return nil
    }

    conn := db.Conn()
    if conn == nil {
        log.Println("failed to retrieve ids")
        return nil
    }

    if conn.Model(&models.Transaction{}).Where("customer_id = ?", id).Pluck("id", &ids).Error != nil {
        log.Println("query for ids failed")
        return nil
    }
    return ids
}

func GetTotal(retriever Retriever, id uint) uint {
    if retriever == nil {
        log.Println("retriever uninitialized")
        return 0
    }

    ids := getCustomerTransactionIds(id)
    transactions := retriever.GetTransactions(ids)
    productArray := [][]models.TransactionProduct{}
    for _, transaction := range transactions {
        productArray = append(productArray, transaction.Products)
    }

    url := os.Getenv(PRODUCTS_SERVICE_ENV)
    products := helpers.GetProducts(url, helpers.ExtractProductIds(productArray...))

    result := uint(0)
    for _, product := range products {
        result += product.Price
    }

    return result
}

func GetCustomerTransactions(retriever Retriever, id uint) []models.Transaction {
    return retriever.GetTransactions(getCustomerTransactionIds(id))
}

func GetRetriver() *Retriever {
    if retriever != nil {
        return &retriever
    }

    mode := os.Getenv("MODE")
    if mode == "" {
        mode = "CACHE"
    }

    if mode == "CACHE" {
        retriever = cache.Get()
    } else {
        retriever = db.Get()
    }

    return &retriever
}

func GetSortedCustomerIds() []uint {
    var ids []uint
    db := db.Get()
    if db == nil {
        log.Println("failed to connect to database")
        return nil
    }
    conn := db.Conn()
    if conn == nil {
        log.Println("failed to connect to database")
        return nil
    }
    conn.Model(&models.Transaction{}).
        Group("customer_id").
        Order("sum(total) DESC").
        Pluck("customer_id", &ids)

    return ids
}
