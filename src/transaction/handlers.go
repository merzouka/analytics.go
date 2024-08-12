package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merzouka/analytics.go/transaction/data"
	"github.com/merzouka/analytics.go/transaction/responses"
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
	ids, err := getIds(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, responses.New(nil).AddError(err))
		return
	}
	clientId := ids[0]

	retriever := data.GetRetriever()
    if retriever.IsNil() {
		ctx.JSON(http.StatusBadRequest, responses.New(nil).AddError(errors.New("retriever is nil")))
		return
    }
	ctx.JSON(http.StatusOK, responses.New(data.GetTotal(retriever, clientId)))
}

func getTransaction(ctx *gin.Context) {
	ids, err := getIds(ctx.Param("id"))
	if err != nil || len(ids) == 0 {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, responses.New(nil).AddError(err))
		return
	}
	transactionId := ids[0]

	retriever := data.GetRetriever()
    if retriever.IsNil() {
		ctx.JSON(http.StatusBadRequest, responses.New(nil).AddError(errors.New("retriever is nil")))
		return
    }
	ctx.JSON(http.StatusOK, responses.New(data.GetTransaction(retriever, transactionId)))
}

func getTransactions(ctx *gin.Context) {
	ids, err := getIds(ctx.Query("ids"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, responses.New(nil).AddError(err))
		return
	}

	retriever := data.GetRetriever()
    if retriever.IsNil() {
		ctx.JSON(http.StatusBadRequest, responses.New(nil).AddError(errors.New("retriever is nil")))
		return
    }
	ctx.JSON(http.StatusOK, responses.New(retriever.GetTransactions(ids)))
}

func getCustomerTransactions(ctx *gin.Context) {
	ids, err := getIds(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, responses.New(nil).AddError(err))
		return
	}

	clientId := ids[0]
	retriever := data.GetRetriever()
    if retriever.IsNil() {
		ctx.JSON(http.StatusBadRequest, responses.New(nil).AddError(errors.New("retriever is nil")))
		return
    }
	ctx.JSON(http.StatusOK, responses.New(data.GetCustomerTransactions(retriever, clientId)))
}

func getSortedCustomerIds(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, responses.New(data.GetSortedCustomerIds()))
}
