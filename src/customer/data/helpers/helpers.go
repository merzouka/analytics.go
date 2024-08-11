package helpers

import (
	"log"
	"strconv"

	"github.com/merzouka/analytics.go/customer/data/models"
	"gorm.io/gorm"
)

func Paginate(pageSizeStr, pageStr string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
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

        page, err := strconv.Atoi(pageStr)
        if err != nil {
            log.Println(err)
            page = 1
        }


        return db.Offset((page - 1) * pageSize).Limit(pageSize)
    }
}

// parsePagination parses pageSize and page
func ParsePagination(pageSizeStr, pageStr string) (int, int) {
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

    page, err := strconv.Atoi(pageStr)
    if err != nil {
        log.Println(err)
        page = 1
    }
    return pageSize, page
}

func SubsetIds(ids []uint, pageSizeStr, pageStr string) []uint {
    pageSize, page := ParsePagination(pageSizeStr, pageStr) 
    return ids[(page - 1) * pageSize:page * pageSize]
}

func CompleteIds(original []uint) []uint {
    var ids []uint
    db := models.GetConn()
    if db == nil {
        log.Println("failed to connect to database")
        return original
    }
    err := db.
        Model(&models.Customer{}).
        Where("id not in (?)", original).
        Order("id DESC").
        Pluck("id", &ids).
    Error
    if err != nil {
        log.Println(err)
        return original
    }

    original = append(original, ids...)
    return original
}
