package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", err)
	}
	defer f.Close()

	buf := make([]byte, 8) 
	var line string

	for {
		n, err := f.Read(buf)

		readData := buf[:n]

		if i := bytes.IndexByte(readData, '\n'); i != -1 {
			line += string(readData[:i])
			fmt.Printf("read: %s\n", line)
			line = string(readData[i+1:])
		} else {
			line += string(readData)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading from file: %v\n", err)
			os.Exit(1)
		}
	}

	if len(line) > 0 {
		fmt.Printf("read: %s\n", line)
	}
}