package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/merzouka/analytics.go/transaction/data/db"
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

type TransactionProducts struct {
    models.Transaction
    ProductIds []uint
}

func (c Cache) Close() {
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
        ids := extractProductIds(transaction.Products)
        transaction.Products = nil
        value, err := json.Marshal(&TransactionProducts{
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
        log.Println("failed to connection to database")
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

func cacheQueryTransactions(cache *redis.Client, ids []uint) ([]TransactionProducts, []uint) {
    if cache == nil {
        return nil, nil
    }
    misses := []uint{}
    transactions := []TransactionProducts{}
    ctx := context.Background()
    for _, id := range ids {
        result := cache.Get(ctx, getKey(id))
        if result.Err() != nil {
            misses = append(misses, id)
        }

        var transaction TransactionProducts
        err := json.Unmarshal([]byte(result.Val()), &transaction)
        if err != nil {
            misses = append(misses, id)
        }

        transactions = append(transactions, transaction)
    }

    return transactions, misses
}

func extractTransactions(transactionProducts []TransactionProducts) []models.Transaction {
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

type Product struct {
        ID    uint   `gorm:"primaryKey,autoIncrement" json:"id"`
        Name  string `json:"name"`
        Price uint   `json:"price"`
}

func getProducts(url string, ids []uint) []Product {
    resp, err := http.Get(url)
    if err != nil || resp.Body == nil{
        log.Println(fmt.Sprintf("request failed for products: %v", ids))
        return nil
    }
    defer resp.Body.Close()

    buffer := new(bytes.Buffer)
    _, err = io.Copy(buffer, resp.Body)
    if err != nil {
        log.Println(fmt.Sprintf("request failed for products: %v", ids))
        return nil
    }

    var products []Product
    if json.Unmarshal(buffer.Bytes(), &products) != nil {
        log.Println(fmt.Sprintf("request failed for products: %v", ids))
        return nil
    }

    return products
}

func extractProductIds(productsArrays ...[]models.Product) []uint {
    m := map[uint]bool{}
    for _, products := range productsArrays {
        for _, product := range products {
            m[product.ProductID] = true
        }
    }
    ids := []uint{}
    for id := range m {
        ids = append(ids, id)
    }

    return ids
}

func (c Cache) GetTransaction(id uint) map[string]interface{} {
    // TODO: move common code out
    cache := c.conn
    if cache == nil {
        return nil
    }

    var transaction TransactionProducts
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

        transaction = TransactionProducts{
            Transaction: result,
            ProductIds: extractProductIds(result.Products),
        }
    }

    str, err := json.Marshal(transaction.Transaction)
    if err != nil {
        log.Println(err)
        return nil
    }
    var resp map[string]interface{}
    err = json.Unmarshal(str, &resp)
    if err != nil {
        log.Println(err)
        return nil
    }

    url := os.Getenv("PRODUCTS_URL")
    resp["products"] = getProducts(url, transaction.ProductIds)
    return resp
}
