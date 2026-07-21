package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
	result, msg := isPrime(0)
	if result {
		t.Errorf("with %d as test paramter, got true, but expected false", 0)
	}

	if msg != "0 is not prime by definition" {
		t.Error("wrong message returned:", msg)
	}

	result, msg = isPrime(7)

	if !result {
		t.Errorf("with %d as test paramter, got false, but expected true", 7)
	}

	if msg != "7 is a prime number" {
		t.Error("wrong message returned:", msg)
	}
}

func Test_isPrimeTable(t *testing.T) {
	primeTests := []struct {
		name        string
		testNum     int
		expected    bool
		expectedMsg string
	}{
		{"Test with 0", 0, false, "0 is not prime by definition"},
		{"Test with 1", 1, false, "1 is not prime by definition"},
		{"Test with -5", -5, false, "-5 is negative number is not prime, by definition"},
		{"Test with 3", 3, true, "3 is a prime number"},
		{"Test with 4", 4, false, "4 is not prime because it is divisible by 2"},
	}

	for _, tt := range primeTests {
		t.Run(tt.name, func(t *testing.T) {
			result, msg := isPrime(tt.testNum)
			if result != tt.expected {
				t.Errorf("with %d as test parameter, got %v, but expected %v", tt.testNum, result, tt.expected)
			}
			if msg != tt.expectedMsg {
				t.Errorf("with %d as test parameter, got message '%s', but expected '%s'", tt.testNum, msg, tt.expectedMsg)
			}
		})
	}
}

func Test_prompt(t *testing.T) {
	// save os.Stdout
	oldOut := os.Stdout

	// create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// call the prompt function
	prompt()

	// read the output from the pipe
	w.Close()

	// read the output of prompt()
	output, _ := io.ReadAll(r)

	// check if the output is as expected

	expectedOutput := "->:"
	if string(output) != expectedOutput {
		t.Errorf("Expected output '%s', but got '%s'", expectedOutput, output)
	}

	// restore os.Stdout
	os.Stdout = oldOut
}

func Test_intro(t *testing.T) {
	// save os.Stdout
	oldOut := os.Stdout

	// create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// call the prompt function
	intro()

	// read the output from the pipe
	w.Close()

	// read the output of prompt()
	output, _ := io.ReadAll(r)

	// check if the output is as expected

	if !strings.Contains(string(output), "Welcome to the prime number checker") {
		t.Errorf("Expected output \"Welcome to the prime number checker\", not in %s", string(output))
	}

	// restore os.Stdout
	os.Stdout = oldOut
}

func Test_checkNumbers(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty input", input: "", expected: "Please enter a number!"},
		{name: "3", input: "3", expected: "3 is a prime number"},
		{name: "quit", input: "q", expected: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := strings.NewReader(test.input)
			reader := bufio.NewScanner(input)
			result, _ := checkNumbers(reader)
			if !strings.EqualFold(result, test.expected) {
				t.Errorf("incorrect response, expected: %s got %s", test.expected, result)
			}
		})
	}

}

func Test_readUserInput(t *testing.T) {
	doneChan := make(chan bool)
	// create a reference to a bytes.Buffer
	var stdin bytes.Buffer
	stdin.Write([]byte("3\nq\n"))

	go readUserInput(&stdin, doneChan)
	<-doneChan
	close(doneChan)
}
