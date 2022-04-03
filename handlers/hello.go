package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	logger *log.Logger
}

func NewHello(logger *log.Logger) *Hello {
	return &Hello{logger}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Hello world", r.Method)

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Opps", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Hello %s", d)
}
