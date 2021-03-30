package adapter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dionisiusrizky/binance-test/entity"
)

func FetchSymbols(quoteAsset string) ([]string, error) {
	//Fetching all symbols
	url := "https://api.binance.com/api/v3/exchangeInfo"
	var exci entity.ExchangeInfo
	var result []string

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&exci)
	if err != nil {
		return nil, err
	}

	//Filtering symbols by quoteAsset
	for _, symbol := range exci.Symbols {
		if symbol.QuoteAsset == quoteAsset {
			result = append(result, symbol.Name)
		}
	}

	return result, nil
}

func Fetch24HRPrice() ([]entity.Price, error) {
	//Fetching ticker
	url := "https://api.binance.com/api/v3/ticker/24hr"
	var prices []entity.Price

	resp, err := http.Get(url)
	if err != nil {
		return prices, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return prices, err
	}

	err = json.Unmarshal(body, &prices)
	if err != nil {
		return prices, err
	}

	return prices, nil
}

func FetchOrderBook(symbol string, limit int) (entity.Order, error) {
	var limitParam string
	var result entity.Order

	switch {
	case limit < 5:
		limitParam = "5"
	case limit < 10:
		limitParam = "10"
	case limit < 20:
		limitParam = "20"
	case limit < 50:
		limitParam = "50"
	case limit < 100:
		limitParam = "100"
	case limit < 500:
		limitParam = "500"
	case limit < 1000:
		limitParam = "1000"
	case limit < 5000:
		limitParam = "5000"
	}

	url := fmt.Sprintf("https://api.binance.com/api/v3/depth?symbol=%s&limit=%s", symbol, limitParam)
	var order entity.Order

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &order)
	if err != nil {
		return result, err
	}

	//If both bids and asks length is larger than limit
	if len(order.Bids) >= limit && len(order.Asks) >= limit {
		result = entity.Order{
			Bids: order.Bids[:limit],
			Asks: order.Asks[:limit],
		}
	} else if len(order.Bids) >= limit { //if only bids is larger than limit
		result = entity.Order{
			Bids: order.Bids[:limit],
			Asks: order.Asks,
		}
	} else if len(order.Asks) >= limit { //if only asks is larger than limit
		result = entity.Order{
			Bids: order.Bids,
			Asks: order.Asks[:limit],
		}
	} else {
		result = entity.Order{
			Bids: order.Bids,
			Asks: order.Asks,
		}
	}

	return result, nil
}
