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
	Timezone  string
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		xff := r.Header.Get("X-Forwarded-For")

		ip := r.RemoteAddr

		if xff != "" {
			ip = strings.TrimSpace(strings.Split(xff, ",")[0])
		}

		timezone := r.URL.Query().Get("tz")

		currentTime := time.Now()

		if timezone != "" {
			location, err := time.LoadLocation(timezone)

			if err == nil {
				currentTime = currentTime.In(location)
			}
		}

		visit := Visit{
			IP:        ip,
			Referrer:  r.Referer(),
			UserAgent: r.UserAgent(),
			Time:      currentTime.Format("2006-01-02 15:04:05 MST"),
			Timezone:  timezone,
		}

		fmt.Fprintln(w, "Ciao da Gio Go!")
		fmt.Fprintln(w)

		fmt.Fprintf(w, "Tracker per %s\n\n", fromSite)
		fmt.Fprintf(w, "IP visitatore: %s\n", visit.IP)
		fmt.Fprintf(w, "Referrer: %s\n", visit.Referrer)
		fmt.Fprintf(w, "Browser: %s\n", visit.UserAgent)
		fmt.Fprintf(w, "Ora locale: %s\n", visit.Time)
		fmt.Fprintf(w, "Timezone: %s\n", visit.Timezone)
	})

	fmt.Printf("Server avviato sulla porta %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
