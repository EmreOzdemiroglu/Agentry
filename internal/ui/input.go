package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// InputReader handles user input operations
type InputReader struct {
	reader *bufio.Reader
}

// NewInputReader creates a new input reader
func NewInputReader() *InputReader {
	return &InputReader{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ReadInput reads a line of text from stdin after printing a prompt
func (r *InputReader) ReadInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := r.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ReadNonEmptyInput reads input and ensures it's not empty
func (r *InputReader) ReadNonEmptyInput(prompt string, errorMsg string) string {
	for {
		input := r.ReadInput(prompt)
		if strings.TrimSpace(input) != "" {
			return input
		}
		if errorMsg != "" {
			fmt.Println(errorMsg)
		}
	}
}

// ReadChoice reads a choice from a set of valid options (case insensitive)
func (r *InputReader) ReadChoice(prompt string, validChoices []string) string {
	for {
		input := strings.ToUpper(r.ReadInput(prompt))
		for _, choice := range validChoices {
			if input == strings.ToUpper(choice) {
				return input
			}
		}
		fmt.Printf("Invalid choice. Please enter one of: %s\n", strings.Join(validChoices, ", "))
	}
}

// ReadNumber reads a number within a specified range
func (r *InputReader) ReadNumber(prompt string, min, max int) int {
	for {
		input := r.ReadInput(prompt)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < min || choice > max {
			fmt.Printf("Invalid selection. Please enter a number between %d and %d.\n", min, max)
			continue
		}
		return choice
	}
}

// ReadYesNo reads a yes/no answer
func (r *InputReader) ReadYesNo(prompt string) bool {
	choice := r.ReadChoice(prompt+" [Y/N]", []string{"Y", "N"})
	return choice == "Y"
}