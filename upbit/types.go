package upbit

type RequestParams struct {
	Market              string `json:"market,omitempty"`
	State               string `json:"state,omitempty"`
	Page                int    `json:"page,omitempty"`
	Limit               int    `json:"limit,omitempty"`
	OrderBy             string `json:"order_by,omitempty"`
	Uuid                string `json:"uuid,omitempty"`
	Identifier          string `json:"identifier,omitempty"`
	Side                string `json:"side,omitempty"`
	Volume              string `json:"volume,omitempty"`
	Price               string `json:"price,omitempty"`
	OrdType             string `json:"ord_type,omitempty"`
	Currency            string `json:"currency,omitempty"`
	Txid                string `json:"txid,omitempty"`
	Amount              string `json:"amount,omitempty"`
	To                  string `json:"to,omitempty"`
	Count               int    `json:"count,omitempty"`
	Cursor              string `json:"cursor,omitempty"`
	DaysAgo             int    `json:"daysAgo,omitempty"`
	Unit                int    `json:"unit,omitempty"`
	ConvertingPriceUnit string `json:"convertingPriceUnit,omitempty"`
}

type Account struct {
	Currency            string `json:"currency"`
	Balance             string `json:"balance"`
	Locked              string `json:"locked"`
	AvgBuyPrice         string `json:"avg_buy_price"`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified"`
	UnitCurrency        string `json:"unit_currency"`
}

type Order struct {
	Uuid            string  `json:"uuid"`
	Side            string  `json:"side"`
	OrdType         string  `json:"ord_type"`
	Price           string  `json:"price"`
	State           string  `json:"state"`
	Market          string  `json:"market"`
	CreatedAt       string  `json:"created_at"`
	Volume          string  `json:"volume"`
	RemainingVolume string  `json:"remaining_volume"`
	ReservedFee     string  `json:"reserved_fee"`
	RemainingFee    string  `json:"remaining_fee"`
	PaidFee         string  `json:"paid_fee"`
	Locked          string  `json:"locked"`
	ExecutedVolume  string  `json:"executed_volume"`
	TradesCount     int     `json:"trades_count"`
	Trades          []Trade `json:"trades,omitempty"` // 상세 조회 시에만 존재
}

type Trade struct {
	Market string `json:"market"`
	Uuid   string `json:"uuid"`
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Funds  string `json:"funds"`
	Side   string `json:"side"`
}

type ApiKey struct {
	AccessKey string `json:"access_key"`
	ExpireAt  string `json:"expire_at"`
}

type WalletStatus struct {
	Currency       string `json:"currency"`
	WalletState    string `json:"wallet_state"`
	BlockState     string `json:"block_state"`
	BlockHeight    int    `json:"block_height"`
	BlockUpdatedAt string `json:"block_updated_at"`
}

type Deposit struct {
	Type            string `json:"type"`
	Uuid            string `json:"uuid"`
	Currency        string `json:"currency"`
	Txid            string `json:"txid"`
	State           string `json:"state"`
	CreatedAt       string `json:"created_at"`
	DoneAt          string `json:"done_at"`
	Amount          string `json:"amount"`
	Fee             string `json:"fee"`
	TransactionType string `json:"transaction_type"`
}

type CoinAddress struct {
	Currency         string `json:"currency"`
	DepositAddress   string `json:"deposit_address"`
	SecondaryAddress string `json:"secondary_address"`
}

type GenerateCoinAddressResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Ticker struct {
	Market             string  `json:"market"`
	TradeDate          string  `json:"trade_date"`
	TradeTime          string  `json:"trade_time"`
	TradeDateKst       string  `json:"trade_date_kst"`
	TradeTimeKst       string  `json:"trade_time_kst"`
	TradeTimestamp     int64   `json:"trade_timestamp"`
	OpeningPrice       float64 `json:"opening_price"`
	HighPrice          float64 `json:"high_price"`
	LowPrice           float64 `json:"low_price"`
	TradePrice         float64 `json:"trade_price"`
	PrevClosingPrice   float64 `json:"prev_closing_price"`
	Change             string  `json:"change"`
	ChangePrice        float64 `json:"change_price"`
	ChangeRate         float64 `json:"change_rate"`
	SignedChangePrice  float64 `json:"signed_change_price"`
	SignedChangeRate   float64 `json:"signed_change_rate"`
	TradeVolume        float64 `json:"trade_volume"`
	AccTradePrice      float64 `json:"acc_trade_price"`
	AccTradePrice24h   float64 `json:"acc_trade_price_24h"`
	AccTradeVolume     float64 `json:"acc_trade_volume"`
	AccTradeVolume24h  float64 `json:"acc_trade_volume_24h"`
	Highest52WeekPrice float64 `json:"highest_52_week_price"`
	Highest52WeekDate  string  `json:"highest_52_week_date"`
	Lowest52WeekPrice  float64 `json:"lowest_52_week_price"`
	Lowest52WeekDate   string  `json:"lowest_52_week_date"`
	Timestamp          int64   `json:"timestamp"`
}

type Candle struct {
	Market               string  `json:"market"`
	CandleDateTimeUtc    string  `json:"candle_date_time_utc"`
	CandleDateTimeKst    string  `json:"candle_date_time_kst"`
	OpeningPrice         float64 `json:"opening_price"`
	HighPrice            float64 `json:"high_price"`
	LowPrice             float64 `json:"low_price"`
	TradePrice           float64 `json:"trade_price"`
	Timestamp            int64   `json:"timestamp"`
	CandleAccTradePrice  float64 `json:"candle_acc_trade_price"`
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	PrevClosingPrice     float64 `json:"prev_closing_price,omitempty"` // 일봉 등에서 사용
	ChangePrice          float64 `json:"change_price,omitempty"`
	ChangeRate           float64 `json:"change_rate,omitempty"`
	Unit                 int     `json:"unit,omitempty"`
}

type MarketCode struct {
	Market        string `json:"market"`
	KoreanName    string `json:"korean_name"`
	EnglishName   string `json:"english_name"`
	MarketWarning string `json:"market_warning"`
}

type OrderBook struct {
	Market         string          `json:"market"`
	Timestamp      int64           `json:"timestamp"`
	TotalAskSize   float64         `json:"total_ask_size"`
	TotalBidSize   float64         `json:"total_bid_size"`
	OrderbookUnits []OrderBookUnit `json:"orderbook_units"`
}

type OrderBookUnit struct {
	AskPrice float64 `json:"ask_price"`
	BidPrice float64 `json:"bid_price"`
	AskSize  float64 `json:"ask_size"`
	BidSize  float64 `json:"bid_size"`
}

type Chance struct {
	BidFee     string        `json:"bid_fee"`
	AskFee     string        `json:"ask_fee"`
	Market     ChanceMarket  `json:"market"`
	BidAccount ChanceAccount `json:"bid_account"`
	AskAccount ChanceAccount `json:"ask_account"`
}

type ChanceMarket struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	OrderTypes []string    `json:"order_types"`
	OrderSides []string    `json:"order_sides"`
	Bid        ChanceLimit `json:"bid"`
	Ask        ChanceLimit `json:"ask"`
	MaxTotal   string      `json:"max_total"`
	State      string      `json:"state"`
}

type ChanceLimit struct {
	Currency  string      `json:"currency"`
	PriceUnit interface{} `json:"price_unit"`
	MinTotal  int         `json:"min_total"`
}

type ChanceAccount struct {
	Currency            string `json:"currency"`
	Balance             string `json:"balance"`
	Locked              string `json:"locked"`
	AvgBuyPrice         string `json:"avg_buy_price"`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified"`
	UnitCurrency        string `json:"unit_currency"`
}

type Tick struct {
	Market           string  `json:"market"`
	TradeDateUtc     string  `json:"trade_date_utc"`
	TradeTimeUtc     string  `json:"trade_time_utc"`
	Timestamp        int64   `json:"timestamp"`
	TradePrice       float64 `json:"trade_price"`
	TradeVolume      float64 `json:"trade_volume"`
	PrevClosingPrice float64 `json:"prev_closing_price"`
	ChangePrice      float64 `json:"change_price"`
	AskBid           string  `json:"ask_bid"`
}
