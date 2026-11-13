package main

import (
	"fmt"
	"net/http"
)

type errItemNotFound struct {
	id string
}

func (err errItemNotFound) Error() string {
	return fmt.Sprintf("item %s not found", err.id)
}

func (err errItemNotFound) Status() int {
	return http.StatusNotFound
}

type errEmptyOrder struct{}

func (err errEmptyOrder) Error() string {
	return "empty order"
}

func (err errEmptyOrder) Status() int {
	return http.StatusBadRequest
}
