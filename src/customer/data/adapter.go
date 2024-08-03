package data

import (
	"github.com/merzouka/analytics.go/customer/data/cache"
	"github.com/merzouka/analytics.go/customer/data/db"
	"github.com/merzouka/analytics.go/customer/data/models"
)

type DataSource interface{
    GetTransactionIds(id uint) []uint;
    GetSortedCustomers(pageSize int, page int) []models.Customer;
    GetCustomersForTransactions(ids []uint) []models.Customer;
    GetCustomersByName(name string) []models.Customer;
}

var retrivers DataSource

func UseCache() *DataSource {
    retrivers = cache.GetInstance()
    return &retrivers
}

func UseDB() *DataSource {
    retrivers = db.GetInstance()
    return &retrivers
}

func GetRetrievers() *DataSource {
    return &retrivers
}
