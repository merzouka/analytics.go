package data

import (
	"github.com/merzouka/analytics.go/customer/data/cache"
	"github.com/merzouka/analytics.go/customer/data/db"
	"github.com/merzouka/analytics.go/customer/data/models"
)

type DataSource interface{
    GetSortedCustomers(pageSize string, page string) []models.Customer;
    Close()
}

var source DataSource

func UseCache() *DataSource {
    source = cache.GetInstance()
    return &source
}

func UseDB() *DataSource {
    source = db.GetInstance()
    return &source
}

func GetSource() *DataSource {
    return &source
}
