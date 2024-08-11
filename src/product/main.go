package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merzouka/analytics.go/api/product/data"
	"github.com/merzouka/analytics.go/api/product/data/db"
	"github.com/merzouka/analytics.go/api/product/data/models"
)

var down bool = false

func getIds(query string) []uint {
    if query == "" {
        var ids []uint
        db.Get().Conn().Model(&models.Product{}).Pluck("id", &ids)
        return ids
    }

    ids := []uint{}
    for _, strId := range strings.Split(query, ",") {
        id, err := strconv.ParseUint(strId, 10, 64)
        if err != nil {
            log.Println(err)
            return nil
        }
        ids = append(ids, uint(id))
    }

    return ids
}

func main() {
    router := gin.Default()

    var retriever data.Retriever
    router.GET("/ping", func(ctx *gin.Context) {
        ctx.String(http.StatusOK, "PONG\n")
    })

    router.GET("/products", func(ctx *gin.Context) {
        r := data.GetRetriever()
        retriever = r
        ids := getIds(ctx.Query("ids"))
        if ids == nil || r == nil {
            ctx.JSON(http.StatusInternalServerError, map[string]string{
                "error": "failed to retrieve products",
            })
            return 
        }
        ctx.JSON(http.StatusOK, map[string]interface{}{
            "result": r.GetProducts(ids),
        })
    })

    router.GET("/product/:id", func(ctx *gin.Context) {
        r := data.GetRetriever()
        retriever = r
        ids := getIds(ctx.Param("id"))
        if ids == nil || r == nil {
            ctx.JSON(http.StatusInternalServerError, map[string]string{
                "error": "failed to retrieve products",
            })
            return 
        }

        products := r.GetProducts(ids)
        if len(products) == 0 {
            ctx.JSON(http.StatusOK, nil)
            return
        }
        ctx.JSON(http.StatusOK, products[0])
    })

    router.Run(":8080")
    if retriever != nil {
        retriever.Close()
    }
}
