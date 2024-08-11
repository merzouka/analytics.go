package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/merzouka/analytics.go/customer/data/cache"
	"github.com/merzouka/analytics.go/customer/data/db"
	"github.com/merzouka/analytics.go/customer/data/helpers"
	"github.com/merzouka/analytics.go/customer/data/models"
)

type DataSource interface{
    GetCustomersInOrder(ids []uint) []models.Customer
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

func GetSortedCustomers(pageSize, page string) []models.Customer {
    resp, err := http.Get(fmt.Sprintf("%s/transactions/customers/sorted", os.Getenv("TRANSACTION_SERVICE")))
    if err != nil || resp.Body == nil {
        log.Println("get request failed")
        return nil
    }

    buffer := new(bytes.Buffer)
    if _, err := io.Copy(buffer, resp.Body); err != nil {
        log.Println("failed to retrieve ids")
        return nil
    }

    var result map[string][]uint
    if json.Unmarshal(buffer.Bytes(), &result) != nil {
        log.Println("failed to parse ids")
        return nil
    }

    ids := result["result"]
    ids = helpers.SubsetIds(helpers.CompleteIds(ids), pageSize, page)
    return source.GetCustomersInOrder(ids)
}
