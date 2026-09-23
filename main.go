package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)
/*struct che rappresenta una singola visita*/
type Visit struct {
	IP        string
	Referrer  string
	UserAgent string
	Time      string
	Timezone  string
	City      string
	Region    string
	Country   string
	Path      string
	DeviceType string
}
/*struct che rappresenta una singola location*/
type GeoInfo struct {
	City    string `json:"city"`
	Region  string `json:"regionName"`
	Country string `json:"country"`
}

var db *sql.DB
/*funzione che a partire dall'IP Prendi l'IP tramite un servizio online restituisce città, regione e Paese*/
func getGeoInfo(ip string) GeoInfo {
	var geo GeoInfo
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return geo
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&geo)
	return geo
}
/*funzione per scrivere messaggi sul browser*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Gio Visitor Tracker")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Tracking endpoint:")
	fmt.Fprintln(w, "/track")
}
/*funzione per leggere la richiesta e recuperare l'IP, geolocalizzare, salvare su DB e restituire un messaggio di conferma*/
func trackHandler(w http.ResponseWriter, r *http.Request) {
	/*recupero Header HTTP*/
	xff := r.Header.Get("X-Forwarded-For")
	/*recupero l'IP reale prendendo il primo della stringa suddivisa da virgola*/
	ip := r.RemoteAddr
	if xff != "" {
		ip = strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	/*recupero la timezone se presente nell'URL della richiesta*/
	timezone := r.URL.Query().Get("tz")
	currentTime := time.Now()
	if timezone != "" {
		location, err := time.LoadLocation(timezone)
		if err == nil {
			currentTime = currentTime.In(location)
		}
	}
	/*recupero la location da IP con la funzione getGeoInfo*/
	geo := getGeoInfo(ip)
	deviceType := getDeviceType(r.UserAgent())
    /*recupero info della singola visita*/
	visit := Visit{
		IP:        ip,
		Referrer:  r.Referer(),
		UserAgent: r.UserAgent(),
		Time:      currentTime.Format("2006-01-02 15:04:05 MST"),
		Timezone:  timezone,
		City:      geo.City,
		Region:    geo.Region,
		Country:   geo.Country,
		Path:      r.URL.Path,
		DeviceType: deviceType,
	}

	// Anti-doppione entro 5 secondi
	var count int
	err := db.QueryRow(
		`SELECT COUNT(*)
		 FROM visits
		 WHERE ip = ?
		 AND visit_time > datetime('now','-5 seconds')`,
		visit.IP,
	).Scan(&count)
	if err == nil && count > 0 {
		fmt.Fprintln(w, "Duplicate visit ignored")
		return
	}
    /*inserisco in DB*/
	result, err := db.Exec(
		`INSERT INTO visits
		(
			visit_time,
			ip,
			country,
			region,
			city,
			timezone,
			browser,
			referrer,
			path,
			device_type
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?)`,
		visit.Time,
		visit.IP,
		visit.Country,
		visit.Region,
		visit.City,
		visit.Timezone,
		visit.UserAgent,
		visit.Referrer,
		visit.Path,
		visit.DeviceType,
	)
	if err != nil {
		fmt.Fprintf(w, "DB ERROR: %v\n", err)
		return
	}
	/*Prendo l'ID della riga appena inserita nel database e la salvo nella variabile id*/
	id, _ := result.LastInsertId()

	fmt.Fprintf(w, "INSERT OK ID=%d\n\n", id)
	fmt.Fprintf(w, "IP: %s\n", visit.IP)
	fmt.Fprintf(w, "Country: %s\n", visit.Country)
	fmt.Fprintf(w, "Region: %s\n", visit.Region)
	fmt.Fprintf(w, "City: %s\n", visit.City)
	fmt.Fprintf(w, "Timezone: %s\n", visit.Timezone)
	fmt.Fprintf(w, "Path: %s\n", visit.Path)
	fmt.Fprintf(w, "Device: %s\n", visit.DeviceType)
}
/*funzione per distinguere il tipo di dispositivo*/
func getDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	// Bot
	if strings.Contains(ua, "bot") ||
		strings.Contains(ua, "crawler") ||
		strings.Contains(ua, "spider") ||
		strings.Contains(ua, "google") {
		return "Bot"
	}
	// Tablet
	if strings.Contains(ua, "ipad") ||
		strings.Contains(ua, "tablet") {
		return "Tablet"
	}
	// Smartphone
	if strings.Contains(ua, "iphone") ||
		strings.Contains(ua, "android") ||
		strings.Contains(ua, "mobile") {
		return "Smartphone"
	}
	// Smart TV
	if strings.Contains(ua, "smart-tv") ||
		strings.Contains(ua, "hbbtv") ||
		strings.Contains(ua, "tv") {
		return "TV"
	}
	// Console
	if strings.Contains(ua, "playstation") ||
		strings.Contains(ua, "xbox") ||
		strings.Contains(ua, "nintendo") {
		return "Console"
	}
	// Desktop
	if strings.Contains(ua, "windows") ||
		strings.Contains(ua, "macintosh") ||
		strings.Contains(ua, "linux") {
		return "Desktop"
	}
	return "Unknown"
}

func main() {
	/*capisco su quale porta deve mettersi in ascolto il server Go*/
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
    /*prendo i due valori dall'Enviroment*/
	dbURL := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	/*Apro il DB e ricevo l'oggetto connessione db*/
	var err error
	db, err = sql.Open(
		"libsql",
		dbURL+"?authToken="+authToken,
	)
	if err != nil {
		panic(err)
	}
	/*eseguo le due funzioni per scrivere sul browser e sul db*/
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/track", trackHandler)
	/*Scrivo nei log del server*/
	fmt.Printf("Server avviato sulla porta %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println(err)
	}
}
