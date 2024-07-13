package validation

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

const CAT_API_URL = "https://api.thecatapi.com/v1/breeds"

func ValidateBreed(breed string) bool {
	resp, err := http.Get(CAT_API_URL)

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var breeds []struct {
		Name string `json:"name"`
	}
	if err := render.DecodeJSON(resp.Body, &breeds); err != nil {
		return false
	}

	for _, b := range breeds {
		if strings.EqualFold(breed, b.Name) {
			return true
		}
	}

	return false
}
