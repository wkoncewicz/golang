package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Load data and configuration at the start
	if err := loadCityData("worldcities.csv"); err != nil {
		fmt.Printf("Błąd ładowania danych miast: %v\n", err)
		os.Exit(1)
	}

	if err := loadConfig(); err != nil {
		fmt.Printf("Błąd ładowania konfiguracji: %v\n", err)
		// Don't exit, but warn if config fails (extreme weather analysis won't work)
	}

	// --- Command-line flag setup ---
	presentCmd := flag.NewFlagSet("aktualna", flag.ExitOnError)
	cityPtrPresent := presentCmd.String("city", "Warszawa", "City to check weather (default Warsaw)")

	futureCmd := flag.NewFlagSet("prognoza", flag.ExitOnError)
	cityPtrFuture := futureCmd.String("city", "Warszawa", "City to check weather (default Warsaw)")
	daysPtr := futureCmd.Int("days", 5, "Days in future to check weather (default 5)")

	pastCmd := flag.NewFlagSet("historia", flag.ExitOnError)
	cityPtrPast := pastCmd.String("city", "Warszawa", "City to check weather (default Warsaw)")
	var dates stringArrayFlag
	pastCmd.Var(&dates, "dates", "Comma-separated list of dates (e.g., 2023-01-01,2023-01-02)")

	// --- Command-line argument parsing and execution ---
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "aktualna":
		presentCmd.Parse(os.Args[2:])
		lat, lng, err := getCityCoordinates(*cityPtrPresent)
		if err != nil {
			fmt.Printf("Błąd: %v\n", err)
			return
		}
		fmt.Printf("Współrzędne dla %s (Polska): Lat %s, Lng %s\n", *cityPtrPresent, lat, lng)
		fmt.Printf("Aktualna pogoda dla miasta: %s\n", *cityPtrPresent)
		if err := getWeatherAPI(lat, lng); err != nil {
			fmt.Printf("Błąd pobierania danych pogodowych: %v\n", err)
		}

	case "prognoza":
		futureCmd.Parse(os.Args[2:])
		lat, lng, err := getCityCoordinates(*cityPtrFuture)
		if err != nil {
			fmt.Printf("Błąd: %v\n", err)
			return
		}
		fmt.Printf("Współrzędne dla %s (Polska): Lat %s, Lng %s\n", *cityPtrFuture, lat, lng)
		plotData, apiErr := getForecastWeatherAPI(lat, lng, *daysPtr, *cityPtrFuture)
		if apiErr != nil {
			fmt.Printf("Błąd pobierania prognozy pogody: %v\n", apiErr)
			return
		}
		if len(plotData.Dates) > 0 {
			plotFilename := fmt.Sprintf("prognoza_%s.png", strings.ToLower(strings.ReplaceAll(*cityPtrFuture, " ", "_")))
			if err := createTemperaturePlot(plotData, plotFilename, fmt.Sprintf("Prognoza temperatury dla %s", *cityPtrFuture)); err != nil {
				fmt.Printf("Błąd generowania wykresu: %v\n", err)
			}
		}

	case "historia":
		pastCmd.Parse(os.Args[2:])
		lat, lng, err := getCityCoordinates(*cityPtrPast)
		if err != nil {
			fmt.Printf("Błąd: %v\n", err)
			return
		}
		fmt.Printf("Współrzędne dla %s (Polska): Lat %s, Lng %s\n", *cityPtrPast, lat, lng)
		plotData, apiErr := getHistoricalWeatherAPI(lat, lng, dates, *cityPtrPast)
		if apiErr != nil {
			fmt.Printf("Błąd pobierania historii pogody: %v\n", apiErr)
			return
		}
		if len(plotData.Dates) > 0 {
			plotFilename := fmt.Sprintf("historia_%s.png", strings.ToLower(strings.ReplaceAll(*cityPtrPast, " ", "_")))
			if err := createTemperaturePlot(plotData, plotFilename, fmt.Sprintf("Historia temperatury dla %s", *cityPtrPast)); err != nil {
				fmt.Printf("Błąd generowania wykresu: %v\n", err)
			}
		}

	default:
		fmt.Printf("Nieznana komenda: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Użycie: program <komenda> [flagi]")
	fmt.Println("Komendy:")
	fmt.Println("   aktualna    - Aktualna pogoda")
	fmt.Println("   prognoza    - Prognoza pogody")
	fmt.Println("   historia    - Historia pogody")
}
