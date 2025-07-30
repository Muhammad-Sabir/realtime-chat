package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server listening on:", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Accepted connection from: %v\n", conn.RemoteAddr())
	_, err := conn.Write([]byte("'^]' is the escape character.\n"))
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}

	buffer := make([]byte, 1024)

	for {
		length, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			} else {
				log.Println("Error reading from connection:", err)
			}

			return
		}

		receivedData := strings.TrimSpace(string(buffer[:length]))
		fmt.Printf("Received %d bytes: %s\n", length, receivedData)

		if receivedData == "^]" {
			fmt.Fprintln(conn, "Goodbye!")
			return
		}

		_, err = conn.Write([]byte(receivedData + "\n"))
		if err != nil {
			log.Println("Error writing to connection:", err)
			return
		}
	}
}
