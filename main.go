package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const fromSite = "https://glsse.jimdofree.com"

type Visit struct {
	IP        string
	Referrer  string
	UserAgent string
	Time      string
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Recupera l'IP reale del visitatore
		xff := r.Header.Get("X-Forwarded-For")

		ip := r.RemoteAddr

		if xff != "" {
			ip = strings.TrimSpace(strings.Split(xff, ",")[0])
		}

		visit := Visit{
			IP:        ip,
			Referrer:  r.Referer(),
			UserAgent: r.UserAgent(),
			Time:      time.Now().Format("2006-01-02 15:04:05"),
		}

		fmt.Fprintln(w, "Ciao da Gio Go!")
		fmt.Fprintln(w)

		fmt.Fprintf(w, "Tracker per %s\n\n", fromSite)

		fmt.Fprintf(w, "IP visitatore: %s\n", visit.IP)
		fmt.Fprintf(w, "Referrer: %s\n", visit.Referrer)
		fmt.Fprintf(w, "Browser: %s\n", visit.UserAgent)
		fmt.Fprintf(w, "Ora: %s\n", visit.Time)
	})

	fmt.Printf("Server avviato sulla porta %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
