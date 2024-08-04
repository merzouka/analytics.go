package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/merzouka/analytics.go/customer/data"
)

var source data.DataSource

func getSource() data.DataSource {
    if source != nil {
        return source
    }
    mode := os.Getenv("MODE")
    if mode == "" {
        mode = "CACHE"
    }

    switch mode {
    case "CACHE":
        source = *data.UseCache()
        break
    case "DB":
        source = *data.UseDB()
        break
    }

    return source
}

func customerTransactions(ctx *gin.Context) {
    id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, map[string]string{
            "error": "failed to retrieve customer transactions",
        })
        return
    }

    ctx.JSON(http.StatusOK, map[string]interface{}{
        "result": getSource().GetTransactionIds(uint(id)),
    })
}

func customerTotal(ctx *gin.Context) {
    ctx.String(http.StatusOK, "unimplemented\n")
}

func sortedCustomers(ctx *gin.Context) {
    ctx.JSON(http.StatusOK, getSource().GetSortedCustomers(ctx.Query("pageSize"), ctx.Query("page")))
}
