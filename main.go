package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "<h1>CommandName messenger</h1>")
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
