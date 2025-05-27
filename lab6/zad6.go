package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type NVIDIA struct {
	date      time.Time
	closeLast float64
	volume    int64
	open      float64
	high      float64
	low       float64
}

func parseDollarsToFloat(str string) float64 {

	cleanStr := strings.Replace(str, "$", "", 1)

	floatVal, err := strconv.ParseFloat(cleanStr, 64)
	if err != nil {
		fmt.Println("Error parsing float")
		return 0
	}
	return floatVal
}

func readCSVLineByLine(filepath string) []NVIDIA {
	var historical_data []NVIDIA

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	_, err = reader.Read()

	for i := 0; i < 249; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		date := record[0]
		layout := "01/02/2006"
		parsedTime, err := time.Parse(layout, date)
		if err != nil {
			fmt.Println("Error parsing date:", date)
			return []NVIDIA{}
		}

		closeLast := parseDollarsToFloat(record[1])
		volume, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			fmt.Println("Error parsing volume")
			return []NVIDIA{}
		}
		open := parseDollarsToFloat(record[3])
		high := parseDollarsToFloat(record[4])
		low := parseDollarsToFloat(record[5])

		historical_data = append(historical_data, NVIDIA{parsedTime, closeLast, volume, open, high, low})
	}
	return historical_data
}

func getClose(data []NVIDIA) []float64 {
	var closeArr []float64
	for _, i := range data {
		closeArr = append(closeArr, i.closeLast)
	}
	return closeArr
}

func calculateEMA(prices []float64, period int) []float64 {
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

func calculateATR(data []NVIDIA, period int) []float64 {
	if len(data) < period {
		return nil
	}

	tr := make([]float64, len(data))

	for i := 1; i < len(data); i++ {
		high := data[i].high
		low := data[i].low
		prevClose := data[i-1].closeLast

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

func calculateRSI(prices []float64, period int) []float64 {
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

		avgGain = (avgGain*(float64(period-1)) + gain) / float64(period)
		avgLoss = (avgLoss*(float64(period-1)) + loss) / float64(period)

		if avgLoss == 0 {
			rsi = append(rsi, 100.0)
		} else {
			rs := avgGain / avgLoss
			rsi = append(rsi, 100.0-(100.0/(1+rs)))
		}
	}

	return rsi
}

func getMax(data []float64) float64 {
	max := data[0]
	for _, value := range data {
		if value > max {
			max = value
		}
	}
	return max
}

func getMin(data []float64) float64 {
	min := data[0]
	for _, value := range data {
		if value < min {
			min = value
		}
	}
	return min
}

func main() {
	data := readCSVLineByLine("./NVIDIA_data.csv")
	closeArr := getClose(data)

	ema := calculateEMA(closeArr, 50)
	atr := calculateATR(data, 14)
	rsi := calculateRSI(closeArr, 14)

	fmt.Println("Witamy w analizie działalności spółki NVIDIA na giełdzie w roku 2024")
	fmt.Println("Wybierz dane do wyświetlenia:")
	fmt.Println("1 - Wykładnicza średnia ruchoma (EMA)")
	fmt.Println("2 - Average True Range (ATR)")
	fmt.Println("3 - Wskaźnik siły względnej (RSI)")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Wpisz numer opcji: ")
	scanner.Scan()
	choice := scanner.Text()

	switch choice {
	case "1":
		fmt.Println("\nWykładnicza średnia ruchoma:")
		for _, i := range ema {
			fmt.Println(i)
		}
		fmt.Println("Maksymalna wartość:", getMax(ema))
		fmt.Println("Minimalna wartość:", getMin(ema))
	case "2":
		fmt.Println("\nAverage True Range:")
		for _, i := range atr {
			fmt.Println(i)
		}
		fmt.Println("Maksymalna wartość:", getMax(atr))
		fmt.Println("Minimalna wartość:", getMin(atr))
	case "3":
		fmt.Println("\nWskaźnik siły względnej (RSI):")
		for _, i := range rsi {
			fmt.Println(i)
		}
		fmt.Println("Maksymalna wartość:", getMax(rsi))
		fmt.Println("Minimalna wartość:", getMin(rsi))
	default:
		fmt.Println("Nieprawidłowy wybór.")
	}
}
