package main

import (
	"fmt"
	"log"
	"sync"
	"uuid"
)

type inMemItemStore struct {
	mu     sync.RWMutex
	items  []ResItemCreate
	backup ResourceBackup[ResItemCreate]
}

func NewInMemItemStore(b ResourceBackup[ResItemCreate]) *inMemItemStore {
	items, err := b.Fetch()
	if err != nil {
		log.Fatalln("Failed to create in memory item store:", err)
	}
	return &inMemItemStore{
		items:  items,
		backup: b,
	}
}

func (is *inMemItemStore) AddItem(item ReqItemCreate) (ResItemCreate, error) {
	resItem := ResItemCreate{
		ItemId:      uuid.New().String(),
		Name:        item.Name,
		Description: item.Description,
		Cost:        item.Cost,
	}

	is.mu.Lock()
	defer is.mu.Unlock()

	is.items = append(is.items, resItem)

	go is.backup.Backup(is.items)

	return resItem, nil
}

func (is *inMemItemStore) GetItems(ids ...string) ([]ResItemCreate, error) {
	is.mu.RLock()
	defer is.mu.RUnlock()
	if len(ids) == 0 {
		return is.items, nil
	}

	items, err := getItemsByIds(ids, is.items)
	if err != nil {
		return nil, fmt.Errorf("error fetching items by ids: %w", err)
	}
	return items, nil
}

func (is *inMemItemStore) GetItemById(id string) (ResItemCreate, error) {
	is.mu.RLock()
	defer is.mu.RUnlock()

	return getById(is.items, id)
}

type inMemOrderStore struct {
	itemGetter ItemGetter
	mu         sync.RWMutex
	items      []ResOrderCreate
}

func NewInMemOrderStore(itemGetter ItemGetter) *inMemOrderStore {
	return &inMemOrderStore{
		mu:         sync.RWMutex{},
		items:      []ResOrderCreate{},
		itemGetter: itemGetter,
	}
}

func (os *inMemOrderStore) AddOrder(order ReqOrderCreate) (ResOrderCreate, error) {
	_, err := os.itemGetter.GetItems(order.ItemIds...)
	if err != nil {
		return ResOrderCreate{}, fmt.Errorf("error fetching items: %w", err)
	}
	os.mu.Lock()
	defer os.mu.Unlock()
	resOrder := ResOrderCreate{
		OrderId: uuid.New().String(),
		ItemIds: order.ItemIds,
	}
	os.items = append(os.items, resOrder)

	return resOrder, nil
}

func (os *inMemOrderStore) GetOrders(ids ...string) ([]ResOrderCreate, error) {
	os.mu.RLock()
	defer os.mu.RUnlock()
	if len(ids) == 0 {
		return os.items, nil
	}
	items, err := getItemsByIds(ids, os.items)
	if err != nil {
		return nil, fmt.Errorf("error fetching order by ids: %w", err)
	}
	return items, nil
}

func (os *inMemOrderStore) GetOrderById(id string) (ResOrderCreate, error) {
	os.mu.RLock()
	defer os.mu.RUnlock()

	return getById(os.items, id)
}

func getById[T idable](items []T, id string) (T, error) {
	var t T
	v, ok := findFunc(items, func(i T) bool {
		return i.Id() == id
	})

	if !ok {
		return t, errItemNotFound{id: id}
	}

	return v, nil
}
