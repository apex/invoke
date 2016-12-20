package invoke_test

import (
	"log"

	"github.com/apex/invoke"
)

// Input for the function.
type Input struct {
	Host string `json:"host"`
}

// Output for the function.
type Output struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Region      string `json:"region"`
}

// Example invoking a geoip function.
func Example() {
	var out Output
	err := invoke.Sync("geoip", Input{"apex.sh"}, &out)
	if err != nil {
		log.Fatal(err)
	}
}
