package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err:= os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", err)
	}
	defer f.Close()

	data := make([]byte, 8)
	var line string	
	for {
		n, err := f.Read(data)
		data = data[:n]

		if err != nil && err != io.EOF {
			fmt.Printf("Error reading from file: %v\n", err)
			os.Exit(1)
		}

		if i := bytes.IndexByte(data, '\n'); i != -1 {
			line += string(data[:i])
			fmt.Printf("read: %s\n", line)
			line = string(data[i+1:])
		} else {
			line += string(data)
		}

		if err == io.EOF {
			break
		}
	}
	if len(line) > 0 {
		fmt.Printf("read: %s\n", line)
	}
}	