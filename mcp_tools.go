package main

import (
	"context"
	"fmt"
	"upbit-mcp-server/upbit"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// upbitClientKey는 context 내에서 Upbit 클라이언트를 식별하기 위한 키
type upbitClientKey struct{}

type GetAccountsResult struct {
	Accounts []upbit.Account `json:"accounts"`
}

type GetAvailableOrderInfoRequest struct {
	Market string `json:"market" jsonschema:"Trading pair code representing the market (e.g. KRW-BTC, KRW-ETH ...)"`
}

type GetClosedOrderHistoryRequest struct {
	Market  string `json:"market" jsonschema:"Trading pair code representing the market (e.g. KRW-BTC, KRW-ETH ...)"`
	State   string `json:"state" jsonschema:"Status of the order (allowed value: 'done', 'cancel')"`
	Limit   int    `json:"limit" jsonschema:"요청 개수 (default: 100, max: 1000)"`
	OrderBy string `json:"order_by" jsonschema:"Sorting method for query results. Returns a list of orders sorted according to the specified method based on the order creation time. The available values are 'desc' (descending, latest orders first) or 'asc' (ascending, oldest orders first). The default value is 'desc'. Allowed: 'asc', 'desc'"`

	// TODO.
	// StartTime string `json:"start_time" jsonschema:"Start time of the query period. Only orders created within the specified time range are returned. Maximum range is 7 days. Can be ISO 8601 format with timezone e.g. 2025-06-24T04:56:53Z, 2025-06-24T13:56:53+09:00"`
	// EndTime   string `json:"end_time" jsonschema:"End time of the query period. Only orders created from 'start_time' up to this time are returned. Maximum range is 7 days. Can be ISO 8601 format with timezone e.g. 2025-06-24T04:56:53Z, 2025-06-24T13:56:53+09:00"`
}

type GetOpenOrderHistoryRequest struct {
	Market  string `json:"market" jsonschema:"Trading pair code representing the market."`
	Page    int    `json:"page" jsonschema:"Page number for pagination. A parameter for pagination that allows you to specify the page to retrieve. If not specified, the default value is 1."`
	Limit   int    `json:"limit" jsonschema:"요청 개수(default: 100, max: 100). 요청 당 조회할 주문 개수를 지정합니다. 한번에 최대 100개의 항목을 조회할 수 있으며, 미지정시 기본값은 100입니다."`
	OrderBy string `json:"order_by" jsonschema:"Sorting method for query results. Returns a list of orders sorted according to the specified method based on the order creation time. The available values are 'desc' (descending, latest orders first) or 'asc' (ascending, oldest orders first). The default value is 'desc'. Allowed: 'asc', 'desc'"`
}

type GetClosedOrderHistoryResult struct {
	Orders []upbit.Order `json:"orders"`
}

type GetOpenOrderHistoryResult struct {
	Orders []upbit.Order `json:"orders"`
}

type PlaceBuyOrderByLimitRequest struct {
	Market string `json:"market" jsonschema:"Trading pair code representing the market."`
	Price  string `json:"price" jsonschema:"Order price based on the quote currency. e.g. when buying 1 BTC at 100,000,000 KRW per BTC in the KRW-BTC market, enter 100000000."`
	Volume string `json:"volume" jsonschema:"Order quantity. e.g. to buy 0.1 BTC in the KRW-BTC market, enter 0.1"`
}

type PlaceBuyOrderByMarketRequest struct {
	Market string `json:"market" jsonschema:"Trading pair code representing the market."`
	Price  string `json:"price" jsonschema:"Total order amount based on the quote currency. For example, entering 100000000 in the KRW-BTC pair will buy BTC worth 100,000,000 KRW at market price."`
}

type PlaceSellOrderByLimitRequest struct {
	Market string `json:"market" jsonschema:"Trading pair code representing the market."`
	Price  string `json:"price" jsonschema:"Order price based on the quote currency. For example, when selling 1 BTC at 100,000,000 KRW per BTC in the KRW-BTC market, enter 100000000."`
	Volume string `json:"volume" jsonschema:"Order quantity e.g. to sell 0.1 BTC in the KRW-BTC market, enter 0.1"`
}

type PlaceSellOrderByMarketRequest struct {
	Market string `json:"market" jsonschema:"Trading pair code representing the market."`
	Volume string `json:"volume" jsonschema:"Sell order quantity. For example, entering 0.1 in the KRW-BTC pair will sell 0.1 BTC at market price"`
}

type CancelOrderRequest struct {
	UUID string `json:"uuid" jsonschema:"Unique identifier (UUID) for the order to cancel."`
}

type CancelOrderResult struct {
	Canceled bool `json:"canceled" jsonschema:"Whether the order is canceled or not"`
}

func GetAccounts(ctx context.Context, req *mcp.CallToolRequest, params any) (
	*mcp.CallToolResult,
	*GetAccountsResult,
	error,
) {
	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	var res mcp.CallToolResult

	accounts, err := client.GetAccounts()
	if err != nil {
		return nil, nil, err
	}

	return &res, &GetAccountsResult{Accounts: accounts}, nil
}

func PlaceBuyOrderByLimit(ctx context.Context, req *mcp.CallToolRequest, params *PlaceBuyOrderByLimitRequest) (
	*mcp.CallToolResult,
	*upbit.Order,
	error,
) {
	var res mcp.CallToolResult

	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	orderResult, err := client.PlaceOrder(upbit.RequestParams{
		Market:  params.Market,
		Side:    "bid",
		OrdType: "limit",
		Price:   params.Price,
		Volume:  params.Volume,
		SmpType: "cancel_maker",
	})
	if err != nil {
		return nil, nil, err
	}

	return &res, &orderResult, nil
}

func PlaceBuyOrderByMarket(ctx context.Context, req *mcp.CallToolRequest, params *PlaceBuyOrderByMarketRequest) (
	*mcp.CallToolResult,
	*upbit.Order,
	error,
) {
	var res mcp.CallToolResult

	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	orderResult, err := client.PlaceOrder(upbit.RequestParams{
		Market:  params.Market,
		Side:    "bid",
		OrdType: "price",
		Price:   params.Price,
		SmpType: "cancel_maker",
	})
	if err != nil {
		return nil, nil, err
	}

	return &res, &orderResult, nil
}

func PlaceSellOrderByLimit(ctx context.Context, req *mcp.CallToolRequest, params *PlaceSellOrderByLimitRequest) (
	*mcp.CallToolResult,
	*upbit.Order,
	error,
) {
	var res mcp.CallToolResult

	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	orderResult, err := client.PlaceOrder(upbit.RequestParams{
		Market:  params.Market,
		Side:    "ask",
		OrdType: "limit",
		Price:   params.Price,
		Volume:  params.Volume,
		SmpType: "cancel_maker",
	})
	if err != nil {
		return nil, nil, err
	}

	return &res, &orderResult, nil
}

func PlaceSellOrderByMarket(ctx context.Context, req *mcp.CallToolRequest, params *PlaceSellOrderByMarketRequest) (
	*mcp.CallToolResult,
	*upbit.Order,
	error,
) {
	var res mcp.CallToolResult

	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	orderResult, err := client.PlaceOrder(upbit.RequestParams{
		Market:  params.Market,
		Side:    "ask",
		OrdType: "market",
		Volume:  params.Volume,
		SmpType: "cancel_maker",
	})
	if err != nil {
		return nil, nil, err
	}

	return &res, &orderResult, nil
}

func CancelOrder(ctx context.Context, req *mcp.CallToolRequest, params *CancelOrderRequest) (
	*mcp.CallToolResult,
	*CancelOrderResult,
	error,
) {
	var res mcp.CallToolResult

	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	canceled, err := client.CancelOrder(params.UUID)
	if err != nil {
		return nil, nil, err
	}

	return &res, &CancelOrderResult{Canceled: canceled}, nil
}

func GetAvailableOrderInfo(ctx context.Context, req *mcp.CallToolRequest, params *GetAvailableOrderInfoRequest) (
	*mcp.CallToolResult,
	*upbit.Chance,
	error,
) {
	var res mcp.CallToolResult

	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	chance, err := client.GetChance(params.Market)
	if err != nil {
		return nil, nil, err
	}

	return &res, &chance, nil
}

// TODO. 지금은 최근 7일 동안의 조회만 가능함
func GetClosedOrderHistory(ctx context.Context, req *mcp.CallToolRequest, params *GetClosedOrderHistoryRequest) (
	*mcp.CallToolResult,
	*GetClosedOrderHistoryResult,
	error,
) {
	var res mcp.CallToolResult

	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	orderHistory, err := client.GetOrderHistory(upbit.RequestParams{
		Market:  params.Market,
		State:   params.State,
		OrderBy: params.OrderBy,
		Limit:   params.Limit,
	})
	if err != nil {
		return nil, nil, err
	}

	return &res, &GetClosedOrderHistoryResult{Orders: orderHistory}, nil
}

func GetOpenOrders(ctx context.Context, req *mcp.CallToolRequest, params *GetOpenOrderHistoryRequest) (
	*mcp.CallToolResult,
	*GetOpenOrderHistoryResult,
	error,
) {
	var res mcp.CallToolResult

	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	orderHistory, err := client.GetOpenOrders(upbit.RequestParams{
		Market:  params.Market,
		Page:    params.Page,
		Limit:   params.Limit,
		OrderBy: params.OrderBy,
	})
	if err != nil {
		return nil, nil, err
	}

	return &res, &GetOpenOrderHistoryResult{Orders: orderHistory}, nil
}
