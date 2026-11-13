package main

type dependencies struct {
	OrderStore
	ItemStore
}

func NewDependencies(ib ResourceBackup[ResItemCreate]) (dependencies, error) {
	itemStore := NewInMemItemStore(ib)
	orderStore := NewInMemOrderStore(itemStore)
	return dependencies{
		OrderStore: orderStore,
		ItemStore:  itemStore,
	}, nil
}
