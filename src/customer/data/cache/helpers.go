package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/merzouka/analytics.go/customer/data/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func queryFailure(query string, err error) string {
    return fmt.Sprintf("query failed: %s, error %s", query, err.Error())
}

func customerCacheKey(id uint) string {
    return fmt.Sprintf("customer:%d", uint(id))
}

func getCustomersFromCache(conn *redis.Client, ids []uint) ([]models.Customer, []uint) {
    ctx := context.Background()
    missed := []uint{}
    var customers []models.Customer

    for _, id := range ids {
        key := customerCacheKey(id)
        result := conn.Get(ctx, key)
        if result.Err() != nil  {
            missed = append(missed, id)
            continue
        }

        var customer models.Customer
        if json.Unmarshal([]byte(result.Val()), &customer) != nil {
            missed = append(missed, id)
            continue
        }

        customers = append(customers, customer)
    }

    return customers, missed
}

func cacheCustomers(client *redis.Client, customers []models.Customer) error {
    failed := []uint{}
    ctx := context.Background()
    for _, customer := range customers {
        jsonCustomer, err := json.Marshal(&customer)
        if err != nil {
            failed = append(failed, customer.ID)
            continue
        }

        result := client.Set(ctx, customerCacheKey(customer.ID), string(jsonCustomer), 0)
        if result.Err() != nil {
            failed = append(failed, customer.ID)
        }
    }

    if len(failed) > 0 {
        return errors.New(fmt.Sprintf("failed to cache ids: %v", failed))
    }

    return nil
}

func getCustomersFromDB(db *gorm.DB, cache *redis.Client, ids []uint) []models.Customer {
    var customers []models.Customer
    result := db.Where("id in (?)", ids).Find(&customers)
    if result.Error != nil {
        log.Println(queryFailure(fmt.Sprintf("retrieval of customers with ids: %v", ids), result.Error))
        return nil
    }

    go cacheCustomers(cache, customers)
    return customers
}
