package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/merzouka/analytics.go/transaction/data/db"
	"github.com/merzouka/analytics.go/transaction/data/helpers"
	"github.com/merzouka/analytics.go/transaction/data/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Cache struct {
    conn *redis.Client
}

var cache *Cache

func Get() *Cache  {
    if cache != nil {
        return cache
    }

    url := os.Getenv("CACHE_URL")
    password := os.Getenv("CACHE_PASSWORD")

    client := redis.NewClient(&redis.Options{
        Addr: url,
        Password: password,
        DB: 0,
    })

    if client.Ping(context.Background()).Err() != nil {
        log.Println("failed to connect to cache")
        return nil
    }

    cache = &Cache{
        conn: client,
    }

    return cache
}

func (c *Cache) Close() {
    if c == nil || c.conn == nil {
        return
    }

    c.conn.Close()
}

func getKey(id uint) string {
    return fmt.Sprintf("transaction:%d", uint(id))
}

func cacheTransactions(cache *redis.Client, transactions []models.Transaction) {
    if cache == nil {
        log.Println("failed to connection to cache")
        return
    }

    failed := []uint{}
    ctx := context.Background()
    for _, transaction := range transactions {
        ids := helpers.ExtractProductIds(transaction.Products)
        transaction.Products = nil
        value, err := json.Marshal(&models.TransactionProductIDs{
            ProductIds: ids,
            Transaction: transaction,
        })

        if err != nil {
            failed = append(failed, transaction.ID)
            continue
        }

        if cache.Set(ctx, getKey(transaction.ID), value, 0).Err() != nil {
            failed = append(failed, transaction.ID)
        }
    }

    log.Println(fmt.Sprintf("failed to cache transactions: %v", failed))
}

func dbQueryTransactions(db *gorm.DB, cache *redis.Client, ids []uint) []models.Transaction {
    if db == nil {
        log.Println("connection to database failed")
        return nil
    }

    var transactions []models.Transaction
    if db.Where("id in (?)", ids).Preload("Products").Find(&transactions).Error != nil {
        log.Println("failed to retrieve transactions")
        return nil
    }
    go cacheTransactions(cache, transactions)

    return transactions
}

func cacheQueryTransactions(cache *redis.Client, ids []uint) ([]models.TransactionProductIDs, []uint) {
    if cache == nil {
        return nil, nil
    }
    misses := []uint{}
    transactions := []models.TransactionProductIDs{}
    ctx := context.Background()
    for _, id := range ids {
        result := cache.Get(ctx, getKey(id))
        if result.Err() != nil {
            misses = append(misses, id)
        }

        var transaction models.TransactionProductIDs
        err := json.Unmarshal([]byte(result.Val()), &transaction)
        if err != nil {
            misses = append(misses, id)
        }

        transactions = append(transactions, transaction)
    }

    return transactions, misses
}

func extractTransactions(transactionProducts []models.TransactionProductIDs) []models.Transaction {
    result := []models.Transaction{}
    for _, transaction := range transactionProducts {
        result = append(result, transaction.Transaction)
    }

    return result
}

func (c Cache) GetTransactions(ids []uint) []models.Transaction {
    cache := c.conn
    if cache == nil {
        return nil
    }
    result, misses := cacheQueryTransactions(cache, ids)
    transactions := extractTransactions(result)
    if len(misses) != 0 {
        db := db.Get()
        if db == nil {
            return nil
        }
        transactions = append(transactions, dbQueryTransactions(db.Conn(), cache, misses)...)
    }

    return transactions
}

func (c Cache) GetTransaction(id uint) *models.TransactionProductIDs {
    cache := c.conn
    if cache == nil {
        return nil
    }

    var transaction models.TransactionProductIDs
    result := cache.Get(context.Background(), getKey(id))
    lookup := false
    if result.Err() == nil {
        if json.Unmarshal([]byte(result.Val()), &transaction) != nil {
            lookup = true
        }
    } else {
        lookup = true
    }

    if lookup {
        db := db.Get()
        if db == nil {
            log.Println("failed to retrieve transaction")
            return nil
        }

        var result models.Transaction
        if db.Conn().Preload("Products").First(&result, id).Error != nil {
            log.Println("failed to retrieve transaction")
            return nil
        }

        transaction = models.TransactionProductIDs{
            Transaction: result,
            ProductIds: helpers.ExtractProductIds(result.Products),
        }
    }
    
    return &transaction
}

