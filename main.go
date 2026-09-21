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
		ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {ip = r.RemoteAddr}
		referrer := r.Referer()
		fmt.Fprintf(w, "IP visitatore: %s\n", ip)
		fmt.Fprintf(w, "Referrer: %s\n", referrer)
	})

	http.ListenAndServe(":"+port, nil)
}
