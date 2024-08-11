package main

import (
	"bytes"
	"fmt"
	"io"
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
		"result": getSource().GetTransaction(uint(id)),
	})
}

func customerTotal(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "bad id provided",
		})
		return
	}

    resp, err := http.Get(fmt.Sprintf("%s/transactions/customers/%d/total", os.Getenv("TRANSACTION_SERVICE"), uint(id)))
    if err != nil || resp.Body == nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to connect to transactions service",
		})
        return 
    }
    defer resp.Body.Close()

    buffer := new(bytes.Buffer)
    _, err = io.Copy(buffer, resp.Body)
    if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to connect to transactions service",
		})
        return 
    }

	ctx.JSON(http.StatusOK, buffer.String())
}

func sortedCustomers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, getSource().GetSortedCustomers(ctx.Query("pageSize"), ctx.Query("page")))
}
