package http

import (
	"encoding/json"
	"net/http"

	"github.com/reveald/reveald/v2"
)

func QueryReader(r *http.Request) ([]reveald.Parameter, error) {
	params := []reveald.Parameter{}
	for k, v := range r.URL.Query() {
		params = append(params, reveald.NewParameter(k, v...))
	}

	return params, nil
}

func JSONBodyReader(r *http.Request) ([]reveald.Parameter, error) {
	var body map[string][]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	params := []reveald.Parameter{}

	for k, v := range body {
		params = append(params, reveald.NewParameter(k, v...))
	}

	return params, nil
}
