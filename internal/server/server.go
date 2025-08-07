package server

import (
	"bufio"
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

	clientInput := make(chan string)
	writeToClient(conn, "Enter your name: ")
	go readFromClient(conn, clientInput)

	clientName := <-clientInput
	fmt.Println(clientName + " connected.")

	for {
		receivedData := <-clientInput
		fmt.Printf("Received: %s\n", receivedData)

		if receivedData == "exit()" {
			fmt.Println("Disconnecting...")
			writeToClient(conn, "Goodbye...")
			return
		}

		receivedData = clientName + ": " + receivedData
		go writeToClient(conn, receivedData)
	}
}

func writeToClient(conn net.Conn, data string) {
	_, err := conn.Write([]byte(data + "\n"))
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}

func readFromClient(conn net.Conn, clientInput chan<- string) {
	clientInputReader := bufio.NewReader(conn)

	for {
		receivedData, err := clientInputReader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			} else {
				log.Println("Error reading from connection:", err)
			}

			clientInput <- "exit()"
			return
		}
		clientInput <- strings.TrimSpace(receivedData)
	}
}
