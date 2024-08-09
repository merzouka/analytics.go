package data

type Retriever interface {
    Close()
    GetTransaction(id uint)
    GetTransactions(ids []uint)
    GetTotal(ids []uint)
}


