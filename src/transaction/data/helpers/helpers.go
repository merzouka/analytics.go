package helpers

import (
	"log"
	"strconv"

	"github.com/merzouka/analytics.go/transaction/data/models"
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
