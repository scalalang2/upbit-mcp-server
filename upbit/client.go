package upbit

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const BaseURL = "https://api.upbit.com/v1/"

type Client struct {
	AccessKey  string
	SecretKey  string
	HttpClient *http.Client
}

// NewClient 업비트 클라이언트 생성
func NewClient(accessKey, secretKey string) *Client {
	return &Client{
		AccessKey: accessKey,
		SecretKey: secretKey,
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// structToMap 구조체를 map[string]string 으로 변환
func structToMap(item interface{}) map[string]string {
	res := map[string]string{}
	if item == nil {
		return res
	}
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i)
		val := v.Field(i)

		// 값이 비어있으면 스킵
		if val.IsZero() {
			continue
		}

		// json 태그 값을 키로 사용
		tag := field.Tag.Get("json")
		key := strings.Split(tag, ",")[0]
		if key == "" {
			key = field.Name
		}

		// 값 변환
		var strVal string
		switch val.Kind() {
		case reflect.Int, reflect.Int64:
			strVal = strconv.FormatInt(val.Int(), 10)
		case reflect.Float64:
			strVal = strconv.FormatFloat(val.Float(), 'f', -1, 64)
		case reflect.Bool:
			strVal = strconv.FormatBool(val.Bool())
		default:
			strVal = fmt.Sprintf("%v", val.Interface())
		}
		res[key] = strVal
	}
	return res
}

// generateQueryString: 맵을 정렬된 쿼리 스트링으로 변환
func generateQueryString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Upbit는 파라미터 정렬을 권장함

	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	return strings.Join(parts, "&")
}

// generateToken: JWT 토큰 생성
func (c *Client) generateToken(queryParams map[string]string) (string, error) {
	claims := jwt.MapClaims{
		"access_key": c.AccessKey,
		"nonce":      uuid.New().String(),
	}

	if len(queryParams) > 0 {
		queryString := generateQueryString(queryParams)
		hash := sha512.Sum512([]byte(queryString))
		hashHex := hex.EncodeToString(hash[:])

		claims["query_hash"] = hashHex
		claims["query_hash_alg"] = "SHA512"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.SecretKey))
}

// doRequest: 실제 요청 수행
func (c *Client) doRequest(method, endpoint string, params interface{}, result interface{}) error {
	paramMap := structToMap(params)

	var body io.Reader
	urlString := BaseURL + endpoint

	// GET/DELETE는 쿼리 스트링에 파라미터 추가
	if method == http.MethodGet || method == http.MethodDelete {
		if len(paramMap) > 0 {
			q := url.Values{}
			for k, v := range paramMap {
				q.Add(k, v)
			}
			urlString += "?" + q.Encode()
		}
	} else {
		// POST는 JSON Body 사용
		if params != nil {
			jsonBytes, err := json.Marshal(paramMap)
			if err != nil {
				return err
			}
			body = bytes.NewBuffer(jsonBytes)
		}
	}

	req, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	// 인증 토큰 추가 (Auth가 필요한 경우)
	// paramMap은 Hash 생성을 위해 사용됨
	token, err := c.generateToken(paramMap)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Result가 nil이 아니고 포인터일 때만 언마샬링
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("json unmarshal error: %w, body: %s", err, string(respBody))
		}
	}

	return nil
}

