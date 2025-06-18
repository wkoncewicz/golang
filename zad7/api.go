package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

func getWeatherAPI(lat, lng string) error {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,apparent_temperature,is_day,rain,precipitation", lat, lng)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("błąd podczas wykonywania zapytania HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nieoczekiwany status odpowiedzi: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("błąd podczas odczytu odpowiedzi: %w", err)
	}

	var weatherData CurrentWeatherResponse
	if err = json.Unmarshal(body, &weatherData); err != nil {
		return fmt.Errorf("błąd podczas parsowania JSON: %w", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Parametr", "Wartość"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	dayNight := "Noc"
	if weatherData.Current.IsDay == 1 {
		dayNight = "Dzień"
	}

	table.Append([]string{"Temperatura", fmt.Sprintf("%.2f %s", weatherData.Current.Temperature2M, weatherData.CurrentUnits.Temperature2M)})
	table.Append([]string{"Wilgotność względna", fmt.Sprintf("%d %s", weatherData.Current.RelativeHumidity2M, weatherData.CurrentUnits.RelativeHumidity2M)})
	table.Append([]string{"Temperatura odczuwalna", fmt.Sprintf("%.2f %s", weatherData.Current.ApparentTemperature, weatherData.CurrentUnits.ApparentTemperature)})
	table.Append([]string{"Dzień/Noc", dayNight})
	table.Append([]string{"Deszcz", fmt.Sprintf("%.2f %s", weatherData.Current.Rain, weatherData.CurrentUnits.Rain)})
	table.Append([]string{"Opady", fmt.Sprintf("%.2f %s", weatherData.Current.Precipitation, weatherData.CurrentUnits.Precipitation)})
	table.Render()

	return nil
}

func getForecastWeatherAPI(lat, lng string, days int, cityName string) (PlotData, error) {
	plotData := PlotData{}

	if days > 16 {
		fmt.Println("Ostrzeżenie: prognoza jest dostępna tylko na maksymalnie 16 dni. Ograniczam do 16 dni.")
		days = 16
	}

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=temperature_2m_max,temperature_2m_min,apparent_temperature_max,apparent_temperature_min,sunrise,sunset,precipitation_sum&timezone=Europe%%2FBerlin&forecast_days=%d", lat, lng, days)

	resp, err := http.Get(url)
	if err != nil {
		return plotData, fmt.Errorf("błąd podczas wykonywania zapytania HTTP dla prognozy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return plotData, fmt.Errorf("nieoczekiwany status odpowiedzi dla prognozy: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return plotData, fmt.Errorf("błąd podczas odczytu odpowiedzi dla prognozy: %w", err)
	}

	var forecastData DailyForecastResponse
	if err = json.Unmarshal(body, &forecastData); err != nil {
		return plotData, fmt.Errorf("błąd podczas parsowania JSON dla prognozy: %w", err)
	}

	fmt.Printf("\n--- Prognoza pogody na %d dni dla %s---\n", days, cityName)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Data", "Temp Max", "Temp Min", "Odcz. Max", "Odcz. Min", "Wschód", "Zachód", "Opady Suma"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for i := 0; i < len(forecastData.Daily.Time); i++ {
		t, parseErr := time.Parse("2006-01-02", forecastData.Daily.Time[i])
		if parseErr == nil {
			plotData.Dates = append(plotData.Dates, t)
			plotData.MaxTemps = append(plotData.MaxTemps, forecastData.Daily.Temperature2MMax[i])
			plotData.MinTemps = append(plotData.MinTemps, forecastData.Daily.Temperature2MMin[i])
			plotData.ApparentMaxTemps = append(plotData.ApparentMaxTemps, forecastData.Daily.ApparentTemperatureMax[i])
			plotData.ApparentMinTemps = append(plotData.ApparentMinTemps, forecastData.Daily.ApparentTemperatureMin[i])
		}

		table.Append([]string{
			forecastData.Daily.Time[i],
			fmt.Sprintf("%.1f%s", forecastData.Daily.Temperature2MMax[i], forecastData.DailyUnits.Temperature2MMax),
			fmt.Sprintf("%.1f%s", forecastData.Daily.Temperature2MMin[i], forecastData.DailyUnits.Temperature2MMin),
			fmt.Sprintf("%.1f%s", forecastData.Daily.ApparentTemperatureMax[i], forecastData.DailyUnits.ApparentTemperatureMax),
			fmt.Sprintf("%.1f%s", forecastData.Daily.ApparentTemperatureMin[i], forecastData.DailyUnits.ApparentTemperatureMin),
			forecastData.Daily.Sunrise[i][11:16],
			forecastData.Daily.Sunset[i][11:16],
			fmt.Sprintf("%.1f %s", forecastData.Daily.PrecipitationSum[i], forecastData.DailyUnits.PrecipitationSum),
		})

		analyzeExtremeWeather(forecastData.Daily.Time[i], forecastData.Daily.Temperature2MMax[i], forecastData.Daily.Temperature2MMin[i], forecastData.Daily.PrecipitationSum[i])
	}

	table.Render()
	return plotData, nil
}

func getHistoricalWeatherAPI(lat, lng string, dates []string, cityName string) (PlotData, error) {
	plotData := PlotData{}

	if len(dates) == 0 {
		return plotData, fmt.Errorf("nie podano dat dla historii pogody")
	}

	fmt.Printf("\n--- Historia pogody dla %s---\n", cityName)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Data", "Temp Max", "Temp Min", "Odcz. Max", "Odcz. Min", "Wschód", "Zachód"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, dateStr := range dates {
		t, parseErr := time.Parse("2006-01-02", dateStr)
		if parseErr != nil {
			fmt.Printf("Ostrzeżenie: Nieprawidłowy format daty '%s'. Oczekiwano YYYY-MM-DD. Pomijam.\n", dateStr)
			continue
		}

		if t.After(time.Now()) {
			fmt.Printf("Ostrzeżenie: Data '%s' jest w przyszłości. Pominę pobieranie danych historycznych.\n", dateStr)
			continue
		}

		url := fmt.Sprintf("https://historical-forecast-api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&start_date=%s&end_date=%s&daily=temperature_2m_max,temperature_2m_min,apparent_temperature_max,apparent_temperature_min,sunrise,sunset&timezone=Europe%%2FBerlin", lat, lng, dateStr, dateStr)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Błąd podczas wykonywania zapytania HTTP dla daty %s: %v\n", dateStr, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Nieoczekiwany status odpowiedzi dla daty %s: %s\n", dateStr, resp.Status)
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Odpowiedź API: %s\n", string(body))
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Błąd podczas odczytu odpowiedzi dla daty %s: %v\n", dateStr, err)
			continue
		}

		var historicalData DailyForecastResponse
		if err = json.Unmarshal(body, &historicalData); err != nil {
			fmt.Printf("Błąd podczas parsowania JSON dla daty %s: %v\n", dateStr, err)
			continue
		}

		if len(historicalData.Daily.Time) > 0 {
			plotData.Dates = append(plotData.Dates, t)
			plotData.MaxTemps = append(plotData.MaxTemps, historicalData.Daily.Temperature2MMax[0])
			plotData.MinTemps = append(plotData.MinTemps, historicalData.Daily.Temperature2MMin[0])
			plotData.ApparentMaxTemps = append(plotData.ApparentMaxTemps, historicalData.Daily.ApparentTemperatureMax[0])
			plotData.ApparentMinTemps = append(plotData.ApparentMinTemps, historicalData.Daily.ApparentTemperatureMin[0])

			table.Append([]string{
				historicalData.Daily.Time[0],
				fmt.Sprintf("%.1f%s", historicalData.Daily.Temperature2MMax[0], historicalData.DailyUnits.Temperature2MMax),
				fmt.Sprintf("%.1f%s", historicalData.Daily.Temperature2MMin[0], historicalData.DailyUnits.Temperature2MMin),
				fmt.Sprintf("%.1f%s", historicalData.Daily.ApparentTemperatureMax[0], historicalData.DailyUnits.ApparentTemperatureMax),
				fmt.Sprintf("%.1f%s", historicalData.Daily.ApparentTemperatureMin[0], historicalData.DailyUnits.ApparentTemperatureMin),
				historicalData.Daily.Sunrise[0][11:16],
				historicalData.Daily.Sunset[0][11:16],
			})
		} else {
			fmt.Printf("Brak danych historycznych prognoz dla daty: %s\n", dateStr)
		}
	}
	table.Render()
	return plotData, nil
}
