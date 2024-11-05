package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hellofresh/health-go/v5"
	"github.com/kelseyhightower/envconfig"
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
	var config Config
	err := envconfig.Process("dogop", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := http.NewServeMux()

	// Register Quote Handler
	router.HandleFunc("POST /api/quote", HandleQuote)

	// Register Hello Handler
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello DogOp!"))
	})

	// Create Health Check
	h, _ := health.New(
		health.WithComponent(health.Component{
			Name:    "dogop",
			Version: "v0.0.1",
		}),
		health.WithChecks(
			health.Config{
				Name:      "check",
				Timeout:   time.Second * 2,
				SkipOnErr: false,
				Check: func(ctx context.Context) error {
					// check implementation goes here
					return nil
				},
			},
		),
	)

	// Register Handler Function
	router.HandleFunc("GET /health", h.HandlerFunc)

	log.Printf("Listening on port %v", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
}
