package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/merzouka/analytics.go/api/product/data/db"
	"github.com/merzouka/analytics.go/api/product/data/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Cache struct {
    conn *redis.Client
}

var cache *Cache

func Get() *Cache {
    if cache != nil {
        return cache
    }

    url := os.Getenv("CACHE_URL")
    if url == "" {
        url = "localhost:6379"
    }

    password := os.Getenv("CACHE_PASSWORD")

    conn := redis.NewClient(&redis.Options{
        Addr: url,
        Password: password,
        DB: 0,
    })

    if conn == nil || conn.Ping(context.Background()).Err() != nil {
        log.Println("failed to connect to cache")
        return nil
    }

    cache = &Cache{
        conn: conn,
    }

    return cache
}

func (cache *Cache) Close() {
    if cache == nil {
        return
    }

    cache.conn.Close()
}

func getCacheKey(id uint) string {
    return fmt.Sprintf("product:%d", uint(id))
}

func cacheProducts(cache *redis.Client, products []models.Product) {
    if cache == nil {
        log.Println("connection to cache failed")
        return
    }

    failed := []uint{}
    ctx := context.Background()
    for _, product := range products {
        cachable, err := json.Marshal(product)
        if err != nil {
            failed = append(failed, product.ID)
            continue
        }
        if cache.Set(ctx, getCacheKey(product.ID), []byte(cachable), 0).Err() != nil {
            failed =append(failed, product.ID)
        }
    }
    log.Println(fmt.Sprintf("failed to cache ids: %v", failed))
}

func getProductsFromDB(db *gorm.DB, cache *redis.Client, ids []uint) []models.Product {
    if db == nil {
        log.Println("connection to database failed")
        return nil
    }
    var products []models.Product
    if db.Where("id in (?)", ids).Find(&products).Error != nil {
        log.Println(fmt.Sprintf("query to database failed, ids: %v", ids))
        return nil
    }

    cacheProducts(cache, products)
    return products
}

func getProductsFromCache(cache *redis.Client, ids []uint) ([]models.Product, []uint) {
    if cache == nil {
        log.Println("connection to cache failed")
        return nil, nil
    }

    products := []models.Product{}
    ctx := context.Background()
    misses := []uint{}
    for _, id := range ids {
        result := cache.Get(ctx, getCacheKey(id))
        if result.Err() != nil {
            misses = append(misses, id)
            continue
        }

        var product models.Product
        if json.Unmarshal([]byte(result.Val()), &product) != nil {
            misses = append(misses, id)
            continue
        }
        products = append(products, product)
    }

    return products, misses
}

func (cache Cache) GetProducts(ids []uint) []models.Product {
    conn := cache.conn
    if conn == nil {
        return nil
    }
    db := db.Get().Conn()

    products, misses := getProductsFromCache(conn, ids)
    products = append(products, getProductsFromDB(db, conn, misses)...)
    return products
}
