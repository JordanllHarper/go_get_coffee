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

type errAlreadyExists struct {
	property string
	value    string
}

func (err errAlreadyExists) Error() string {
	return fmt.Sprintf("resource with %s property of %s value already exists", err.property, err.value)
}

func (err errAlreadyExists) Status() int {
	return http.StatusConflict
}
