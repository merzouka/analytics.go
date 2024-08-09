package data

import (
	"os"

	"github.com/merzouka/analytics.go/api/product/data/cache"
	"github.com/merzouka/analytics.go/api/product/data/db"
	"github.com/merzouka/analytics.go/api/product/data/models"
)

type Retriever interface {
    Close()
    GetProducts(ids []uint) []models.Product
}

var retriever Retriever

func GetRetriever() Retriever {
    if retriever != nil {
        return nil
    }
    mode := os.Getenv("")
    if mode == "" {
        mode = "CACHE"
    }

    if mode == "CACHE" {
        if db.Get() == nil {
            return nil
        }

        retriever = cache.Get()
        return retriever
    }

    retriever = db.Get()
    return retriever
}
