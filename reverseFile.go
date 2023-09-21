package main

import (
	"fmt"
	"io"
	"syscall"
)

const BUFFER_SIZE = 4096

/*
reverseFile takes a filename as input and creates a new file that contains the
reversed content of the input file. It operates in a buffered fashion, reading
and reversing chunks of the file at a time.

Parameters:
inputFilename: The name of the file to be reversed.

Notes:
- It uses low-level system calls for file operations.
- The reversed file will have a "reverse_" prefix.
*/
func reverseFile(inputFilename string) {
	// Prefix for the reversed file.
	prefix := "reverse_"
	outputFilename := getOutputFilename(inputFilename, prefix)

	// Flags and mode for opening the input file.
	inputFile_flag := syscall.O_RDONLY
	inputFile_mode := uint32(0644)
	inputFileDescriptor, err := syscall.Open(inputFilename, inputFile_flag, inputFile_mode)
	if err != nil {
		fmt.Printf("Error opening the input file, %s\n", err)
		return
	}
	defer syscall.Close(inputFileDescriptor)

	// Flags and mode for creating the output file.
	outputFile_flag := syscall.O_WRONLY | syscall.O_CREAT
	outputFile_mode := uint32(0644)
	outputFileDescriptor, _ := syscall.Open(outputFilename, outputFile_flag, outputFile_mode)
	if err != nil {
		fmt.Printf("Error opening the output file, %s\n`", err)
		return
	}
	defer syscall.Close(outputFileDescriptor)

	// Seek the last buffer-size chunk of the input file.
	position, _ := syscall.Seek(inputFileDescriptor, BUFFER_SIZE*-1, io.SeekEnd)
	if position == -1 {
		position, _ = syscall.Seek(inputFileDescriptor, 0, io.SeekStart)
	}

	buffer := make([]byte, BUFFER_SIZE)
	readCount, _ := syscall.Read(inputFileDescriptor, buffer)
	var oldPosition int64 = 0

	// While there's still data to read, reverse and write it to the output.
	for {
		_, _ = syscall.Write(outputFileDescriptor, reverseBuffer(buffer[:readCount], readCount))

		// Move the file cursor back twice the amount read (to get to the previous chunk).
		oldPosition = position
		position, _ = syscall.Seek(inputFileDescriptor, int64(-2*readCount), io.SeekCurrent)

		if position == -1 {
			position, _ = syscall.Seek(inputFileDescriptor, 0, io.SeekStart)
			readCount, _ = syscall.Read(inputFileDescriptor, buffer[:oldPosition])
			_, _ = syscall.Write(outputFileDescriptor, reverseBuffer(buffer[:readCount], readCount))
			break
		} else {
			readCount, _ = syscall.Read(inputFileDescriptor, buffer)
		}
	}
}

/*
getOutputFilename constructs the name of the reversed file based on the input
filename and the prefix. If the input filename has a directory path, it preserves
the path in the reversed filename.

Parameters:
inputFilename: The name of the file to be reversed.
prefix: Prefix for the reversed filename.
*/
func getOutputFilename(inputFilename, prefix string) string {
	idx := len(inputFilename) - 1
	var filename string = ""

	// Extract the base filename from the full path.
	for idx >= 0 && inputFilename[idx] != '/' {
		filename = string(inputFilename[idx]) + filename
		idx--
	}

	ans := prefix + filename
	if idx < 0 {
		return ans
	} else {
		return inputFilename[:idx+1] + ans
	}
}

/*
reverseBuffer reverses a buffer in-place. It takes a slice of bytes and its length
as input, and returns the reversed slice.

Parameters:
buffer: Slice of bytes to be reversed.
length: Length of the buffer.
*/
func reverseBuffer(buffer []byte, length int) []byte {
	var temp byte
	for i := 0; i < length/2; i++ {
		temp = buffer[i]
		buffer[i] = buffer[length-i-1]
		buffer[length-i-1] = temp
	}
	return buffer
}

/*
printUsage prints a usage message when the number of command-line arguments
is incorrect.

Parameters:
num_arg: Number of command-line arguments.
*/
func printUsage(num_arg int) {
	if num_arg < 2 {
		fmt.Println("No filename found")
	} else {
		fmt.Println("Too many args")
	}
}
