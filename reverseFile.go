package main

import (
	"io"
	"syscall"
)

const BUFFER_SIZE = 3

func reverseFile(inputFilename string) {
	prefix := "reverse_"
	outputFilename := getOutputFilename(inputFilename, prefix)

	inputFile_flag := syscall.O_RDONLY
	inputFile_mode := uint32(0644)
	inputFileDescriptor, _ := syscall.Open(inputFilename, inputFile_flag, inputFile_mode)

	defer syscall.Close(inputFileDescriptor)

	outputFile_flag := syscall.O_WRONLY | syscall.O_CREAT
	outputFile_mode := uint32(0644)
	outputFileDescriptor, _ := syscall.Open(outputFilename, outputFile_flag, outputFile_mode)

	defer syscall.Close(outputFileDescriptor)

	position, _ := syscall.Seek(inputFileDescriptor, BUFFER_SIZE*-1, io.SeekEnd)
	if position == -1 {
		position, _ = syscall.Seek(inputFileDescriptor, 0, io.SeekStart)
	}

	buffer := make([]byte, BUFFER_SIZE)
	readCount, _ := syscall.Read(inputFileDescriptor, buffer)
	var oldPosition int64 = 0

	for readCount > 0 {
		_, _ = syscall.Write(outputFileDescriptor, reverseBuffer(buffer[:readCount], readCount))

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

func getOutputFilename(inputFilename, prefix string) string {
	return "reverse_" + inputFilename
}

// func printBuffer(buffer, length int) {

// }

func reverseBuffer(buffer []byte, length int) []byte {
	var temp byte
	for i := 0; i < length/2; i++ {
		temp = buffer[i]
		buffer[i] = buffer[length-i-1]
		buffer[length-i-1] = temp
	}
	return buffer
}

// func printUsage(int num_argc) {

// }
