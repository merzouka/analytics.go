package data

type Retriever interface {
    Close()
    Get()
    GetProductIds(id uint)
    GetTransaction(id uint)
    GetTransactions(ids []uint)
    GetTotal(ids []uint)
}


