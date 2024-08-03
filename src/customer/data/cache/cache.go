package cache

import (
	"context"
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

func (c Cache) GetTransactionIds(id uint) []uint {
    conn := c.conn
    if conn == nil {
        log.Println(CONNECTION_ERROR)
        return nil
    }
    ctx := context.Background()
    result := conn.Get(ctx, fmt.Sprintf("customers:%d", uint(id)))
    if result.Err() != nil {
    }

    return nil
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
