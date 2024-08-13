package cache

import (
	"context"
	"log"
	"os"

	"github.com/merzouka/analytics.go/customer/data/db"
	"github.com/merzouka/analytics.go/customer/data/models"
	"github.com/redis/go-redis/v9"
)

type Cache struct{
    conn *redis.Client
}

const (
    CACHE_CONNECTION_ERROR = "connection to cache failed"
)

type CustomerTransactions struct {
    models.Customer
    TransactionIds []uint
}

func (c Cache) Close() {
    db.GetInstance().Close()
    if c.IsInvalid() {
        return
    }

    c.conn.Close()
}

func (c Cache) GetCustomersInOrder(ids []uint) []models.Customer {
    if c.IsInvalid() {
        log.Println(CACHE_CONNECTION_ERROR)
        return nil
    }
    cache := c.conn

    db := db.GetInstance()
    if db.IsInvalid() {
        return nil
    }

    customers, misses := getCustomersFromCache(cache, ids)
    dbCustomers := db.GetCustomersInOrder(misses)
    missesIdx := 0
    customersIdx := 0
    result := []models.Customer{}
    for _, id := range ids {
        if len(customers) > 0 && customers[customersIdx].ID == id {
            result = append(result, customers[customersIdx])
            customersIdx++
            continue
        }
        if len(dbCustomers) > 0 {
            result = append(result, dbCustomers[missesIdx])
            missesIdx++
        }
    }

    return result
}

func (c *Cache) IsInvalid() bool {
    return c.conn == nil
}

var cache *Cache = &Cache{}

func GetInstance() *Cache {
    if !cache.IsInvalid() {
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

    cache.conn = client

    if client.Ping(context.Background()).Err() != nil {
        log.Println(CACHE_CONNECTION_ERROR)
        cache.conn = nil
    }


    return cache
}
