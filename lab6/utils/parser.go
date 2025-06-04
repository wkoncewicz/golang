package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"lab6/models"
)

func ParseDollarsToFloat(str string) float64 {
	cleanStr := strings.Replace(str, "$", "", 1)
	floatVal, err := strconv.ParseFloat(cleanStr, 64)
	if err != nil {
		fmt.Println("Error parsing float")
		return 0
	}
	return floatVal
}

func ReadCSV(filepath string) []models.NVIDIA {
	var historicalData []models.NVIDIA

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
			return []models.NVIDIA{}
		}

		closeLast := ParseDollarsToFloat(record[1])
		volume, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			fmt.Println("Error parsing volume")
			return []models.NVIDIA{}
		}
		open := ParseDollarsToFloat(record[3])
		high := ParseDollarsToFloat(record[4])
		low := ParseDollarsToFloat(record[5])

		historicalData = append(historicalData, models.NVIDIA{
			Date:      parsedTime,
			CloseLast: closeLast,
			Volume:    volume,
			Open:      open,
			High:      high,
			Low:       low,
		})
	}
	return historicalData
}

func GetClosePrices(data []models.NVIDIA) []float64 {
	var closeArr []float64
	for _, i := range data {
		closeArr = append(closeArr, i.CloseLast)
	}
	return closeArr
}
