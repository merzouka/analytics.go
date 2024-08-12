package responses

import (
	"time"
)

type Response struct {
    Data interface{}
    Errors []error
    TimeStamp time.Time
}

type Product struct {
	ID    uint   `gorm:"primaryKey,autoIncrement" json:"id"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

type Transaction struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  uint      `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Total     uint
	Products  []Product
}

func (resp Response) AddError(err error) Response {
    resp.Errors = append(resp.Errors, err)
    return resp
}

func (resp Response) Timestamp() Response {
    resp.TimeStamp = time.Now().UTC()
    return resp
}
