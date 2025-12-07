package upbit

import (
	"fmt"
	"net/http"
	"time"
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

// GetAccounts: 전체 계좌 조회
func (c *Client) GetAccounts() ([]Account, error) {
	var res []Account
	err := c.doRequest("GET", "accounts", nil, &res)
	return res, err
}

// GetOrderHistory: 완료된 주문 조회
func (c *Client) GetOrderHistory(params RequestParams) ([]Order, error) {
	var res []Order
	err := c.doRequest("GET", "orders/closed", params, &res)
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
	err := c.doRequest("GET", "withdraws", params, &res)
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
