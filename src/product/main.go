package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merzouka/analytics.go/api/product/data"
	"github.com/merzouka/analytics.go/api/product/data/db"
)

var down bool = false

func getIds(query string) []uint {
    if query == "" {
        var ids []uint
        db.Get().Conn().Pluck("id", &ids)
        return ids
    }

    ids := []uint{}
    for _, strId := range strings.Split(query, ",") {
        id, err := strconv.ParseUint(strId, 10, 64)
        if err == nil {
            return nil
        }
        ids = append(ids, uint(id))
    }

    return ids
}

func main() {
    router := gin.Default()

    retriever := data.GetRetriever()
    if retriever == nil {
        router.Use(func(ctx *gin.Context) {
            ctx.JSON(http.StatusInternalServerError, map[string]string{
                "error": "service is down",
            })
        })

        router.Run(":8080")
        return
    }
    defer retriever.Close()

    router.GET("/ping", func(ctx *gin.Context) {
        ctx.String(http.StatusOK, "PONG\n")
    })

    router.GET("/products", func(ctx *gin.Context) {
        ids := getIds(ctx.Query("ids"))
        if ids == nil {
            ctx.JSON(http.StatusInternalServerError, map[string]string{
                "error": "failed to retrieve products",
            })
            return 
        }
        ctx.JSON(http.StatusOK, map[string]interface{}{
            "result": retriever.GetProducts(ids),
        })
    })

    router.Run(":8080")
}
