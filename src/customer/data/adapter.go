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
	"github.com/merzouka/analytics.go/customer/data/models"
)

type DataSource interface{
    GetCustomers(ids []uint, pageSize, page string) []models.Customer
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

func addQuery(original string, values map[string]string) string {
    parts := []string{}
    result := new(strings.Builder)
    result.WriteString(original)
    for key, value := range values {
        if value == "" {
            continue
        }
        parts = append(parts, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
    }
    if len(parts) > 0 {
        result.WriteString("?")
        result.WriteString(strings.Join(parts, "&"))
    }
    return result.String()
}

func GetSortedCustomers(pageSize, page string) []models.Customer {
    resp, err := http.Get(addQuery(fmt.Sprintf("%s/transactions/customers/sorted", os.Getenv("TRANSACTION_SERVICE")), map[string]string{
        "pageSize": pageSize,
        "page": page,
    }))
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
    return source.GetCustomers(ids, pageSize, page)
}
