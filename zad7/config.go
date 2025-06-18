package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var cityData = make(map[string]Cities)
var extremeThresholds ExtremeWeatherThresholds

func loadCityData(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("nie udało się otworzyć pliku worldcities.csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err = reader.Read(); err != nil {
		return fmt.Errorf("nie udało się odczytać nagłówka z pliku worldcities.csv: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("błąd podczas czytania rekordu z pliku worldcities.csv: %w", err)
		}
		if len(record) < 5 {
			continue
		}

		city := Cities{
			city:       record[0],
			city_ascii: record[1],
			lat:        record[2],
			lng:        record[3],
			country:    record[4],
		}
		cityData[strings.ToLower(city.city_ascii)] = city
	}
	return nil
}

func getCityCoordinates(cityName string) (string, string, error) {
	city, found := cityData[strings.ToLower(cityName)]
	if !found {
		return "", "", fmt.Errorf("nie znaleziono miasta: %s", cityName)
	}
	if city.country != "Poland" {
		return "", "", fmt.Errorf("miasto %s nie znajduje się w Polsce", cityName)
	}
	return city.lat, city.lng, nil
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("nie udało się wczytać pliku konfiguracyjnego: %w", err)
	}

	if err := viper.UnmarshalKey("extreme_weather_thresholds", &extremeThresholds); err != nil {
		return fmt.Errorf("nie udało się zdekodować progów ekstremalnych zjawisk pogodowych: %w", err)
	}
	return nil
}
