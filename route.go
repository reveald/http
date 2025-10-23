package http

import (
	"net/http"

	"github.com/reveald/reveald/v2"
)

type Route struct {
	pattern    string
	indices    []string
	features   []reveald.Feature
	middleware []func(http.Handler) http.Handler
}

type RouteOption func(*Route)

func WithIndex(idx string) RouteOption {
	return func(r *Route) {
		r.indices = append(r.indices, idx)
	}
}

func WithFeatures(fs ...reveald.Feature) RouteOption {
	return func(r *Route) {
		r.features = append(r.features, fs...)
	}
}

func WithMiddleware(h func(http.Handler) http.Handler) RouteOption {
	return func(r *Route) {
		r.middleware = append(r.middleware, h)
	}
}

func NewRoute(pattern string, opts ...RouteOption) *Route {
	var indices []string
	var features []reveald.Feature
	var middleware []func(http.Handler) http.Handler

	r := &Route{
		pattern,
		indices,
		features,
		middleware,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}
