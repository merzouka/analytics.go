package responses

import (
	"time"

	"github.com/merzouka/analytics.go/transaction/comm"
	"github.com/merzouka/analytics.go/transaction/data/models"
)

type Response struct {
	Data      interface{} `json:"data"`
	Errors    []error     `json:"errors"`
	TimeStamp time.Time   `json:"timestamp"`
}

type Transaction struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ClientID  uint      `json:"clientId"`
	CreatedAt time.Time `json:"createdAt"`
	Total     uint
	Products  []comm.Product
}

func (resp Response) AddError(err error) Response {
	resp.Errors = append(resp.Errors, err)
	return resp
}

func (resp Response) Timestamp() Response {
	resp.TimeStamp = time.Now().UTC()
	return resp
}

func New(data interface{}) Response {
    return Response{
        Data: data,
    }.Timestamp()
}

func (t Transaction) AddProducts(transaction models.Transaction) Transaction {
	ids := []uint{}
	for _, product := range transaction.Products {
		ids = append(ids, product.ProductID)
	}

	products := comm.GetProducts(ids)
	t.Products = append(t.Products, products...)
	return t
}

func (t Transaction) Set(transaction models.Transaction) Transaction {
	t = Transaction{
		ID:        transaction.ID,
		ClientID:  transaction.ClientID,
		Total:     transaction.Total,
		CreatedAt: transaction.CreatedAt,
	}.AddProducts(transaction)
	return t
}
