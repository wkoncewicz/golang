package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type Nomenclature struct {
	code       string
	de_label   string
	en_label   string
	es_label   string
	fr_label   string
	pt_label   string
	short_code string
}

func readCSVLineByLine(filepath string) []string {
	var accidents_slice []Nomenclature

	// otwarcie pliku
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	//ustawienie separatora
	reader.Comma = ';'

	//pominiecie pierwszej lini (naglowek)
	_, err = reader.Read()

	// czytanie linia po linii
	for i := 0; i < 5; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file
		}

		accidents_slice = append(accidents_slice, Nomenclature{record[0], record[1], record[2], record[3], record[4], record[5], record[6]})
	}

	// sort.Slice(accidents_slice, func(i, j int) bool { return accidents_slice[i].de_label < accidents_slice[j].de_label })
	// var labels []string
	// for _, item := range accidents_slice {
	// 	labels = append(labels, item.de_label)
	// }

	sort.Slice(accidents_slice, func(i, j int) bool { return accidents_slice[i].en_label < accidents_slice[j].en_label })
	var labels []string
	for _, item := range accidents_slice {
		labels = append(labels, item.en_label)
	}

	return labels
}

func main() {
	var new_data []string

	new_data = readCSVLineByLine("./nomenclature-cpv.csv")

	fmt.Println(new_data)

}
