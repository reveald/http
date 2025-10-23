package http

import (
	"net/http"

	"github.com/reveald/reveald/v2"
)

type HTTPAPI struct {
	backend reveald.Backend
	routes  []*Route
	logger  Logger
	reader  ParamReader
}

type Option func(*HTTPAPI)

func WithLogger(l Logger) Option {
	return func(h *HTTPAPI) {
		h.logger = l
	}
}

type ParamReader = func(*http.Request) ([]reveald.Parameter, error)

func WithParamReader(readers ...ParamReader) Option {
	return func(h *HTTPAPI) {
		h.reader = func(r *http.Request) ([]reveald.Parameter, error) {
			params := []reveald.Parameter{}
			for _, reader := range readers {
				ps, err := reader(r)
				if err != nil {
					return nil, err
				}

				params = append(params, ps...)
			}
			return params, nil
		}
	}
}

func New(backend reveald.Backend, opts ...Option) *HTTPAPI {
	var routes []*Route
	h := &HTTPAPI{
		backend: backend,
		routes:  routes,
		logger:  &noopLogger{},
		reader:  QueryReader,
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
		handler = queryResultHandler(h.backend, route, h.logger, h.reader)

		for _, middleware := range route.middleware {
			handler = middleware(handler)
		}
		router.Handle(route.pattern, handler)
	}

	router.ServeHTTP(w, r)
}
