package main

import (
	"net/http"

	"github.com/dionisiusrizky/binance-test/handler"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	handler.Init()
	r := chi.NewRouter()

	r.Route("/answer", func(r chi.Router) {
		r.Get("/1", handler.Question1)
		r.Get("/2", handler.Question2)
		r.Get("/3", handler.Question3)
		r.Get("/4", handler.Question4)
		r.HandleFunc("/5", handler.Question5)
	})
	handler.Question6()
	r.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", r)
}
