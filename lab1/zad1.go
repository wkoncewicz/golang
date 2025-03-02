package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// GeneratePESEL: geneuje numer PESEL
// Parametry:
// - birthDate: time.Time: reprezentacja daty urodzenia
// - płeć: znak "M" lub "K"
// Wyjscie:
// Tablica z cyframi numeru PESEL
func GenerujPESEL(birthDate time.Time, gender string) [11]int {

	// tablica zawierajaca kolejne cyfry numeru PESEL
	var cyfryPESEL [11]int

	// konwersja daty na dane skladowe
	year := birthDate.Year()
	month := PoliczMiesiac(int(birthDate.Month()), year)
	day := birthDate.Day()
	fmt.Println(month)

	// losowy numer
	randomSerial := rand.Intn(900) + 100 // 3 cyfrowy losowy numer z zakresu 100-999

	var genderNum int
	if gender == "M" {
		genderNum = ((rand.Intn(5) + 1) * 2) - 1
	} else {
		genderNum = (rand.Intn(5) + 1) * 2
	}
	yearArr := LiczbaDoListy(year)
	monthArr := LiczbaDoListy(month)
	if len(monthArr) == 1 {
		monthArr = []int{0, monthArr[0]}
	}
	dayArr := LiczbaDoListy(day)
	randomSerialArr := LiczbaDoListy(randomSerial)
	cyfraKontrolna := ObliczCyfre([10]int{yearArr[2], yearArr[3], monthArr[0], monthArr[1], dayArr[0], dayArr[1], randomSerialArr[0], randomSerialArr[1], randomSerialArr[2], genderNum})

	cyfryPESEL = [11]int{yearArr[2], yearArr[3], monthArr[0], monthArr[1], dayArr[0], dayArr[1], randomSerialArr[0], randomSerialArr[1], randomSerialArr[2], genderNum, cyfraKontrolna}

	return cyfryPESEL
}

// WeryfikujPESEL: weryfikuje poprawność numeru PESEL
// Parametry:
// - cyfryPESEL: Tablica z cyframi numeru PESEL
// Wyjscie:
//zmienna bool

func PoliczMiesiac(month int, year int) int {
	if year < 1900 {
		return month + 80
	} else if year < 2000 {
		return month
	} else if year < 2100 {
		return month + 20
	} else if year < 2200 {
		return month + 40
	} else if year < 2300 {
		return month + 60
	} else {
		return month
	}
}

func LiczbaDoListy(n int) []int {
	str := strconv.Itoa(n)
	arr := make([]int, len(str))

	for i, char := range str {
		arr[i] = int(char - '0')
	}
	return arr
}

func ObliczCyfre(cyfry [10]int) int {
	num1 := Liczba(int(cyfry[0]), 1)
	num2 := Liczba(int(cyfry[1]), 3)
	num3 := Liczba(int(cyfry[2]), 7)
	num4 := Liczba(int(cyfry[3]), 9)
	num5 := Liczba(int(cyfry[4]), 1)
	num6 := Liczba(int(cyfry[5]), 3)
	num7 := Liczba(int(cyfry[6]), 7)
	num8 := Liczba(int(cyfry[7]), 9)
	num9 := Liczba(int(cyfry[8]), 1)
	num10 := Liczba(int(cyfry[9]), 3)
	sum := num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9 + num10
	if sum >= 10 {
		return 10 - LiczbaDoListy(sum)[0]
	} else {
		return 10 - sum
	}
}

func Liczba(liczba int, mnoznik int) int {
	num := liczba * mnoznik
	if num >= 10 {
		num = num % 10
	}
	return num
}

func WeryfikujPESEL(cyfryPESEL [11]int) bool {

	var czyPESEL bool
	if cyfryPESEL[9]%2 == 0 {
		czyPESEL = false
	} else {
		czyPESEL = true
	}

	return czyPESEL
}

// Przykład użycia
func main() {
	//
	birthDate := time.Date(2005, 3, 10, 0, 0, 0, 0, time.UTC)
	pesel := GenerujPESEL(birthDate, "M")

	fmt.Println("Wygenerowany PESEL:", pesel)

	fmt.Println("Czy numer PESEL jest poprawny:", WeryfikujPESEL(pesel))
}
