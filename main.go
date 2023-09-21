package main

import (
	"os"
)

func main() {
	// Retrieve command-line arguments
	args := os.Args
	// Get the number of arguments provided
	size := len(args)

	// Check if the number of arguments is not exactly 2
	// (The first argument is the program name itself, and the second is the expected file name)
	if size != 2 {
		// If not, print the usage message
		printUsage(size)
	} else {
		// If correct number of arguments, proceed to reverse the file
		reverseFile(args[1])
	}
}
