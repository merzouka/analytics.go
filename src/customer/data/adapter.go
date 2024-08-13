package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/merzouka/analytics.go/customer/data/cache"
	"github.com/merzouka/analytics.go/customer/data/db"
	"github.com/merzouka/analytics.go/customer/data/helpers"
	"github.com/merzouka/analytics.go/customer/data/models"
)

type DataSource interface{
    GetCustomersInOrder(ids []uint) []models.Customer
    IsInvalid() bool
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

func getIds(data []interface{}) []uint {
    ids := []uint{}
    for _, d := range data {
        ids = append(ids, uint(d.(float64)))
    }
    return ids
}

func GetSortedCustomers(pageSize, page string) []models.Customer {
    if source.IsInvalid() {
        return nil
    }

    resp, err := http.Get(fmt.Sprintf("http://%s/transactions/customers/sorted", os.Getenv("TRANSACTION_SERVICE")))
    if err != nil || resp.Body == nil {
        log.Println(err)
        return nil
    }

    buffer := new(bytes.Buffer)
    if _, err := io.Copy(buffer, resp.Body); err != nil {
        log.Println(err)
        return nil
    }

    var result map[string]interface{}
    if json.Unmarshal(buffer.Bytes(), &result) != nil {
        log.Println("failed to parse ids")
        return nil
    }

    ids := getIds(result["data"].([]interface{}))
    ids = helpers.SubsetIds(helpers.CompleteIds(ids), pageSize, page)
    return source.GetCustomersInOrder(ids)
}
