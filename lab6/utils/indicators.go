package utils

import (
	"math"

	"lab6/models"
)

func CalculateEMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	ema := make([]float64, len(prices))
	k := 2.0 / float64(period+1)

	var sum float64
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i] * k) + (ema[i-1] * (1 - k))
	}
	return ema[period-1:]
}

func CalculateATR(data []models.NVIDIA, period int) []float64 {
	if len(data) < period {
		return nil
	}

	tr := make([]float64, len(data))

	for i := 1; i < len(data); i++ {
		high := data[i].High
		low := data[i].Low
		prevClose := data[i-1].CloseLast

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)

		tr[i] = math.Max(tr1, math.Max(tr2, tr3))
	}

	atr := make([]float64, len(data))
	var sum float64

	for i := 1; i <= period; i++ {
		sum += tr[i]
	}
	atr[period] = sum / float64(period)

	for i := period + 1; i < len(data); i++ {
		atr[i] = (atr[i-1]*(float64(period)-1) + tr[i]) / float64(period)
	}

	return atr[period:]
}

func CalculateRSI(prices []float64, period int) []float64 {
	if len(prices) < period+1 {
		return nil
	}

	rsi := make([]float64, 0, len(prices)-period)
	gainSum := 0.0
	lossSum := 0.0

	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gainSum += change
		} else {
			lossSum -= change
		}
	}

	avgGain := gainSum / float64(period)
	avgLoss := lossSum / float64(period)

	if avgLoss == 0 {
		rsi = append(rsi, 100.0)
	} else {
		rs := avgGain / avgLoss
		rsi = append(rsi, 100.0-(100.0/(1+rs)))
	}

	for i := period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		gain := 0.0
		loss := 0.0

		if change > 0 {
			gain = change
		} else {
			loss = -change
		}

		avgGain = (avgGain*(float64(period)-1) + gain) / float64(period)
		avgLoss = (avgLoss*(float64(period)-1) + loss) / float64(period)

		if avgLoss == 0 {
			rsi = append(rsi, 100.0)
		} else {
			rs := avgGain / avgLoss
			rsi = append(rsi, 100.0-(100.0/(1+rs)))
		}
	}

	return rsi
}
