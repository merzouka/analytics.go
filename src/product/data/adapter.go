package data

import (
	"os"

	"github.com/merzouka/analytics.go/api/product/data/cache"
	"github.com/merzouka/analytics.go/api/product/data/db"
	"github.com/merzouka/analytics.go/api/product/data/models"
)

type Retriever interface {
    Close()
    IsInvalid() bool
    GetProducts(ids []uint) []models.Product
}

var retriever Retriever

func GetRetriever() Retriever {
    mode := os.Getenv("MODE")
    if mode == "" {
        mode = "CACHE"
    }

    if mode == "CACHE" {
        retriever = cache.Get()
    } else {
        retriever = db.Get()
    }

    return retriever
}
