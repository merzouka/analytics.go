package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/merzouka/analytics.go/customer/data/models"
	"github.com/redis/go-redis/v9"
)

type Cache struct{
    conn *redis.Client
}

const (
    CONNECTION_ERROR = "connection to cache failed"
)

type CustomerTransactions struct {
    models.Customer
    TransactionIds []uint
}

func (c Cache) GetTransactionIds(id uint) []uint {
    conn := c.conn
    if conn == nil {
        log.Println(CONNECTION_ERROR)
        return nil
    }
    db := models.GetConn()

    ctx := context.Background()
    result := conn.Get(ctx, fmt.Sprintf("customer:%d", uint(id)))
    if result.Err() != nil {
        var customer models.Customer
        db.Where("id = ?", id).Preload("Transactions").First(&customer)
        ids := []uint{}
        for _, transaction := range customer.Transactions {
            ids = append(ids, transaction.TransactionID)
        }
        customer.Transactions = nil
        str, err := json.Marshal(&CustomerTransactions{
            Customer: customer,
            TransactionIds: ids,
        })

        if err != nil {
            log.Println(fmt.Sprintf("failed to cache response to transaction ids customer:%d", uint(id)))
        } else {
            conn.Set(ctx, fmt.Sprintf("customer:%d", uint(id)), string(str), 0)
        }

        return ids
    }

    var customer CustomerTransactions
    err := json.Unmarshal([]byte(result.Val()), &customer)
    if err != nil {
        log.Println(fmt.Sprintf("request for customer:%d transaction ids failed", uint(id)))
        return nil
    }

    return customer.TransactionIds
}

func (c Cache) GetSortedCustomers(pageSize int, page int) []models.Customer {
    return nil
}

func (c Cache) GetCustomersForTransactions(ids []uint) []models.Customer {
    return nil
}

func (c Cache) GetCustomersByName(name string) []models.Customer {
    return nil
}

var cache *Cache

func GetInstance() *Cache {
    if cache != nil {
        return cache 
    }

    url := os.Getenv("CACHE_URL")
    if url == "" {
        url = "localhost:6379"
    }

    password := os.Getenv("CACHE_PASSWORD")
    client := redis.NewClient(&redis.Options{
        Addr: url,
        Password: password,
        DB: 0,
    })
    cache = &Cache{
        conn: client,
    }

    return cache
}
