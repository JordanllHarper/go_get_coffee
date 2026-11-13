package main

import (
	"errors"
	"strings"
	"time"
)

// Item creation

type (
	ReqItemCreate struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Cost        Cost   `json:"cost"`
	}

	ResItemCreate struct {
		ItemId      string `json:"itemId"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Cost        Cost   `json:"cost"`
	}
)

// Order creation

type (
	ReqOrderCreate struct {
		ItemIds []string `json:"itemIds"`
	}

	ResOrderCreate struct {
		OrderId   string    `json:"orderId"`
		ItemIds   []string  `json:"itemIds"`
		Cost      int       `json:"cost"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)

// General

type (
	Cost struct {
		// Currency code
		Currency string `json:"currency"`
		// cost in lowest form of currency
		Amount int `json:"amount"`
	}
)

// Validation

func (r ReqOrderCreate) Validate() []error {
	errs := make([]error, 0)
	if len(r.ItemIds) == 0 {
		errs = append(errs, errors.New("no item ids provided"))
	}
	return errs
}

func (r ReqItemCreate) Validate() []error {
	errs := make([]error, 0)
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		errs = append(errs, errors.New("empty name"))
	}
	r.Cost.Currency = strings.TrimSpace(r.Cost.Currency)
	if r.Cost.Currency == "" {
		errs = append(errs, errors.New("empty currency"))
	}
	return errs
}

// Idable

func (item ResItemCreate) Id() string {
	return item.ItemId
}

func (item ResOrderCreate) Id() string {
	return item.OrderId
}
