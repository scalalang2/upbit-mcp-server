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

	// Add MCP tools
	AddUpbitTools(server, client)

	log.Println("MCP server started")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
