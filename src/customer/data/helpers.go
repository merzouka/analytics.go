package data

import "gorm.io/gorm"

func Paginate(pageSize, page int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return nil
    }
}
