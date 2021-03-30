package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/dionisiusrizky/binance-test/entity"
	"github.com/dionisiusrizky/binance-test/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	absDelta = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "binance",
		Name:      "binance_absolute_delta",
		Help:      "Absolute delta of the price spread",
	}, []string{"symbol"})
	deltaMap6 map[string]float64
	deltaMap5 map[string]entity.SpreadDelta
)

func Question1(w http.ResponseWriter, r *http.Request) {
	topVolumes, err := usecase.HighestPrices("BTC", "volume", 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
	}

	result, err := json.Marshal(topVolumes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
	}
	w.Write(result)
}

func Question2(w http.ResponseWriter, r *http.Request) {
	topTrades, err := usecase.HighestPrices("USDT", "trade", 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
	}

	result, err := json.Marshal(topTrades)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
	}
	w.Write(result)
}

func Question3(w http.ResponseWriter, r *http.Request) {
	var totalNotionals []entity.Notional
	topVolumes, err := usecase.HighestPrices("BTC", "volume", 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
	}

	for _, price := range topVolumes {
		totalNotional, err := usecase.TotalNotional(price.Symbol, 200)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
		}
		totalNotionals = append(totalNotionals, totalNotional)
	}

	result, err := json.Marshal(totalNotionals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
	}
	w.Write(result)
}

func Question4(w http.ResponseWriter, r *http.Request) {
	var priceSpreads []entity.PriceSpread
	var spread float64

	topTrades, err := usecase.HighestPrices("USDT", "trade", 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
	}

	for _, price := range topTrades {
		spread = price.AskPrice - price.BidPrice
		priceSpreads = append(priceSpreads, entity.PriceSpread{
			Symbol: price.Symbol,
			Spread: fmt.Sprintf("%.8f", spread),
		})
	}

	result, err := json.Marshal(priceSpreads)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
	}

	w.Write(result)
}

func Question5(w http.ResponseWriter, r *http.Request) {
	var spread, delta, pastDelta float64
	var deltastr string

	topTrades, err := usecase.HighestPrices("USDT", "trade", 5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for _, price := range topTrades {
		spread = price.AskPrice - price.BidPrice
		pastDelta, _ = strconv.ParseFloat(deltaMap5[price.Symbol].Delta, 64)
		delta = math.Abs(pastDelta - spread)
		deltastr = fmt.Sprintf("%.8f", delta)
		deltaMap5[price.Symbol] = entity.SpreadDelta{
			Spread: fmt.Sprintf("%.8f", spread),
			Delta:  deltastr,
		}
	}

	filePath := path.Join("5.html")
	tmpl := template.Must(template.ParseFiles(filePath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, deltaMap5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Question6() {
	go func() {
		for {
			var spread float64

			topTrades, err := usecase.HighestPrices("USDT", "trade", 5)
			if err != nil {
			}

			for _, price := range topTrades {
				spread = price.AskPrice - price.BidPrice
				deltaMap6[price.Symbol] = math.Abs(deltaMap6[price.Symbol] - spread)
			}

			for symbol, delta := range deltaMap6 {
				absDelta.WithLabelValues(symbol).Set(delta)
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

func Init() {
	deltaMap5 = make(map[string]entity.SpreadDelta)
	deltaMap6 = make(map[string]float64)
}
