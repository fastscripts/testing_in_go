package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// welcome message
	intro()

	// channel wate vor quit
	doneChan := make(chan bool)

	// start a goroutine to read user input
	go readUserInput(os.Stdin, doneChan)

	// wait for the goroutine to finish
	<-doneChan
	close(doneChan)

	fmt.Println("Goodbye!")

}

func readUserInput(in io.Reader, doneChan chan bool) {
	//func readUserInput(doneChan chan bool) {
	scanner := bufio.NewScanner(in)
	//scanner := bufio.NewScanner(os.Stdin)

	for {
		res, done := checkNumbers(scanner)

		if done {
			doneChan <- true
			return
		}

		fmt.Println(res)
		prompt()
	}
}

func checkNumbers(scanner *bufio.Scanner) (string, bool) {
	// read user input
	scanner.Scan()

	// check if user wants to quit
	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}

	// try to convert input to int

	numToCheck, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "Please enter a number!", false
	}

	_, msg := isPrime(numToCheck)

	return msg, false

}

func intro() {
	fmt.Println("Welcome to the prime number checker!")
	fmt.Println("This program will check if a number is prime or not.")
	fmt.Println("Enter \"q\" to quit")

	prompt()
}

func prompt() {
	fmt.Print("->:")
}

func isPrime(n int) (bool, string) {
	// 0 and 1 are not prime by definition
	if n == 1 || n == 0 {
		return false, fmt.Sprintf("%d is not prime by definition", n)
	}

	// negative numbers are not prime
	if n < 0 {
		return false, fmt.Sprintf("%d is negative number is not prime, by definition", n)
	}

	// use the modulus operator to to see if we have a prime number
	// nur um irgendwas zum testen zu haben

	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false, fmt.Sprintf("%d is not prime because it is divisible by %d", n, i)
		}
	}
	return true, fmt.Sprintf("%d is a prime number", n)
}
