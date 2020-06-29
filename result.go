package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/reveald/reveald"
)

type Result struct {
	Duration     int64                    `json:"duration_in_ms"`
	NumberOfHits int64                    `json:"total_hit_count"`
	Hits         []map[string]interface{} `json:"hits"`
	Buckets      map[string][]*Bucket     `json:"buckets,omitempty"`
	Pages        *Pagination              `json:"pages,omitempty"`
	Sort         []*SortOption            `json:"sort_options,omitempty"`
}

type Bucket struct {
	Value interface{} `json:"value"`
	Count int64       `json:"count"`
	Query string      `json:"query"`
}

type Pagination struct {
	Count    int     `json:"count"`
	Current  int     `json:"current"`
	Previous *string `json:"previous,omitempty"`
	Next     *string `json:"next,omitempty"`
}

type SortOption struct {
	Name      string `json:"name"`
	Selected  bool   `json:"selected"`
	Ascending bool   `json:"asc"`
	Query     string `json:"query"`
}

func NewResult(r *reveald.Result) *Result {
	request := r.Request()

	buckets := make(map[string][]*Bucket)
	for key, aggs := range r.Aggregations {
		var b []*Bucket
		for _, bs := range aggs {
			request.Set(key, fmt.Sprintf("%v", bs.Value))

			b = append(b, &Bucket{
				Value: bs.Value,
				Count: bs.HitCount,
				Query: makeRequestURL(request),
			})
		}

		if len(b) > 0 {
			buckets[key] = b
		} else {
			buckets[key] = []*Bucket{}
		}
	}

	return &Result{
		Duration:     int64(r.Duration / time.Millisecond),
		NumberOfHits: r.TotalHitCount,
		Hits:         r.Hits,
		Buckets:      buckets,
		Pages:        NewPagination(r),
		Sort:         NewSortOptions(r),
	}
}

func NewPagination(r *reveald.Result) *Pagination {
	if r.Pagination == nil || r.Pagination.PageSize == 0 {
		return nil
	}

	count := int(r.TotalHitCount / int64(r.Pagination.PageSize))
	if count < 1 {
		count = 1
	}

	current := 1
	if r.Pagination.Offset >= r.Pagination.PageSize {
		current = (r.Pagination.Offset / r.Pagination.PageSize) + 1
	}

	p := &Pagination{
		Count:   count,
		Current: current,
	}

	request := r.Request()

	if current > 1 {
		request.Set("offset", fmt.Sprintf("%d", r.Pagination.Offset-r.Pagination.PageSize))
		request.Set("size", fmt.Sprintf("%d", r.Pagination.PageSize))

		url := makeRequestURL(request)
		p.Previous = &url
	}
	if current < count {
		request.Set("offset", fmt.Sprintf("%d", r.Pagination.Offset+r.Pagination.PageSize))
		request.Set("size", fmt.Sprintf("%d", r.Pagination.PageSize))

		url := makeRequestURL(request)
		p.Next = &url
	}

	return p
}

func NewSortOptions(r *reveald.Result) []*SortOption {
	if r.Sorting == nil {
		return nil
	}

	request := r.Request()

	var options []*SortOption
	for _, v := range r.Sorting.Options {
		request.Set(r.Sorting.Param, v.Name)

		options = append(options, &SortOption{
			Name:      v.Name,
			Selected:  v.Selected,
			Ascending: v.Ascending,
			Query:     makeRequestURL(request),
		})
	}

	return options
}

func makeRequestURL(r *reveald.Request) string {
	var parts []string

	for k, p := range r.GetAll() {
		if p.IsRangeValue() {
			min, wmin := p.Min()
			if wmin {
				parts = append(parts, fmt.Sprintf("%s.%s=%.0f", k, reveald.RangeMinParameterName, min))
			}

			max, wmax := p.Max()
			if wmax {
				parts = append(parts, fmt.Sprintf("%s.%s=%.0f", k, reveald.RangeMaxParameterName, max))
			}
		} else {
			for _, v := range p.Values() {
				parts = append(parts, fmt.Sprintf("%s=%v", k, v))
			}
		}
	}

	return strings.Join(parts, "&")
}
