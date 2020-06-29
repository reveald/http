package http

import (
	"net/http"

	"github.com/reveald/reveald"
)

type HTTPAPI struct {
	backend reveald.Backend
	routes  []*Route
	logger  Logger
}

type Option func(*HTTPAPI)

func WithLogger(l Logger) Option {
	return func(h *HTTPAPI) {
		h.logger = l
	}
}

func New(backend reveald.Backend, opts ...Option) *HTTPAPI {
	var routes []*Route
	h := &HTTPAPI{
		backend: backend,
		routes:  routes,
		logger:  &noopLogger{},
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *HTTPAPI) Route(pattern string, opts ...RouteOption) {
	r := NewRoute(pattern, opts...)
	h.routes = append(h.routes, r)
}

func (h *HTTPAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := http.NewServeMux()

	for _, route := range h.routes {
		var handler http.Handler
		handler = queryResultHandler(h.backend, route, h.logger)

		for _, middleware := range route.middleware {
			handler = middleware(handler)
		}
		router.Handle(route.pattern, handler)
	}

	router.ServeHTTP(w, r)
}
