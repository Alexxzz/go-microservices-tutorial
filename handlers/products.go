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

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle GET Products")

	lp := data.GetProducts()

	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
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
