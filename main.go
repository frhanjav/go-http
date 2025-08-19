package main

import (
	"fmt"
	"log"
	"os"
	"io"
)

func main() {
	f, err:= os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", err)
	}
	defer f.Close()

	data := make([]byte, 8)
	
	for {
		n, err := f.Read(data)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("Error reading from file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("read: %s\n", string(data[:n]))

	}
}