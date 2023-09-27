package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/ryansheppard/weather/internal/cache"
)

var (
	getsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "weather_get_calls_total",
		Help: "The total number of processed events",
	}, []string{"caller"})

	cacheHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "weather_cache_hits_total",
		Help: "The total number of cache hits",
	}, []string{"caller"})

	cacheMisses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "weather_cache_misses_total",
		Help: "The total number of cache misses",
	}, []string{"caller"})
)

type HttpRequest struct {
	Endpoint  string
	UserAgent string
	Headers   map[string]string
	Caller    string
	Cache     *cache.Cache
}

type HttpResponseOption func(*HttpRequest)

func WithUserAgent(userAgent string) HttpResponseOption {
	return func(r *HttpRequest) {
		r.UserAgent = userAgent
	}
}

func WithCaller(caller string) HttpResponseOption {
	return func(r *HttpRequest) {
		r.Caller = caller
	}
}

func WithHeaders(headers map[string]string) HttpResponseOption {
	return func(r *HttpRequest) {
		r.Headers = headers
	}
}

func WithCache(cache *cache.Cache) HttpResponseOption {
	return func(r *HttpRequest) {
		r.Cache = cache
	}
}

func NewHttpRequest(endpoint string, opts ...HttpResponseOption) *HttpRequest {
	r := &HttpRequest{
		Endpoint:  endpoint,
		UserAgent: "",
		Headers:   nil,
		Cache:     nil,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *HttpRequest) Get() ([]byte, error) {
	getsProcessed.With(prometheus.Labels{"caller": r.Caller}).Inc()
	if r.Cache != nil {
		rawBody, err := r.Cache.GetKey(r.Endpoint)
		if err != nil {
			slog.Error("Error getting key from cache:", err)
			return []byte{}, err
		}

		if rawBody != nil {
			cacheHits.With(prometheus.Labels{"caller": r.Caller}).Inc()
			return []byte(rawBody.(string)), nil
		}
	}

	cacheMisses.With(prometheus.Labels{"caller": r.Caller}).Inc()

	req, err := http.NewRequest("GET", r.Endpoint, nil)
	if err != nil {
		slog.Error("Error creating request:", err)
		return []byte{}, err
	}

	if r.UserAgent != "" {
		req.Header.Set("User-Agent", r.UserAgent)
	}

	if r.Headers != nil {
		for key, value := range r.Headers {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request:", err)
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	if r.Cache != nil {
		r.Cache.SetKey(r.Endpoint, body, 3600)
	}

	return body, nil
}

func Decode(body []byte, v interface{}) error {
	return json.Unmarshal(body, &v)
}
