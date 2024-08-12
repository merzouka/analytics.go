package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/merzouka/analytics.go/transaction/data/models"
	"github.com/merzouka/analytics.go/transaction/responses"
	"gorm.io/gorm"
)

// ExtractProductIds extracts ids from products in productsArrays removing duplicates. Returns the list of ids
func ExtractProductIds(productsArrays ...[]models.TransactionProduct) []uint {
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

func StringifyArray(ids []uint) []string {
    result := []string{}
    for _, id := range ids {
        result = append(result, strconv.FormatUint(uint64(id), 10))
    }

    return result
}

// GetProducts queries the 'products' service and returns products matching ids.
// Returns nil on failure
func GetProducts(target string, ids []uint) []responses.Product {
    url := fmt.Sprintf("%s?ids=%s", target, strings.Join(StringifyArray(ids), ","))
    log.Println(url)
    resp, err := http.Get(url)
    if err != nil || resp.Body == nil{
        log.Println(err)
        return nil
    }
    defer resp.Body.Close()

    buffer := new(bytes.Buffer)
    _, err = io.Copy(buffer, resp.Body)
    if err != nil {
        log.Println(err)
        return nil
    }

    var products []responses.Product
    if json.Unmarshal(buffer.Bytes(), &products) != nil {
        log.Println(err)
        return nil
    }

    return products
}

func Paginate(pageSizeStr, pageStr string) func(*gorm.DB) *gorm.DB {
    return func(d *gorm.DB) *gorm.DB {
        var pageSize int
        var page int
        if pageSizeStr == "" {
            pageSizeStr = "10"
        }
        if pageStr == "" {
            pageStr = "1"
        }

        pageSize, err := strconv.Atoi(pageSizeStr)
        if err != nil {
            log.Println(err)
            pageSize = 10
        }

        page, err = strconv.Atoi(pageStr)
        if err != nil {
            log.Println(err)
            page = 1
        }

        return d.Limit(pageSize).Offset((page - 1) * pageSize)
    }
}
