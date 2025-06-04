package main

import (
	"bufio"
	"fmt"
	"os"

	"lab6/utils"
)

func main() {
	data := utils.ReadCSV("./data/NVIDIA_data.csv")
	closeArr := utils.GetClosePrices(data)

	ema := utils.CalculateEMA(closeArr, 50)
	atr := utils.CalculateATR(data, 14)
	rsi := utils.CalculateRSI(closeArr, 14)

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
		fmt.Println("Maksymalna wartość:", utils.GetMax(ema))
		fmt.Println("Minimalna wartość:", utils.GetMin(ema))
	case "2":
		fmt.Println("\nAverage True Range:")
		for _, i := range atr {
			fmt.Println(i)
		}
		fmt.Println("Maksymalna wartość:", utils.GetMax(atr))
		fmt.Println("Minimalna wartość:", utils.GetMin(atr))
	case "3":
		fmt.Println("\nWskaźnik siły względnej (RSI):")
		for _, i := range rsi {
			fmt.Println(i)
		}
		fmt.Println("Maksymalna wartość:", utils.GetMax(rsi))
		fmt.Println("Minimalna wartość:", utils.GetMin(rsi))
	default:
		fmt.Println("Nieprawidłowy wybór.")
	}
}
