package main

import (
	"fmt"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func createTemperaturePlot(data PlotData, filename string, title string) error {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = "Data"
	p.Y.Label.Text = "Temperatura (°C)"
	p.X.Tick.Marker = plot.ConstantTicks(generateDateTicks(data.Dates))

	maxTempPoints := make(plotter.XYs, len(data.Dates))
	minTempPoints := make(plotter.XYs, len(data.Dates))
	apparentMaxTempPoints := make(plotter.XYs, len(data.Dates))
	apparentMinTempPoints := make(plotter.XYs, len(data.Dates))

	for i := range data.Dates {
		maxTempPoints[i] = plotter.XY{X: float64(data.Dates[i].Unix()), Y: data.MaxTemps[i]}
		minTempPoints[i] = plotter.XY{X: float64(data.Dates[i].Unix()), Y: data.MinTemps[i]}
		apparentMaxTempPoints[i] = plotter.XY{X: float64(data.Dates[i].Unix()), Y: data.ApparentMaxTemps[i]}
		apparentMinTempPoints[i] = plotter.XY{X: float64(data.Dates[i].Unix()), Y: data.ApparentMinTemps[i]}
	}

	err := plotutil.AddLinePoints(p,
		"Temperatura Max", maxTempPoints,
		"Temperatura Min", minTempPoints,
		"Odczuwalna Max", apparentMaxTempPoints,
		"Odczuwalna Min", apparentMinTempPoints,
	)
	if err != nil {
		return fmt.Errorf("nie udało się dodać punktów do wykresu: %w", err)
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		return fmt.Errorf("nie udało się zapisać wykresu: %w", err)
	}
	fmt.Printf("Wykres zapisano do pliku: %s\n", filename)
	return nil
}

func generateDateTicks(dates []time.Time) []plot.Tick {
	ticks := make([]plot.Tick, len(dates))
	for i, t := range dates {
		ticks[i] = plot.Tick{Value: float64(t.Unix()), Label: t.Format("02.01")} // Format as DD.MM
	}
	return ticks
}
