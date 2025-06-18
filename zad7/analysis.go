package main

import "fmt"

func analyzeExtremeWeather(date string, maxTemp, minTemp, precipitation float64) {
	if extremeThresholds.HighTemperatureCelsius > 0 && maxTemp > extremeThresholds.HighTemperatureCelsius {
		fmt.Printf("!!! OSTRZEŻENIE !!! Dnia %s przewidywana jest ekstremalnie wysoka temperatura: %.1f°C (powyżej %.1f°C)\n", date, maxTemp, extremeThresholds.HighTemperatureCelsius)
	}
	if extremeThresholds.LowTemperatureCelsius != 0 && minTemp < extremeThresholds.LowTemperatureCelsius {
		fmt.Printf("!!! OSTRZEŻENIE !!! Dnia %s przewidywana jest ekstremalnie niska temperatura: %.1f°C (poniżej %.1f°C)\n", date, minTemp, extremeThresholds.LowTemperatureCelsius)
	}
	if extremeThresholds.MaxRainMmPerDay > 0 && precipitation > extremeThresholds.MaxRainMmPerDay {
		fmt.Printf("!!! OSTRZEŻENIE !!! Dnia %s przewidywane są intensywne opady: %.1f mm (powyżej %.1f mm)\n", date, precipitation, extremeThresholds.MaxRainMmPerDay)
	}
}
