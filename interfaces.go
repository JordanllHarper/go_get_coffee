package main

// Http

type HttpError interface {
	error
	Status() int
}

// Idables

type idable interface {
	Id() string
}

// Validation

type Validator interface {
	// Validates the implementor.
	//
	// Implementations should return an empty slice to indicate no errors occurred.
	Validate() []error
}

// Stores

type (
	ItemStore interface {
		ItemAdder
		ItemGetter
	}

	ItemAdder interface {
		AddItem(item ReqItemCreate) (ResItemCreate, error)
	}

	ItemGetter interface {
		GetItems(...string) ([]ResItemCreate, error)
		GetItemById(id string) (ResItemCreate, error)
	}

	OrderStore interface {
		OrderAdder
		OrderGetter
	}

	OrderAdder interface {
		AddOrder(order ReqOrderCreate) (ResOrderCreate, error)
	}

	OrderGetter interface {
		GetOrders(...string) ([]ResOrderCreate, error)
	}
)

type ResourceBackup[T any] interface {
	Backup(items []T)
	Fetch() ([]T, error)
}
