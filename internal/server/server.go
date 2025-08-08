package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/Muhammad-Sabir/realtime-chat/internal/models"
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

	clientInputReader := bufio.NewReader(conn)

	clientMsg := make(chan models.Message)

	go readFromClient(clientMsg, clientInputReader)

	for {
		message := <-clientMsg
		fmt.Printf("Received: %s\n", message.Content)

		if message.Content == "exit()" {
			fmt.Println("Disconnecting...")
			return
		}

		go writeToClient(conn, message)
	}
}

func writeToClient(conn net.Conn, msg models.Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}

func readFromClient(clientMsg chan<- models.Message, reader *bufio.Reader) {
	for {
		var msg models.Message

		receivedData, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("readFromClient EOF", err)
			} else {
				log.Println("Error reading from connection:", err)
			}
			return
		}

		err = json.Unmarshal([]byte(receivedData), &msg)
		if err != nil {
			log.Println("Error unmarshalling:", err)
			continue
		}

		clientMsg <- msg
	}
}
