package comm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	PRODUCTS_SERVICE_ENV = "PRODUCT_SERVICE"
)

func getUrl(endpoint string) string {
	url := os.Getenv(PRODUCTS_SERVICE_ENV)
    return fmt.Sprintf("http://%s%s", url, endpoint)
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
func GetProducts(ids []uint) []Product {
	url := fmt.Sprintf("%s?ids=%s", getUrl("/products"), strings.Join(StringifyArray(ids), ","))
	resp, err := http.Get(url)
	if err != nil || resp.Body == nil {
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

	var products []Product
	if json.Unmarshal(buffer.Bytes(), &products) != nil {
		log.Println(err)
		return nil
	}

	return products
}
