package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello world", r.Method)

		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Opps", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Hello %s", d)
	})

	http.HandleFunc("/goodbye", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Goodbye world", r.Method)
	})

	http.ListenAndServe(":9090", nil)
}
