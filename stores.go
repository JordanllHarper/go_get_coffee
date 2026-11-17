package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
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
	is.mu.Lock()
	defer is.mu.Unlock()
	_, found := findFunc(is.items, func(i ResItemCreate) bool {
		return strings.TrimSpace(i.Name) == strings.TrimSpace(item.Name)
	})

	if found {
		return ResItemCreate{}, errAlreadyExists{"Name", item.Name}
	}

	resItem := ResItemCreate{
		ItemId:      uuid.New().String(),
		Name:        item.Name,
		Description: item.Description,
		Cost:        item.Cost,
	}

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
	items, err := os.itemGetter.GetItems(order.ItemIds...)
	if err != nil {
		return ResOrderCreate{}, fmt.Errorf("error fetching items: %w", err)
	}

	os.mu.Lock()
	defer os.mu.Unlock()

	resOrder := ResOrderCreate{
		OrderId: uuid.New().String(),
		ItemIds: order.ItemIds,
		TotalCost: Cost{
			Amount: aggFunc(items, 0, func(acc int, i ResItemCreate) int {
				return acc + i.Cost.Amount
			}),
			Currency: items[0].Cost.Currency,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    StatusInProgress,
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

func (os *inMemOrderStore) CompleteOrder(id string) (ResOrderCreate, error) {
	os.mu.Lock()
	defer os.mu.Unlock()

	idx := findIndexFunc(os.items, func(item ResOrderCreate) bool {
		return item.Id() == id
	})

	if idx == -1 {
		return ResOrderCreate{}, errItemNotFound{id: id}
	}

	os.items[idx].Status = StatusCompleted
	return os.items[idx], nil
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
