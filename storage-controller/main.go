package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/hello", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
} 