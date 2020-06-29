package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/reveald/reveald"
)

func queryResultHandler(b reveald.Backend, r *Route, l Logger) http.HandlerFunc {
	endpoint := reveald.NewEndpoint(b,
		reveald.WithIndices(r.indices...))

	err := endpoint.Register(r.features...)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request, err := makeRequest(r)
		if err != nil {
			l.Errorf("request failed: %v", err)
			w.WriteHeader(400)
			return
		}

		result, err := endpoint.Execute(context.Background(), request)
		if err != nil {
			l.Errorf("searching failed: %v", err)
			w.WriteHeader(400)
			return
		}

		response := NewResult(result)
		out, err := json.Marshal(response)
		if err != nil {
			l.Errorf("generating response failed: %v", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("content-type", "application/json")
		_, err = w.Write(out)
		if err != nil {
			l.Errorf("writing response failed: %v", err)
			w.WriteHeader(500)
			return
		}
	}
}

func makeRequest(r *http.Request) (*reveald.Request, error) {
	q := reveald.NewRequest()

	for k, v := range r.URL.Query() {
		p := reveald.NewParameter(k, v...)
		q = q.Append(p)
	}

	return q, nil
}
