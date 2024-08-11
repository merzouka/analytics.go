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
    cache := c.conn
    db := models.GetConn()
    cache.Close()
    sqlDB, err := db.DB()
    if err != nil {
        log.Println("closing database connection failed")
        return
    }

    sqlDB.Close()
}

func (c Cache) GetCustomersInOrder(ids []uint) []models.Customer {
    cache := c.conn
    if cache == nil {
        log.Println(CACHE_CONNECTION_ERROR)
        return nil
    }
    db := db.GetInstance()
    if db == nil {
        return nil
    }

    customers, misses := getCustomersFromCache(cache, ids)
    dbCustomers := db.GetCustomersInOrder(misses)
    missesIdx := 0
    customersIdx := 0
    result := []models.Customer{}
    for _, id := range ids {
        if customers[customersIdx].ID == id {
            result = append(result, customers[customersIdx])
            customersIdx++
            continue
        }
        result = append(result, dbCustomers[missesIdx])
        missesIdx++
    }

    return result
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

    err := client.Ping(context.Background()).Err()
    if err != nil {
        log.Println(CACHE_CONNECTION_ERROR)
        return nil
    }
    log.Println("connected to cache successfully")

    cache = &Cache{
        conn: client,
    }

    return cache
}
