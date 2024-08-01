package cache

import "github.com/merzouka/analytics.go/api/customer/models"

type Cache struct{
    conn interface{}
}

func (c Cache) GetTransactionIds(id uint) []uint {
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

    return cache
}
