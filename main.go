package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// Flight struct
type Flight struct {
	ID           int
	Airline      string
	FlightNumber string
	From         string
	To           string
	Departure    string
	Arrival      string
	Price        float64
}

// Booking struct
type Booking struct {
	Flight     Flight
	Name       string
	Email      string
	Passengers int
	BookingID  string
}

// Sample flights
var flights = []Flight{
	{1, "Air Go", "GO123", "New York", "London", "2026-03-10 08:00", "2026-03-10 20:00", 499.99},
	{2, "Sky Airlines", "SK456", "New York", "London", "2026-03-10 12:00", "2026-03-10 23:50", 529.99},
	{3, "JetFly", "JF789", "New York", "London", "2026-03-10 16:00", "2026-03-11 04:00", 549.99},
}

// Templates with 'mul' function
var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"mul": func(a, b float64) float64 { return a * b },
}).ParseGlob("templates/*.html"))

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/results", resultsHandler)
	http.HandleFunc("/booking", bookingHandler)

	println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Homepage
func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

// Flight search results
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")

	var matched []Flight
	for _, f := range flights {
		if f.From == from && f.To == to {
			matched = append(matched, f)
		}
	}

	templates.ExecuteTemplate(w, "results.html", matched)
}

// Booking page & form submission
func bookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		for _, f := range flights {
			if f.ID == id {
				templates.ExecuteTemplate(w, "booking.html", f)
				return
			}
		}
		http.NotFound(w, r)
	} else if r.Method == http.MethodPost {
		flightID, _ := strconv.Atoi(r.FormValue("flight_id"))
		name := r.FormValue("name")
		email := r.FormValue("email")
		passengers, _ := strconv.Atoi(r.FormValue("passengers"))

		var bookedFlight Flight
		for _, f := range flights {
			if f.ID == flightID {
				bookedFlight = f
				break
			}
		}

		booking := Booking{
			Flight:     bookedFlight,
			Name:       name,
			Email:      email,
			Passengers: passengers,
			BookingID:  "BK" + strconv.Itoa(1000+flightID),
		}

		templates.ExecuteTemplate(w, "confirmation.html", booking)
	}
}
