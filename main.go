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

	// 1. Port initialisieren
	listenAddr := ":8080"
	customHandlerPort, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if ok {
		listenAddr = ":" + customHandlerPort
	}

	// 2. Quote Handler registrieren
	http.HandleFunc("/api/quote", HandleQuote)

	// 3. Web Server starten
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
