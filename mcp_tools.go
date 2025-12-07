package main

import (
	"context"
	"upbit-mcp-server/upbit"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

func AddUpbitTools(server *mcp.Server, client *upbit.Client) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "GetAccounts",
			Description: "전체 계좌 조회"},
		func(ctx context.Context, req *mcp.CallToolRequest, params any) (
			*mcp.CallToolResult,
			*GetAccountsResult,
			error,
		) {
			var res mcp.CallToolResult

			accounts, err := client.GetAccounts()
			if err != nil {
				return nil, nil, err
			}

			return &res, &GetAccountsResult{Accounts: accounts}, nil
		})

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "PlaceBuyOrder",
			Description: "지정가/시장가 매수 주문하기"},
		func(ctx context.Context, req *mcp.CallToolRequest, params any) (
			*mcp.CallToolResult,
			*upbit.Order,
			error,
		) {
			var res mcp.CallToolResult

			orderResult, err := client.PlaceOrder()
			if err != nil {
				return nil, nil, err
			}

			return &res, &orderResult, nil
		})

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "PlaceSellOrder",
			Description: "지정가/시장가 매도 주문하기"},
		func(ctx context.Context, req *mcp.CallToolRequest, params any) (
			*mcp.CallToolResult,
			*upbit.Order,
			error,
		) {
			var res mcp.CallToolResult

			orderResult, err := client.PlaceOrder()
			if err != nil {
				return nil, nil, err
			}

			return &res, &orderResult, nil
		})

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "CancelOrder",
			Description: "지정가/시장가 매도 주문하기"},
		func(ctx context.Context, req *mcp.CallToolRequest, params any) (
			*mcp.CallToolResult,
			*GetAccountsResult,
			error,
		) {
			var res mcp.CallToolResult

			accounts, err := client.GetAccounts()
			if err != nil {
				return nil, nil, err
			}

			return &res, &GetAccountsResult{Accounts: accounts}, nil
		})

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: "GetAvailableOrderInfo",
			Description: `Retrieves the order availability information for the specified pair. 
				The response doesn't include current trading pair prices 
				you should consider the current price if you want to decide whether to buy or sell.`,
		},
		func(ctx context.Context, req *mcp.CallToolRequest, params GetAvailableOrderInfoRequest) (
			*mcp.CallToolResult,
			*upbit.Chance,
			error,
		) {
			var res mcp.CallToolResult

			chance, err := client.GetChance(params.Market)
			if err != nil {
				return nil, nil, err
			}

			return &res, &chance, nil
		})

	// TODO. 지금은 최근 7일 동안의 조회만 가능함
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: "GetClosedOrderHistory",
			Description: `Retrieves the order availability information for the specified pair. 
				The response doesn't include current trading pair prices 
				you should consider the current price if you want to decide whether to buy or sell.`,
		},
		func(ctx context.Context, req *mcp.CallToolRequest, params GetClosedOrderHistoryRequest) (
			*mcp.CallToolResult,
			*GetClosedOrderHistoryResult,
			error,
		) {
			var res mcp.CallToolResult

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
		})

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: "GetOpenOrders",
			Description: `Retrieves the order availability information for the specified pair. 
				The response doesn't include current trading pair prices 
				you should consider the current price if you want to decide whether to buy or sell.`,
		},
		func(ctx context.Context, req *mcp.CallToolRequest, params GetOpenOrderHistoryRequest) (
			*mcp.CallToolResult,
			*GetOpenOrderHistoryResult,
			error,
		) {
			var res mcp.CallToolResult

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
		})
}
