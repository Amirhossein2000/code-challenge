package main

import "fmt"

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	length := len(numbers)
	powTwoChan := make(chan int, length)
	sumChan := make(chan int, length)
	outputChan := make(chan int, 1)

	go func() {
		for _, num := range numbers {
			powTwoChan <- num
		}
	}()
	go TwoPow(length, powTwoChan, sumChan)
	go sum(length, sumChan, outputChan)

	fmt.Println(<-outputChan)
}

func TwoPow(inputLen int, inputChan <-chan int, outputChan chan<- int) {
	i := 0
	for num := range inputChan {
		i++
		outputChan <- num * num
		if i >= inputLen {
			return
		}
	}
}

func sum(inputLen int, inputChan <-chan int, outputChan chan<- int) {
	i := 0
	output := 0

	for num := range inputChan {
		i++
		output += num
		if i >= inputLen {
			outputChan <- output
			return
		}
	}
}
