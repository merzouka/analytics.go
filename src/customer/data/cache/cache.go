package cache

import (
	"context"
	"encoding/json"
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

func (c Cache) GetTransactionIds(id uint) []uint {
    conn := c.conn
    if conn == nil {
        log.Println(CACHE_CONNECTION_ERROR)
        return nil
    }
    db := models.GetConn()

    ctx := context.Background()
    result := conn.Get(ctx, customerCacheKey(id))

    if result.Err() != nil {
        var customer models.Customer
        result := db.Where("id = ?", id).Preload("Transactions").First(&customer)
        if result.Error != nil {
            log.Println(queryFailure(fmt.Sprintf("id lookup %d", uint(id)), result.Error))
            return nil
        }

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
            log.Println(fmt.Sprintf("failed to cache response to transaction ids %s", customerCacheKey(id)))
        } else {
            conn.Set(ctx, fmt.Sprintf("customer:%d", uint(id)), string(str), 0)
        }

        return ids
    }

    var customer CustomerTransactions
    err := json.Unmarshal([]byte(result.Val()), &customer)
    if err != nil {
        log.Println(fmt.Sprintf("request for %s transaction ids failed", customerCacheKey(id)))
        return nil
    }

    return customer.TransactionIds
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

func (c Cache) GetCustomersByName(name string) []models.Customer {
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
    result := db.Model(&models.Customer{}).Where("name like %?%", name).Pluck("id", &ids)
    if result.Error != nil {
        log.Println(fmt.Sprintf(""))
    }

    customers, misses := getCustomersFromCache(cache, ids)
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
