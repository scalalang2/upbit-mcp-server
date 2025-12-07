package main

import (
	"context"
	"upbit-mcp-server/upbit"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetAccountsResult struct {
	Accounts []upbit.Account `json:"accounts"`
}

func AddUpbitTools(server *mcp.Server, client *upbit.Client) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "GetAccounts",
			Description: "전체 계좌 조회"},
		func(ctx context.Context, req *mcp.CallToolRequest, params *upbit.RequestParams) (
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
}
