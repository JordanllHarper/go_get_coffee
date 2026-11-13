package main

import (
	"fmt"
	"log"
	"net/http"
)

func setupV1Endpoints(mux *http.ServeMux, d dependencies, srv string) {
	mux.HandleFunc("POST /v1/items", handlePostItem(d, srv))
	mux.HandleFunc("GET /v1/items", handleGetItems(d))
	mux.HandleFunc("GET /v1/items/{id}", handleGetItemById(d))

	mux.HandleFunc("POST /v1/orders", handlePostOrders(d, srv))
	mux.HandleFunc("GET /v1/orders", handleGetOrders(d))
	mux.HandleFunc("GET /v1/orders/{id}", handleGetOrderById(d))
}

func handleGetItems(d ItemGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		ids := query["item"]

		items, err := d.GetItems(ids...)
		if err != nil {
			handleError(w, err)
			return
		}
		if err := jsonEncode(w, http.StatusOK, items); err != nil {
			log.Println("Error sending items", err)
		}
	}
}

func handleGetItemById(d ItemGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getResourceById(w, r, d.GetItemById)
	}
}

func handlePostItem(ia ItemAdder, srv string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReqItemCreate
		if err := jsonDecode(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if ok := validate(w, req); !ok {
			return
		}
		res, err := ia.AddItem(req)
		if err != nil {
			handleError(w, err)
			return
		}
		if err := sendStatusCreated(w, fmt.Sprintf("%s/v1/items/%s", srv, res.ItemId), res); err != nil {
			log.Println("Error sending created item:", err)
		}
	}
}

func handlePostOrders(oa OrderAdder, srv string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReqOrderCreate
		if err := jsonDecode(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if ok := validate(w, req); !ok {
			return
		}
		order, err := oa.AddOrder(req)
		if err != nil {
			handleError(w, err)
			return
		}

		if err := sendStatusCreated(w, fmt.Sprintf("%s/v1/items/%s", srv, order.OrderId), order); err != nil {
			log.Println("Error sending created order:", err)
		}
	}
}

func handleGetOrders(d OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := d.GetOrders()
		if err != nil {
			handleError(w, err)
			return
		}
		if err := jsonEncode(w, http.StatusOK, items); err != nil {
			log.Println("Error sending items", err)
		}
	}
}

func handleGetOrderById(d OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getResourceById(w, r, d.GetOrderById)
	}
}

func validate(w http.ResponseWriter, v Validator) (ok bool) {
	errs := v.Validate()
	if len(errs) == 0 {
		return true
	}

	errStrs := mapFunc(errs, func(err error) string {
		return err.Error()
	})

	if err := jsonEncode(w, http.StatusBadRequest, errStrs); err != nil {
		log.Println("Error encoding errors:", err)
	}
	return false
}

func getResourceById[T any](w http.ResponseWriter, r *http.Request, getResource func(id string) (T, error)) {
	id := r.PathValue("id")
	resource, err := getResource(id)
	if err != nil {
		handleError(w, err)
		return
	}
	if err := jsonEncode(w, http.StatusOK, resource); err != nil {
		log.Println("Error sending items", err)
	}
}
