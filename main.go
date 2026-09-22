package main

import (
	"encoding/json"
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
	City      string
	Region    string
	Country   string
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	type GeoInfo struct {
	City    string `json:"city"`
	Region  string `json:"regionName"`
	Country string `json:"country"`
	}
	func getGeoInfo(ip string) GeoInfo {
		var geo GeoInfo
		url := "http://ip-api.com/json/" + ip
		resp, err := http.Get(url)
		if err != nil {
			return geo
		}
		defer resp.Body.Close()

		json.NewDecoder(resp.Body).Decode(&geo)

		return geo
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		xff := r.Header.Get("X-Forwarded-For")

		ip := r.RemoteAddr

		if xff != "" {
			ip = strings.TrimSpace(strings.Split(xff, ",")[0])
		}

		timezone := r.URL.Query().Get("tz")
		geo := getGeoInfo(ip)

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
			City:      geo.City,
			Region:    geo.Region,
			Country:   geo.Country,
		}

		fmt.Fprintln(w, "Ciao da Gio Go!")
		fmt.Fprintln(w)

		fmt.Fprintf(w, "Tracker per %s\n\n", fromSite)
		fmt.Fprintf(w, "IP visitatore: %s\n", visit.IP)
		fmt.Fprintf(w, "Referrer: %s\n", visit.Referrer)
		fmt.Fprintf(w, "Browser: %s\n", visit.UserAgent)
		fmt.Fprintf(w, "Ora locale: %s\n", visit.Time)
		fmt.Fprintf(w, "Timezone: %s\n", visit.Timezone)
		fmt.Fprintf(w, "Country: %s\n", visit.Country)
		fmt.Fprintf(w, "Region: %s\n", visit.Region)
		fmt.Fprintf(w, "City: %s\n", visit.City)
	})

	fmt.Printf("Server avviato sulla porta %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
