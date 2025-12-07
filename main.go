package main

import (
	"fmt"
	"os"
)

func main() {
	// 1. Ensure we have Input and Output file arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input_file> <output_file>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// 2. Read File
	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// 3. Processing - casting (main logic)
	result := NormalizeText(string(content))

	// 4. Write File
	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
}
