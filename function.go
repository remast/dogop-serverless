package function

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
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

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	functions.HTTP("HandleQuote", HandleQuote)
}
