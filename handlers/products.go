// Package handlers
// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"github.com/Alexxzz/go-microservices-tutorial/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// A list of products in the response
// swagger:response productsResponse
type productsResponse struct {
	// All products in the system
	// in: body
	Body []data.Product
}

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle POST Products")

	prod := r.Context().Value(productContextKey).(*data.Product)

	p.logger.Printf("Prod: %#v", prod)

	data.AddProduct(prod)
}

func (p Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle PUT Product")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(productContextKey).(*data.Product)
	p.logger.Println("\tid:", id)

	data.UpdateProduct(id, prod)
}

type contextKey string

const productContextKey contextKey = "__ProductContextKey"

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.logger.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Unable to unmarshall JSON", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.logger.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Unable to validate product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), productContextKey, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
