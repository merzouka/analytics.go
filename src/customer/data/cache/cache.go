package cache

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/merzouka/analytics.go/customer/data/helpers"
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

func (c Cache) GetSortedCustomers(pageSize, page string) []models.Customer {
    cache := c.conn
    if cache == nil {
        log.Println(CACHE_CONNECTION_ERROR)
        return nil
    }
    db := models.GetConn()
    if db == nil {
        return nil
    }

    var ids []uint
    result := db.
        Model(&models.Transaction{}).
        Select("customer_id").
        Group("customer_id").
        Order("count(*) DESC").
        Scopes(helpers.Paginate(pageSize, page)).
        Pluck("customer_id", &ids)

    if result.Error != nil {
        log.Println(queryFailure("sorted customers by transactions", result.Error))
        return nil
    }

    customers, misses := getCustomersFromCache(cache, ids)
    customers = append(customers, getCustomersFromDB(db, cache, misses)...)
    return customers
}

func (c Cache) GetCustomersForTransactions(transactionIds []uint) []models.Customer {
    cache := c.conn
    if cache == nil {
        log.Println(CACHE_CONNECTION_ERROR)
        return nil
    }
    db := models.GetConn()
    if db == nil {
        return nil
    }

    var customerIds []uint
    result := db.Model(&models.Transaction{}).Where("transaction_id in (?)").Pluck("distinct(customer_id)", &customerIds)
    if result.Error != nil {
        log.Println(fmt.Sprintf("customer id retrieval failed for transactions: %v", transactionIds))
        return nil
    }

    customers, misses := getCustomersFromCache(cache, customerIds)
    customers = append(customers, getCustomersFromDB(db, cache, misses)...)
    return customers
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
