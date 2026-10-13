package main

import "fmt"

// total sums the prices and applies the discount.
func total(prices []float64, discount float64) float64 {
	sum := 0.0
	for _, p := range prices {
		sum += p
	}
	return sum * (1 - discount)
}

func main() {
	prices := []float64{9.99, 19.99, 4.50, 12.00}
	fmt.Println("total:", total(prices, 0.1))
	fmt.Println("done")
}
