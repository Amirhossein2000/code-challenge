package main

import "fmt"

func main() {
	numbers := []int{2, 2, 5, 6, 5}
	exception := findException(numbers)
	if exception > 0 {
		fmt.Println("element occurring once is", exception)
	} else {
		fmt.Println("there isn't any exception in these numbers", numbers)
	}
}

func findException(numbers []int) int {
	output := numbers[0]
	for i := 1; i < len(numbers); i++ {
		output ^= numbers[i]
	}
	return output
}
