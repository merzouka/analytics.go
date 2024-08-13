package main

import (
	"bytes"
	"encoding/json"
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

func request(endpoint string, v any) error {
    resp, err := http.Get(fmt.Sprintf("http://%s/%s", os.Getenv("TRANSACTION_SERVICE"), endpoint))
    if err != nil || resp.Body == nil {
        return errors.New("failed to connect to transactions service")
    }
    defer resp.Body.Close()

    buffer := new(bytes.Buffer)
    _, err = io.Copy(buffer, resp.Body)
    if err != nil {
        return errors.New("failed to connect to transactions service")
    }

    err = json.Unmarshal(buffer.Bytes(), v)
    if err != nil {
        return err
    }

    return nil
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
    var resp map[string]interface{}
    err = request(endpoint, &resp)
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
    var resp map[string]interface{}
    err = request(endpoint, &resp)
    if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to connect to transactions service",
		})
        return
    }

	ctx.JSON(http.StatusOK, resp)
}

func sortedCustomers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
        "result": data.GetSortedCustomers(ctx.Query("pageSize"), ctx.Query("page")),
    })
}
