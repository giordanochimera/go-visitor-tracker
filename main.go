package main

import (
	"fmt"
	"net/http"
	"os"
)

const fromSite = "https://glsse.jimdofree.com"

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Ciao da Gio Go!")
		fmt.Fprintf(w, "Tracker per %s", fromSite)
	})

	http.ListenAndServe(":"+port, nil)
}
