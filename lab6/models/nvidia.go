package models

import "time"

type NVIDIA struct {
	Date      time.Time
	CloseLast float64
	Volume    int64
	Open      float64
	High      float64
	Low       float64
}
