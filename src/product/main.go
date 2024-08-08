package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merzouka/analytics.go/api/product/data"
)

var down bool = false

func main() {
    router := gin.Default()

    retriever := data.GetRetriever()
    if retriever == nil {
        router.Use(func(ctx *gin.Context) {
            ctx.JSON(http.StatusInternalServerError, map[string]string{
                "error": "service is down",
            })
        })
    } else {
        defer retriever.Close()
    }

    router.GET("/products", func(ctx *gin.Context) {
        ids := []uint{}
        for _, strId := range strings.Split(ctx.Query("ids"), ",") {
            id, err := strconv.ParseInt(strId, 10, 64)
            if err == nil {
                ctx.JSON(http.StatusInternalServerError, map[string]string{
                    "error": "failed to retrieve products",
                })
                return
            }
            ids = append(ids, uint(id))
        }
        ctx.JSON(http.StatusOK, map[string]interface{}{
            "result": retriever.GetProducts(ids),
        })
    })

    router.Run(":8080")
}
