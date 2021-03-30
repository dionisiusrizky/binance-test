package entity

type ExchangeInfo struct {
	ServerTime int      `json:"serverTime"`
	Symbols    []Symbol `json:"symbols"`
}

type Symbol struct {
	Name       string `json:"symbol"`
	QuoteAsset string `json:"quoteAsset"`
}

type Price struct {
	Symbol   string  `json:"symbol"`
	Volume   float64 `json:"volume,string"`
	Count    int     `json:"count"`
	BidPrice float64 `json:"bidPrice,string"`
	AskPrice float64 `json:"askPrice,string"`
}

type PriceSpread struct {
	Symbol string `json:"symbol"`
	Spread string `json:"spread"`
}

type SpreadDelta struct {
	Spread string `json:"spread"`
	Delta  string `json:"delta"`
}

type Order struct {
	Bids []Bid `json:"bids"`
	Asks []Ask `json:"asks"`
}

type Bid []string
type Ask []string

type Notional struct {
	Symbol string `json:"symbol"`
	Bid    string `json:"bid"`
	Ask    string `json:"ask"`
}
