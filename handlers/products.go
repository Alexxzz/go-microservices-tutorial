package handlers

import (
	"github.com/Alexxzz/go-microservices-tutorial/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)

	case http.MethodPost:
		p.addProduct(w, r)

	case http.MethodPut:
		{
			rg := regexp.MustCompile(`/(\d+)`)
			g := rg.FindAllStringSubmatch(r.URL.Path, -1)
			if len(g) != 1 || len(g[0]) != 2 {
				p.logger.Println("Invalid URI more than one id or capture group")
				http.Error(w, "Invalid URI", http.StatusBadRequest)
				return
			}

			idString := g[0][1]
			id, err := strconv.Atoi(idString)
			if err != nil {
				p.logger.Println("Invalid URI unable to convert to number")
				http.Error(w, "Invalid URI", http.StatusBadRequest)
				return
			}

			err = p.updateProduct(id, w, r)
			if err == data.ErrProductNotFound {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}

			if err != nil {
				http.Error(w, "Product not found", http.StatusInternalServerError)
				return
			}
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle GET Products")

	lp := data.GetProducts()

	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle POST Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshall JSON", http.StatusBadRequest)
	}

	p.logger.Printf("Prod: %#v", prod)

	data.AddProduct(prod)
}

func (p Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) error {
	p.logger.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall JSON", http.StatusBadRequest)
	}

	return data.UpdateProduct(id, prod)
}
