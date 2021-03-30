package usecase

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/dionisiusrizky/binance-test/adapter"
	"github.com/dionisiusrizky/binance-test/entity"
)

//Factor means by which factor the data is sorted (Volume, Trade Count)
//Length is slice length to be returned (-1 means all data is returned)
func HighestPrices(symbol string, factor string, length int) ([]entity.Price, error) {
	//Fetching symbol with quoteAsset BTC
	symbols, err := adapter.FetchSymbols(symbol)
	if err != nil {
		return nil, err
	}

	//Fetching volumes map and filtering them by fetched symbols above
	//Put them into struct for sorting
	var allPrices []entity.Price
	var filteredPrices []entity.Price

	allPrices, err = adapter.Fetch24HRPrice()
	if err != nil {
		return nil, err
	}

	for _, price := range allPrices {
		if contains(symbols, price.Symbol) {
			filteredPrices = append(filteredPrices, price)
		}
	}

	//Sort the symbol by volumes or trade count
	if factor == "volume" {
		sort.Slice(filteredPrices[:], func(i, j int) bool {
			return filteredPrices[i].Volume > filteredPrices[j].Volume
		})
	} else if factor == "trade" {
		sort.Slice(filteredPrices[:], func(i, j int) bool {
			return filteredPrices[i].Count > filteredPrices[j].Count
		})
	}
	if length >= 0 {
		return filteredPrices[:length], nil
	} else {
		return filteredPrices, nil
	}
}

func TotalNotional(symbol string, limit int) (entity.Notional, error) {
	var bids, asks, price, quantity float64
	order, err := adapter.FetchOrderBook(symbol, limit)
	if err != nil {
		return entity.Notional{}, err
	}
	//Summing bids and asks
	for _, bid := range order.Bids {
		price, _ = strconv.ParseFloat(bid[0], 64)
		quantity, _ = strconv.ParseFloat(bid[1], 64)
		bids += price * quantity
	}
	for _, ask := range order.Asks {
		price, _ = strconv.ParseFloat(ask[0], 64)
		quantity, _ = strconv.ParseFloat(ask[1], 64)
		asks += price * quantity
	}

	return entity.Notional{
		Symbol: symbol,
		Bid:    fmt.Sprintf("%.8f", bids),
		Ask:    fmt.Sprintf("%.8f", asks),
	}, nil

}

func contains(s []string, x string) bool {
	for _, i := range s {
		if i == x {
			return true
		}
	}
	return false
}
