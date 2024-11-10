package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"schneider.vip/problem"
)

type Config struct {
	Port string `default:"8080" envconfig:"PORT"`
}

type Quote struct {
	Age     int      `json:"age" validate:"required"`
	Breed   string   `json:"breed" validate:"required"`
	Tariffs []Tariff `json:"tariffs"`
}

type Tariff struct {
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}

func HandleQuote(w http.ResponseWriter, r *http.Request) {
	var quote Quote
	err := json.NewDecoder(r.Body).Decode(&quote)
	if err != nil {
		problem.New(problem.Wrap(err), problem.Status(http.StatusBadRequest)).WriteTo(w)
		return
	}

	err = validate.Struct(quote)
	if err != nil {
		problem.New(problem.Wrap(err), problem.Status(http.StatusBadRequest)).WriteTo(w)
		return
	}

	log.Println("Calculate quote for", quote)

	tariff := Tariff{Name: "Dog OP _ Basic", Rate: 12.4}
	quote.Tariffs = []Tariff{tariff}

	err = json.NewEncoder(w).Encode(quote)
	if err != nil {
		problem.New(
			problem.Wrap(err),
			problem.Status(http.StatusInternalServerError),
		).WriteTo(w)
	}
}

var validate *validator.Validate

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// 1. Port lesen (Default 8080)
	port := "8080"
	if envPort := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT"); envPort != "" {
		port = envPort
	}

	// 2. Hostname lesen (default 127.0.0.1)
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}

	// 3. Quote Handler registrieren
	http.HandleFunc("/api/quote", HandleQuote)

	// 4. Web Server starten
	log.Fatal(http.ListenAndServe(hostname+":"+port, nil))
}
