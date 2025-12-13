package indicators

import (
	"math"
	"upbit-mcp-server/upbit"
)

// CalculateSMA calculates the Simple Moving Average (SMA) for a given period.
func CalculateSMA(candles []*upbit.Candle, period int) []float64 {
	if len(candles) < period {
		return []float64{}
	}
	var smaValues []float64
	for i := period - 1; i < len(candles); i++ {
		sum := 0.0
		for j := i; j > i-period; j-- {
			sum += candles[j].TradePrice
		}
		smaValues = append(smaValues, sum/float64(period))
	}
	return smaValues
}

// CalculateEMA calculates the Exponential Moving Average (EMA) for a given period.
func CalculateEMA(candles []*upbit.Candle, period int) []float64 {
	if len(candles) < period {
		return []float64{}
	}
	var emaValues []float64
	multiplier := 2.0 / (float64(period) + 1.0)

	// Calculate initial SMA for the first EMA value
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += candles[i].TradePrice
	}
	emaValues = append(emaValues, sum/float64(period))

	// Calculate subsequent EMA values
	for i := period; i < len(candles); i++ {
		ema := (candles[i].TradePrice-emaValues[len(emaValues)-1])*multiplier + emaValues[len(emaValues)-1]
		emaValues = append(emaValues, ema)
	}

	return emaValues
}

// CalculateMACD calculates the Moving Average Convergence Divergence (MACD).
func CalculateMACD(candles []*upbit.Candle, shortPeriod, longPeriod, signalPeriod int) ([]float64, []float64, []float64) {
	if len(candles) < longPeriod {
		return []float64{}, []float64{}, []float64{}
	}

	emaShort := CalculateEMA(candles, shortPeriod)
	emaLong := CalculateEMA(candles, longPeriod)

	// Align EMA slices
	emaShort = emaShort[longPeriod-shortPeriod:]

	var macdLine []float64
	for i := 0; i < len(emaLong); i++ {
		macdLine = append(macdLine, emaShort[i]-emaLong[i])
	}

	// Create a temporary slice of "candles" for the MACD line to calculate the signal line
	var macdCandles []*upbit.Candle
	for _, v := range macdLine {
		macdCandles = append(macdCandles, &upbit.Candle{TradePrice: v})
	}

	signalLine := CalculateEMA(macdCandles, signalPeriod)

	// Align signal line with macd line
	macdLine = macdLine[len(macdLine)-len(signalLine):]

	var histogram []float64
	for i := 0; i < len(signalLine); i++ {
		histogram = append(histogram, macdLine[i]-signalLine[i])
	}

	return macdLine, signalLine, histogram
}

// CalculateBollingerBands calculates the Bollinger Bands.
func CalculateBollingerBands(candles []*upbit.Candle, period int, stdDev float64) ([]float64, []float64, []float64) {
	if len(candles) < period {
		return nil, nil, nil
	}

	sma := CalculateSMA(candles, period)
	var upperBand, lowerBand []float64

	for i := period - 1; i < len(candles); i++ {
		sum := 0.0
		for j := i; j > i-period; j-- {
			sum += candles[j].TradePrice
		}
		mean := sum / float64(period)
		sd := 0.0
		for j := i; j > i-period; j-- {
			sd += math.Pow(candles[j].TradePrice-mean, 2)
		}
		sd = math.Sqrt(sd / float64(period))
		upperBand = append(upperBand, mean+sd*stdDev)
		lowerBand = append(lowerBand, mean-sd*stdDev)
	}

	return sma, upperBand, lowerBand
}

// CalculateRSI calculates the Relative Strength Index (RSI).
func CalculateRSI(candles []*upbit.Candle, period int) []float64 {
	if len(candles) < period+1 {
		return []float64{}
	}

	var rsiValues []float64
	var gains, losses []float64

	for i := 1; i < len(candles); i++ {
		change := candles[i].TradePrice - candles[i-1].TradePrice
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, -change)
		}
	}

	avgGain := 0.0
	avgLoss := 0.0
	for i := 0; i < period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	rsiValues = append(rsiValues, rsi)

	for i := period; i < len(gains); i++ {
		avgGain = (avgGain*float64(period-1) + gains[i]) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + losses[i]) / float64(period)

		if avgLoss == 0 {
			rsiValues = append(rsiValues, 100)
		} else {
			rs = avgGain / avgLoss
			rsi = 100 - (100 / (1 + rs))
			rsiValues = append(rsiValues, rsi)
		}
	}
	return rsiValues
}

// CalculateOBV calculates the On-Balance Volume (OBV).
func CalculateOBV(candles []*upbit.Candle) []float64 {
	if len(candles) == 0 {
		return []float64{}
	}

	obvValues := make([]float64, len(candles))
	obvValues[0] = candles[0].CandleAccTradeVolume

	for i := 1; i < len(candles); i++ {
		if candles[i].TradePrice > candles[i-1].TradePrice {
			obvValues[i] = obvValues[i-1] + candles[i].CandleAccTradeVolume
		} else if candles[i].TradePrice < candles[i-1].TradePrice {
			obvValues[i] = obvValues[i-1] - candles[i].CandleAccTradeVolume
		} else {
			obvValues[i] = obvValues[i-1]
		}
	}

	return obvValues
}
