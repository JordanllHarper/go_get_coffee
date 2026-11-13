package main

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
)

func jsonEncode(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.MarshalWrite(w, v)
}

func jsonDecode(r io.Reader, v any) error {
	return json.UnmarshalRead(r, v)
}

func handleError(w http.ResponseWriter, err error) {
	var httpErr HttpError
	if ok := errors.As(err, &httpErr); ok {
		w.WriteHeader(httpErr.Status())
		io.WriteString(w, fmt.Sprintf("Error: %v", err))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	log.Println("Error happened internally:", err)
}

func sendStatusCreated(w http.ResponseWriter, location string, v any) error {
	w.Header().Add("Location", location)
	return jsonEncode(w, http.StatusCreated, v)
}

func findFunc[T any](items []T, pred func(T) bool) (T, bool) {
	var t T
	for _, item := range items {
		if pred(item) {
			return item, true
		}
	}
	return t, false
}

func mapFunc[Src any, Dst any](items []Src, pred func(Src) Dst) []Dst {
	dsts := make([]Dst, len(items))
	for i, item := range items {
		dsts[i] = pred(item)
	}
	return dsts
}

func getItemsByIds[T idable](ids []string, items []T) ([]T, error) {
	idItems := make([]T, len(ids))
	for i, id := range ids {
		if idx := slices.IndexFunc(items, func(item T) bool {
			return item.Id() == id
		}); idx != -1 {
			idItems[i] = items[idx]
		} else {
			return nil, errItemNotFound{id}
		}
	}
	return idItems, nil
}
