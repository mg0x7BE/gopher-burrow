package main

import "fmt"

func total(prices []float64) float64 {
	sum := 0.0
	fmt.Println("debug:", len(prices))
	for _, p := range prices {
		sum += p
	}
	return sum
}

func main() {
	prices := []float64{9.99, 19.99, 4.50}
	fmt.Println("items:", len(prices))
	fmt.Println("total:", total(prices))
	fmt.Println("done")
}
