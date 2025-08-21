package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		buf := make([]byte, 8) 
		var line string

		for {
			n, err := f.Read(buf)

			readData := buf[:n]

			if i := bytes.IndexByte(readData, '\n'); i != -1 {
				line += string(readData[:i])
				out <- line
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
			out <- line
		}
	}()

	return out

}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}
		for word := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", word)
		}
	}


}