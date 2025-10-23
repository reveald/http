package main

import (
	"fmt"
	"net/http"

	revealdhttp "github.com/reveald/http/v2"
	"github.com/reveald/reveald/v2"
	"github.com/reveald/reveald/v2/featureset"
)

func main() {
	b, err := reveald.NewElasticBackend([]string{"http://127.0.0.1:9200/"})
	if err != nil {
		panic(err)
	}

	h := revealdhttp.New(b)

	h.Route("/search",
		revealdhttp.WithMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				a := r.Header.Get("authorization")
				if a != "" {
					fmt.Println("Authorized request.")
				}

				next.ServeHTTP(w, r)

				fmt.Println("completed request.")
			})
		}),
		revealdhttp.WithIndex("the-idx"),
		revealdhttp.WithFeatures(
			featureset.NewPaginationFeature(),
			featureset.NewSortingFeature("sort", featureset.WithSortOption("by-range", "range_field", true)),
			featureset.NewStaticFilterFeature(featureset.WithRequiredProperty("maybe_field")),
			featureset.NewStaticFilterFeature(featureset.WithRequiredValue("status.keyword", "Active")),
			featureset.NewDynamicFilterFeature("text_field"),
			featureset.NewHistogramFeature("range_field", featureset.WithInterval(1000))))

	fmt.Print("starting server...\n\n")
	fmt.Println("* navigate to http://localhost:8080/search to see search results")
	fmt.Println("* filter search results by adding a querystring, such as http://localhost:8080/search?text_field=Third")

	err = http.ListenAndServe(":8080", h)
	if err != nil {
		panic(err)
	}
}
