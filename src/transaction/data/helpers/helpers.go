package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/merzouka/analytics.go/transaction/data/models"
)

type Product struct {
        ID    uint   `gorm:"primaryKey,autoIncrement" json:"id"`
        Name  string `json:"name"`
        Price uint   `json:"price"`
}

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


// GetProducts queries the 'products' service and returns products matching ids.
// Returns nil on failure
func GetProducts(url string, ids []uint) []Product {
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
