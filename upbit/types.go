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
	SmpType             string `json:"smp_type,omitempty"`
}

type Account struct {
	Currency            string `json:"currency" jsonschema:"Currency code to be queried"`
	Balance             string `json:"balance" jsonschema:"Available amount or volume for orders. For digital assets, this represents the available quantity. For fiat currency, this represents the available amount"`
	Locked              string `json:"locked" jsonschema:"Amount or quantity locked by pending orders or withdrawals"`
	AvgBuyPrice         string `json:"avg_buy_price" jsonschema:"Average buy price of the asset"`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified" jsonschema:"Indicates whether the average buy price has been modified"`
	UnitCurrency        string `json:"unit_currency" jsonschema:"Currency unit used as the basis for avg_buy_price. [Example] KRW, BTC, USDT"`
}

type Order struct {
	Uuid            string  `json:"uuid" jsonschema:"Unique identifier (UUID) for the order."`
	Side            string  `json:"side" jsonschema:"Order side: ask (sell), bid (buy)."`
	OrdType         string  `json:"ord_type" jsonschema:"Order type to create. (limit: Limit buy/sell order, price: market buy order, market: market sell order)"`
	Price           string  `json:"price" jsonschema:"Order unit price or total amount. For limit orders, this is the unit price. For market buy orders, this is the total purchase amount."`
	State           string  `json:"state" jsonschema:"Order status. (done, cancel)"`
	Market          string  `json:"market" jsonschema:"Trading pair code representing the market"`
	CreatedAt       string  `json:"created_at" jsonschema:"Order creation time in KST [Format] yyyy-MM-ddTHH:mm:ss+09:00"`
	Volume          string  `json:"volume" jsonschema:"Order request amount or quantity."`
	RemainingVolume string  `json:"remaining_volume" jsonschema:"Remaining order quantity after execution."`
	ExecutedVolume  string  `json:"executed_volume" jsonschema:"Executed order quantity."`
	ReservedFee     string  `json:"reserved_fee" jsonschema:"Fee amount reserved for the order."`
	RemainingFee    string  `json:"remaining_fee" jsonschema:"Fee amount reserved for the order."`
	PaidFee         string  `json:"paid_fee" jsonschema:"Fee amount paid at the time of execution."`
	Locked          string  `json:"locked" jsonschema:"Amount or quantity locked by pending orders or trades."`
	TradesCount     int     `json:"trades_count" jsonschema:"Number of trades executed for the order."`
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

type MarketTradingPair struct {
	Market      string        `json:"market"`
	KoreanName  string        `json:"korean_name"`
	EnglishName string        `json:"english_name"`
	MarketEvent []MarketEvent `json:"market_warning"`
}

type MarketEvent struct {
	Warning bool          `json:"warning" jsonschema:"Whether the pair has been designated as an 'Investment Caution'' item under Upbit’s market alert system"`
	Caution MarketCaution `json:"caution"`
}

type MarketCaution struct {
	PriceFluctuations            bool `json:"PRICE_FLUCTUATIONS" jsonschema:"Price Surge/Drop Alert"`
	TradingVolumeSoaring         bool `json:"TRADING_VOLUME_SOARING" jsonschema:"Trading Volume Surge Alert"`
	DepositAmountSoaring         bool `json:"DEPOSIT_AMOUNT_SOARING" jsonschema:"Deposit Volume Surge Alert"`
	GlobalPriceDifferences       bool `json:"GLOBAL_PRICE_DIFFERENCES" jsonschema:"Domestic and International Price Difference Alert"`
	ConcentrationOfSmallAccounts bool `json:"CONCENTRATION_OF_SMALL_ACCOUNTS" jsonschema:"Concentrated Trading by a Small Number of Accounts Alert"`
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
	BidFee      string       `json:"bid_fee" jsonschema:"Fee rate applied to buy orders."`
	AskFee      string       `json:"ask_fee" jsonschema:"Fee rate applied to sell orders"`
	MakerBidFee string       `json:"maker_bid_fee" jsonschema:"Fee rate for buy maker orders."`
	MakerAskFee string       `json:"maker_ask_fee" jsonschema:"Fee rate for sell maker orders."`
	Market      ChanceMarket `json:"market"`
	BidAccount  BidAccount   `json:"bid_account"`
	AskAccount  AskAccount   `json:"ask_account"`
}

type ChanceMarket struct {
	Id         string         `json:"id" jsonschema:"Trading pair code representing the market. e.g. KRW-BTC"`
	Name       string         `json:"name" jsonschema:"Trading pair code in the format (base asset)/(quote asset)."`
	OrderSides []string       `json:"order_sides" jsonschema:"Supported order sides: 'bid' (buy), 'ask' (sell)."`
	BidTypes   []string       `json:"bid_types" jsonschema:"Supported buy order types."`
	AskTypes   []string       `json:"ask_types" jsonschema:"Supported sell order types."`
	Bid        BidChanceLimit `json:"bid" jsonschema:"Bid constraints"`
	Ask        AskChanceLimit `json:"ask" jsonschema:"Ask constraints"`
	MaxTotal   string         `json:"max_total" jsonschema:"Maximum available order amount."`
	State      string         `json:"state" jsonschema:"Trading pair operation status."`
}

type BidChanceLimit struct {
	Currency string `json:"currency" jsonschema:"디지털 자산 구매에 사용되는 통화(KRW,BTC,USDT)"`
	MinTotal string `json:"min_total" jsonschema:"매수 시 최소 주문 금액(결제 화폐 기준) [예시] min_total: 5000일 경우, 5000 KRW를 의미합니다."`
}

type AskChanceLimit struct {
	Currency string `json:"currency" jsonschema:"매도 자산 통화 e.g. BTC, ETH"`
	MinTotal string `json:"min_total" jsonschema:"매도 시 최소 주문 금액 ([예시] min_total: 5000일 경우, 5000 KRW를 의미합니다."`
}

type BidAccount struct {
	Currency            string `json:"currency" jsonschema:"Currency code to be queried."`
	Balance             string `json:"balance" jsonschema:"Available amount or volume for orders. For digital assets, this represents the available quantity. For fiat currency, this represents the available amount."`
	Locked              string `json:"locked" jsonschema:"Amount or quantity locked by pending orders or withdrawals."`
	AvgBuyPrice         string `json:"avg_buy_price" jsonschema:"Average buy price of the asset."`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified" jsonschema:"Indicates whether the average buy price has been modified."`
	UnitCurrency        string `json:"unit_currency" jsonschema:"Currency unit used as the basis for avg_buy_price. [Example] KRW, BTC, USDT"`
}

type AskAccount struct {
	Currency            string `json:"currency" jsonschema:"Currency code to be queried."`
	Balance             string `json:"balance" jsonschema:"Available amount or volume for orders. For digital assets, this represents the available quantity. For fiat currency, this represents the available amount."`
	Locked              string `json:"locked" jsonschema:"Amount or quantity locked by pending orders or withdrawals."`
	AvgBuyPrice         string `json:"avg_buy_price" jsonschema:"Average buy price of the asset."`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified" jsonschema:"Indicates whether the average buy price has been modified."`
	UnitCurrency        string `json:"unit_currency" jsonschema:"Currency unit used as the basis for avg_buy_price. [Example] KRW, BTC, USDT"`
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
