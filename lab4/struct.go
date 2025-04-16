package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Struktura modelująca zamówienie
type Order struct {
	ID           int
	CustomerName string
	Items        []string
	TotalAmount  float64
}

// Struktura modelująca prze
type ProcessResult struct {
	OrderID      int
	CustomerName string
	Success      bool
	ProcessTime  time.Duration
	Error        error
}

func makeOrders(ch chan<- Order, orderCount int) {
	namesPool := []string{"Stasiek", "Kuba", "Wiktor", "Staszek", "Stefan", "Małgorzata", "Steve", "Magdalena"}
	itemsPool := []string{"BigMac", "MacChicken", "Frytki", "MacNuggets", "MacRoyale", "WieśMac", "Cheeseburger", "MacDouble"}
	for i := 1; i <= orderCount; i++ {
		nameIndex := rand.Intn(len(namesPool))
		itemsNum := rand.Intn(len(itemsPool)) + 1
		var itemIndex int
		var item string
		var items []string
		for i := 0; i < itemsNum; i++ {
			itemIndex = rand.Intn(len(itemsPool))
			item = itemsPool[itemIndex]
			items = append(items, item)
		}
		name := namesPool[nameIndex]
		totalAmount := rand.Float64()*20 + float64(len(items))
		order := Order{ID: i, CustomerName: name, Items: items, TotalAmount: totalAmount}
		ch <- order
	}
	close(ch)
}

func worker(ordersChan <-chan Order, resultsChan chan<- ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range ordersChan {
		duration := time.Duration(rand.Intn(1000)+500) * time.Millisecond
		time.Sleep(duration)

		var err error
		success := rand.Float32() < 0.85
		for !success {
			err = fmt.Errorf("Zamówienie nieudane!")
			result := ProcessResult{OrderID: order.ID, CustomerName: order.CustomerName, Success: success, ProcessTime: duration, Error: err}
			resultsChan <- result

			duration = time.Duration(rand.Intn(1000)+500) * time.Millisecond
			time.Sleep(duration)

			success = rand.Float32() < 0.85
		}

		result := ProcessResult{OrderID: order.ID, CustomerName: order.CustomerName, Success: success, ProcessTime: duration, Error: err}
		resultsChan <- result
	}
}

func processResults(resultsChan <-chan ProcessResult, done chan<- struct{}) {
	var total, success int
	for result := range resultsChan {
		if !result.Success {
			fmt.Println(result.Error, "Ponawiam próbę dla zamówienia:", result.OrderID)
		} else {
			fmt.Println("Zamówienie udane!")
			fmt.Println("ID:", result.OrderID)
			fmt.Println("Czas działania:", result.ProcessTime)
			success++
		}
		total++
	}
	fmt.Println("Liczba udanych zamówień:", success)
	fmt.Println("Liczba zamówień:", total)
	precent := float64(success) / float64(total) * 100
	fmt.Println("Procent dokładności:", precent)
	done <- struct{}{}
}

func main() {
	workerCount := 3
	orderCount := 15

	ordersChan := make(chan Order, orderCount)
	resultsChan := make(chan ProcessResult, orderCount)
	done := make(chan struct{})

	go makeOrders(ordersChan, orderCount)

	var wg sync.WaitGroup
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(ordersChan, resultsChan, &wg)
	}

	go processResults(resultsChan, done)
	wg.Wait()
	close(resultsChan)
	<-done
}
