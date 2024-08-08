package helpers

import (
	"strconv"

	"gorm.io/gorm"
)

func Paginate(pageSizeStr, pageStr string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        var page int
        var pageSize int

        if pageSizeStr == "" {
            pageSize = 10
        } else {
            r, err := strconv.Atoi(pageSizeStr)
            if err != nil {
                pageSize = 10
            } else {
                pageSize = r
            }
        }

        if pageStr == "" {
            page = 1 
        } else {
            r, err := strconv.Atoi(pageStr)
            if err != nil {
                page = 10
            } else {
                page = r
            }
        }

        return db.Offset((page - 1) * pageSize).Limit(pageSize)
    }
}
