package main

import "fmt"

func main() {
	fmt.Println(Div(1, 2))
}

func Div(a, b float64) float64 {
	return a/b + 1
}
