package main

import (
	"strings"
	"time"
)

// --- Data Structs ---

type Cities struct {
	city       string
	city_ascii string
	lat        string
	lng        string
	country    string
}

type CurrentWeatherResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentUnits         struct {
		Time                string `json:"time"`
		Interval            string `json:"interval"`
		Temperature2M       string `json:"temperature_2m"`
		RelativeHumidity2M  string `json:"relative_humidity_2m"`
		ApparentTemperature string `json:"apparent_temperature"`
		IsDay               string `json:"is_day"`
		Rain                string `json:"rain"`
		Precipitation       string `json:"precipitation"`
	} `json:"current_units"`
	Current struct {
		Time                string  `json:"time"`
		Interval            int     `json:"interval"`
		Temperature2M       float64 `json:"temperature_2m"`
		RelativeHumidity2M  int     `json:"relative_humidity_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		IsDay               int     `json:"is_day"`
		Rain                float64 `json:"rain"`
		Precipitation       float64 `json:"precipitation"`
	} `json:"current"`
}

type DailyForecastResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	DailyUnits           struct {
		Time                   string `json:"time"`
		Temperature2MMax       string `json:"temperature_2m_max"`
		Temperature2MMin       string `json:"temperature_2m_min"`
		ApparentTemperatureMax string `json:"apparent_temperature_max"`
		ApparentTemperatureMin string `json:"apparent_temperature_min"`
		Sunrise                string `json:"sunrise"`
		Sunset                 string `json:"sunset"`
		PrecipitationSum       string `json:"precipitation_sum"`
	} `json:"daily_units"`
	Daily struct {
		Time                   []string  `json:"time"`
		Temperature2MMax       []float64 `json:"temperature_2m_max"`
		Temperature2MMin       []float64 `json:"temperature_2m_min"`
		ApparentTemperatureMax []float64 `json:"apparent_temperature_max"`
		ApparentTemperatureMin []float64 `json:"apparent_temperature_min"`
		Sunrise                []string  `json:"sunrise"`
		Sunset                 []string  `json:"sunset"`
		PrecipitationSum       []float64 `json:"precipitation_sum"`
	} `json:"daily"`
}

type PlotData struct {
	Dates            []time.Time
	MaxTemps         []float64
	MinTemps         []float64
	ApparentMaxTemps []float64
	ApparentMinTemps []float64
}

type ExtremeWeatherThresholds struct {
	HighTemperatureCelsius float64 `mapstructure:"high_temperature_celsius"`
	LowTemperatureCelsius  float64 `mapstructure:"low_temperature_celsius"`
	MaxRainMmPerDay        float64 `mapstructure:"max_rain_mm_per_day"`
}

// --- Flag handling ---
type stringArrayFlag []string

func (i *stringArrayFlag) String() string {
	return strings.Join(*i, ",")
}

func (i *stringArrayFlag) Set(value string) error {
	*i = strings.Split(value, ",")
	return nil
}
