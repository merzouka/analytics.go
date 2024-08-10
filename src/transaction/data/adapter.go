package data

import "github.com/merzouka/analytics.go/transaction/data/models"

type Retriever interface {
    Close()
    GetTransaction(id uint) models.Transaction
    GetTransactions(ids []uint) []models.Transaction
    GetTotal(ids []uint) uint
}


