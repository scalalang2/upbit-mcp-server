package main

import (
	"context"
	"fmt"
	"upbit-mcp-server/indicators"
	"upbit-mcp-server/upbit"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetMovingAverageRequest struct {
	Market string `json:"market" jsonschema:"Trading pair code representing the market (e.g. KRW-BTC, KRW-ETH ...)"`
	Period int    `json:"period" jsonschema:"The period to calculate the moving average for."`
	Count  int    `json:"count,omitempty" jsonschema:"Number of candles to retrieve. Max 200."`
}

type GetMovingAverageResult struct {
	SMA []float64 `json:"sma"`
	EMA []float64 `json:"ema"`
}

type GetMACDRequest struct {
	Market       string `json:"market" jsonschema:"Trading pair code representing the market (e.g. KRW-BTC, KRW-ETH ...)"`
	ShortPeriod  int    `json:"short_period" jsonschema:"The short period for MACD calculation."`
	LongPeriod   int    `json:"long_period" jsonschema:"The long period for MACD calculation."`
	SignalPeriod int    `json:"signal_period" jsonschema:"The signal period for MACD calculation."`
	Count        int    `json:"count,omitempty" jsonschema:"Number of candles to retrieve. Max 200."`
}

type GetMACDResult struct {
	MACDLine   []float64 `json:"macd_line"`
	SignalLine []float64 `json:"signal_line"`
	Histogram  []float64 `json:"histogram"`
}

type GetBollingerBandsRequest struct {
	Market string  `json:"market" jsonschema:"Trading pair code representing the market (e.g. KRW-BTC, KRW-ETH ...)"`
	Period int     `json:"period" jsonschema:"The period to calculate the Bollinger Bands for."`
	StdDev float64 `json:"std_dev" jsonschema:"The standard deviation to use for the Bollinger Bands."`
	Count  int     `json:"count,omitempty" jsonschema:"Number of candles to retrieve. Max 200."`
}

type GetBollingerBandsResult struct {
	SMA       []float64 `json:"sma"`
	UpperBand []float64 `json:"upper_band"`
	LowerBand []float64 `json:"lower_band"`
}

type GetRSIRequest struct {
	Market string `json:"market" jsonschema:"Trading pair code representing the market (e.g. KRW-BTC, KRW-ETH ...)"`
	Period int    `json:"period" jsonschema:"The period to calculate the RSI for."`
	Count  int    `json:"count,omitempty" jsonschema:"Number of candles to retrieve. Max 200."`
}

type GetRSIResult struct {
	RSI []float64 `json:"rsi"`
}

type GetOBVRequest struct {
	Market string `json:"market" jsonschema:"Trading pair code representing the market (e.g. KRW-BTC, KRW-ETH ...)"`
	Count  int    `json:"count,omitempty" jsonschema:"Number of candles to retrieve. Max 200."`
}

type GetOBVResult struct {
	OBV []float64 `json:"obv"`
}

func GetMovingAverage(ctx context.Context, req *mcp.CallToolRequest, params *GetMovingAverageRequest) (*mcp.CallToolResult, *GetMovingAverageResult, error) {
	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	candles, err := client.GetDayCandles(upbit.RequestParams{
		Market: params.Market,
		Count:  params.Count,
	})
	if err != nil {
		return nil, nil, err
	}

	sma := indicators.CalculateSMA(candles, params.Period)
	ema := indicators.CalculateEMA(candles, params.Period)

	return &mcp.CallToolResult{}, &GetMovingAverageResult{SMA: sma, EMA: ema}, nil
}

func GetMACD(ctx context.Context, req *mcp.CallToolRequest, params *GetMACDRequest) (*mcp.CallToolResult, *GetMACDResult, error) {
	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	candles, err := client.GetDayCandles(upbit.RequestParams{
		Market: params.Market,
		Count:  params.Count,
	})
	if err != nil {
		return nil, nil, err
	}

	macd, signal, histogram := indicators.CalculateMACD(candles, params.ShortPeriod, params.LongPeriod, params.SignalPeriod)

	return &mcp.CallToolResult{}, &GetMACDResult{MACDLine: macd, SignalLine: signal, Histogram: histogram}, nil
}

func GetBollingerBands(ctx context.Context, req *mcp.CallToolRequest, params *GetBollingerBandsRequest) (*mcp.CallToolResult, *GetBollingerBandsResult, error) {
	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	candles, err := client.GetDayCandles(upbit.RequestParams{
		Market: params.Market,
		Count:  params.Count,
	})
	if err != nil {
		return nil, nil, err
	}

	sma, upper, lower := indicators.CalculateBollingerBands(candles, params.Period, params.StdDev)

	return &mcp.CallToolResult{}, &GetBollingerBandsResult{SMA: sma, UpperBand: upper, LowerBand: lower}, nil
}

func GetRSI(ctx context.Context, req *mcp.CallToolRequest, params *GetRSIRequest) (*mcp.CallToolResult, *GetRSIResult, error) {
	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	candles, err := client.GetDayCandles(upbit.RequestParams{
		Market: params.Market,
		Count:  params.Count,
	})
	if err != nil {
		return nil, nil, err
	}

	rsi := indicators.CalculateRSI(candles, params.Period)

	return &mcp.CallToolResult{}, &GetRSIResult{RSI: rsi}, nil
}

func GetOBV(ctx context.Context, req *mcp.CallToolRequest, params *GetOBVRequest) (*mcp.CallToolResult, *GetOBVResult, error) {
	client, ok := ctx.Value(upbitClientKey{}).(*upbit.Client)
	if !ok {
		return nil, nil, fmt.Errorf("Upbit client not found in context")
	}

	candles, err := client.GetDayCandles(upbit.RequestParams{
		Market: params.Market,
		Count:  params.Count,
	})
	if err != nil {
		return nil, nil, err
	}

	obv := indicators.CalculateOBV(candles)

	return &mcp.CallToolResult{}, &GetOBVResult{OBV: obv}, nil
}
