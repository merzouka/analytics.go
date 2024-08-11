package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merzouka/analytics.go/transaction/data"
)

func getIds(idsStr string) ([]uint, error) {
    if idsStr == "" {
        return nil, errors.New("ids must be specified")
    }
    ids := []uint{}
    for _, idStr := range strings.Split(idsStr, ",") {
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            return nil, errors.New(fmt.Sprintf("failed to parse id: %s", idsStr))
        }

        ids = append(ids, uint(id))
    }

    return ids, nil
}

func getTransactionsTotal(ctx *gin.Context) {
    ids, err := getIds(ctx.Query("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": "bad ids provided",
        })
        return
    }
    clientId := ids[0]

    retriever := data.GetRetriver()
    ctx.JSON(http.StatusOK, map[string]interface{}{
        "result": data.GetTotal(*retriever, clientId),
    })
}

func getTransaction(ctx *gin.Context) {
    ids, err := getIds(ctx.Param("id"))
    if err != nil || len(ids) == 0 {
        ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": "bad id provided",
        })
        return
    }
    transactionId := ids[0]

    retriever := data.GetRetriver()
    ctx.JSON(http.StatusOK, map[string]interface{}{
        "result": data.GetTransaction(*retriever, transactionId),
    })
}

func getTransactions(ctx *gin.Context) {
    ids, err := getIds(ctx.Query("ids"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": "bad ids provided",
        })
        return
    }

    retriever := data.GetRetriver()
    ctx.JSON(http.StatusOK, map[string]interface{}{
        "result": (*retriever).GetTransactions(ids),
    })
}

func getCustomerTransactions(ctx *gin.Context) {
    ids, err := getIds(ctx.Query("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, map[string]string{
            "error": "invalid id provided",
        })
        return
    }

    clientId := ids[0]
    retriever := data.GetRetriver()
    ctx.JSON(http.StatusOK, map[string]interface{}{
        "result": data.GetCustomerTransactions(*retriever, clientId),
    })
}