// doNonAuthRequest: 인증이 필요 없는 요청 (시세 조회 등)
func (c *Client) doNonAuthRequest(endpoint string, params interface{}, result interface{}) error {
	paramMap := structToMap(params)
	urlString := BaseURL + endpoint

	if len(paramMap) > 0 {
		q := url.Values{}
		for k, v := range paramMap {
			q.Add(k, v)
		}
		urlString += "?" + q.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

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
	PriceUnit interface{} `json:"price_unit"` // object in C#
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

// --- API Implementation ---

// GetAccounts: 전체 계좌 조회
func (c *Client) GetAccounts() ([]Account, error) {
	var res []Account
	err := c.doRequest("GET", "accounts", nil, &res)
	return res, err
}

// GetOrderHistory: 완료된 주문 조회 (C#의 ClosedOrderHistory)
func (c *Client) GetOrderHistory(params RequestParams) ([]Order, error) {
	var res []Order
	// C#에서는 orders/closed 엔드포인트를 호출하지만 Upbit 공식 API V1은 "orders"에 state=done을 보냄
	// C# 코드를 그대로 따르자면:
	err := c.doRequest("GET", "orders/closed", params, &res)
	// 만약 API가 404를 낸다면 "orders"로 변경하고 params.State = "done" 설정 필요
	return res, err
}

// GetOrder: 특정 주문 조회
func (c *Client) GetOrder(uuid string) (Order, error) {
	var res Order
	params := RequestParams{Uuid: uuid}
	err := c.doRequest("GET", "order", params, &res)
	return res, err
}

// GetOrders: 주문 리스트 조회 (Open, Wait 등)
func (c *Client) GetOrders(params RequestParams) ([]Order, error) {
	var res []Order
	err := c.doRequest("GET", "orders", params, &res)
	return res, err
}

// CancelOrder: 주문 취소
func (c *Client) CancelOrder(uuid string) (bool, error) {
	var res Order
	params := RequestParams{Uuid: uuid}
	err := c.doRequest("DELETE", "order", params, &res)
	// 성공하면 uuid가 담긴 객체가 옴
	return err == nil && res.Uuid != "", err
}

// PlaceOrder: 주문하기
func (c *Client) PlaceOrder(params RequestParams) (Order, error) {
	var res Order
	err := c.doRequest("POST", "orders", params, &res)
	return res, err
}

// GetChance: 주문 가능 정보 확인
func (c *Client) GetChance(market string) (Chance, error) {
	var res Chance
	params := RequestParams{Market: market}
	err := c.doRequest("GET", "orders/chance", params, &res)
	return res, err
}

// GetCoinAddresses: 전체 입금 주소 조회
func (c *Client) GetCoinAddresses() ([]CoinAddress, error) {
	var res []CoinAddress
	err := c.doRequest("GET", "deposits/coin_addresses", nil, &res)
	return res, err
}

// GetCoinAddress: 특정 코인 입금 주소 조회
func (c *Client) GetCoinAddress(currency string) (CoinAddress, error) {
	var res CoinAddress
	params := RequestParams{Currency: currency}
	err := c.doRequest("GET", "deposits/coin_address", params, &res)
	return res, err
}

// GenerateCoinAddress: 입금 주소 생성 요청
func (c *Client) GenerateCoinAddress(currency string) (GenerateCoinAddressResponse, error) {
	var res GenerateCoinAddressResponse
	params := RequestParams{Currency: currency}
	err := c.doRequest("POST", "deposits/generate_coin_address", params, &res)
	return res, err
}

// GetWithdraws: 출금 리스트 조회
func (c *Client) GetWithdraws(params RequestParams) ([]Deposit, error) {
	var res []Deposit
	err := c.doRequest("GET", "withdraws", params, &res) // C# 코드는 deposits를 호출하지만 이름이 GetWithdraws임. API 스펙상 출금은 /withdraws
	// 주의: 원본 C# 코드의 GetWithdraws는 endpoint가 "deposits"로 되어있습니다.
	// 의도대로라면 "withdraws"여야 하나, C# 코드를 그대로 포팅하려면 "deposits"를 써야 합니다.
	// 여기서는 올바른 API인 "withdraws"를 쓰거나 C#대로 "deposits"를 쓸지 결정해야 합니다.
	// 위 코드는 메서드 명에 맞춰 수정됨. 원본 C# 동작을 원하면 "deposits"로 변경하세요.
	return res, err
}

// GetWithdraw: 개별 출금 조회
func (c *Client) GetWithdraw(uuid string) (Deposit, error) {
	var res Deposit
	params := RequestParams{Uuid: uuid}
	err := c.doRequest("GET", "withdraw", params, &res)
	return res, err
}

// DepositKrw: 원화 입금하기
func (c *Client) DepositKrw(amount string) (Deposit, error) {
	var res Deposit
	params := RequestParams{Amount: amount}
	err := c.doRequest("POST", "deposits/krw", params, &res)
	return res, err
}

// GetWalletStatus: 지갑 상태 조회
func (c *Client) GetWalletStatus() ([]WalletStatus, error) {
	var res []WalletStatus
	err := c.doRequest("GET", "status/wallet", nil, &res)
	return res, err
}

// GetApiKeys: API 키 리스트 조회
func (c *Client) GetApiKeys() ([]ApiKey, error) {
	var res []ApiKey
	err := c.doRequest("GET", "api_keys", nil, &res)
	return res, err
}

// --- Quotation API (Non-Auth) ---

// GetMarketCodes: 마켓 코드 조회
func (c *Client) GetMarketCodes() ([]MarketCode, error) {
	var res []MarketCode
	err := c.doNonAuthRequest("market/all", map[string]string{"isDetails": "true"}, &res)
	return res, err
}

// GetTicks: 최근 체결 내역
func (c *Client) GetTicks(params RequestParams) ([]Tick, error) {
	var res []Tick
	err := c.doNonAuthRequest("trades/ticks", params, &res)
	return res, err
}

// GetTicker: 현재가 정보
func (c *Client) GetTicker(symbol string) ([]Ticker, error) {
	var res []Ticker
	// endpoint pattern: ticker?markets=KRW-BTC
	err := c.doNonAuthRequest("ticker", map[string]string{"markets": symbol}, &res)
	return res, err
}

// GetOrderBooks: 호가 정보
func (c *Client) GetOrderBooks(symbol string) ([]OrderBook, error) {
	var res []OrderBook
	err := c.doNonAuthRequest("orderbook", map[string]string{"markets": symbol}, &res)
	return res, err
}

// GetDayCandles: 일봉
func (c *Client) GetDayCandles(params RequestParams) ([]Candle, error) {
	var res []Candle
	// C# 코드: candles/days
	// Path Parameter 처리가 필요할 수 있으나 Upbit는 보통 쿼리파라미터나 고정 path 사용
	err := c.doNonAuthRequest("candles/days", params, &res)
	return res, err
}

// GetWeekCandles: 주봉
func (c *Client) GetWeekCandles(params RequestParams) ([]Candle, error) {
	var res []Candle
	err := c.doNonAuthRequest("candles/weeks", params, &res)
	return res, err
}

// GetMonthCandles: 월봉
func (c *Client) GetMonthCandles(params RequestParams) ([]Candle, error) {
	var res []Candle
	err := c.doNonAuthRequest("candles/months", params, &res)
	return res, err
}

// GetMinuteCandles: 분봉
func (c *Client) GetMinuteCandles(unit int, params RequestParams) ([]Candle, error) {
	var res []Candle
	endpoint := fmt.Sprintf("candles/minutes/%d", unit)
	err := c.doNonAuthRequest(endpoint, params, &res)
	return res, err
}
