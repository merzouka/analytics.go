package main

import (
	"bytes"
	"errors"
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

func request(endpoint string) (string, error) {
    resp, err := http.Get(fmt.Sprintf("%s/%s", os.Getenv("TRANSACTION_SERVICE"), endpoint))
    if err != nil || resp.Body == nil {
        return "", errors.New("failed to connect to transactions service")
    }
    defer resp.Body.Close()

    buffer := new(bytes.Buffer)
    _, err = io.Copy(buffer, resp.Body)
    if err != nil {
        return "", errors.New("failed to connect to transactions service")
    }
    return buffer.String(), nil
}

func customerTransactions(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve customer transactions",
		})
		return
	}

    endpoint := fmt.Sprintf("transactions/customers/%d", uint(id))
    resp, err := request(endpoint)
    if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to connect to transactions service",
		})
        return
    }

	ctx.JSON(http.StatusOK, resp)
}

func customerTotal(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "bad id provided",
		})
		return
	}

    endpoint := fmt.Sprintf("transactions/customers/%d/total", uint(id))
    resp, err := request(endpoint)
    if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to connect to transactions service",
		})
        return
    }

	ctx.JSON(http.StatusOK, resp)
}

func sortedCustomers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, data.GetSortedCustomers(ctx.Query("pageSize"), ctx.Query("page")))
}
