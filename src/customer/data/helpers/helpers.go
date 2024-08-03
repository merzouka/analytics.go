package helpers

import "gorm.io/gorm"

func Paginate(pageSize, page int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if page < 1 {
            page = 1
        }
        if pageSize == -1 {
            pageSize = 10
        }

        return db.Offset((page - 1) * pageSize).Limit(page)
    }
}
