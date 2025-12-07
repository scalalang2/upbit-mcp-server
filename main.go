package main

import (
	"context"
	"log"
	"os"
	"upbit-mcp-server/upbit"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "greeter",
		Version: "v1.0.0",
	}, nil)

	server.AddReceivingMiddleware(createLoggingMiddleware())

	accessKey := os.Getenv("UPBIT_ACCESS_KEY")
	secretKey := os.Getenv("UPBIT_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		log.Fatal("UPBIT_ACCESS_KEY and UPBIT_SECRET_KEY must be set")
	}

	client := upbit.NewClient(accessKey, secretKey)
	ctx := context.WithValue(context.Background(), upbitClientKey{}, client)

	// Add MCP tools
	mcp.AddTool(server, &mcp.Tool{Name: "GetAccounts", Description: "전체 계좌 조회"}, GetAccounts)
	mcp.AddTool(server, &mcp.Tool{Name: "PlaceBuyOrderByLimit", Description: "지정가 매수 주문하기"}, PlaceBuyOrderByLimit)
	mcp.AddTool(server, &mcp.Tool{Name: "PlaceBuyOrderByMarket", Description: "시장가 매수 주문하기"}, PlaceBuyOrderByMarket)
	mcp.AddTool(server, &mcp.Tool{Name: "PlaceSellOrderByLimit", Description: "지정가 매도 주문하기"}, PlaceSellOrderByLimit)
	mcp.AddTool(server, &mcp.Tool{Name: "PlaceSellOrderByMarket", Description: "시장가 매도 주문하기"}, PlaceSellOrderByMarket)
	mcp.AddTool(server, &mcp.Tool{Name: "CancelOrder", Description: "주문 취소하기"}, CancelOrder)

	mcp.AddTool(server, &mcp.Tool{
		Name: "GetAvailableOrderInfo",
		Description: `Retrieves the order availability information for the specified pair. 
				The response doesn't include current trading pair prices 
				you should consider the current price if you want to decide whether to buy or sell.`,
	}, GetAvailableOrderInfo)

	mcp.AddTool(server, &mcp.Tool{
		Name: "GetClosedOrderHistory",
		Description: `Retrieves the order availability information for the specified pair. 
				The response doesn't include current trading pair prices 
				you should consider the current price if you want to decide whether to buy or sell.`,
	}, GetClosedOrderHistory)

	mcp.AddTool(server, &mcp.Tool{
		Name: "GetOpenOrders",
		Description: `Retrieves the order availability information for the specified pair. 
				The response doesn't include current trading pair prices 
				you should consider the current price if you want to decide whether to buy or sell.`,
	}, GetOpenOrders)

	log.Println("MCP server started")
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
