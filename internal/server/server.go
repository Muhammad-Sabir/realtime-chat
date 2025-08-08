package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/Muhammad-Sabir/realtime-chat/internal/models"
)

var store *models.UserStore

func Start(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server listening on:", listener.Addr().String())

	store = models.NewUserStore()

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

	user := readClientUser(clientInputReader)

	err := store.AddUser(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer store.RemoveUser(user)

	disconnect := make(chan string)
	clientMsg := make(chan models.Message)
	go readClientMessage(clientMsg, clientInputReader, disconnect)

	for {
		select {
		case message := <-clientMsg:
			fmt.Printf("Received: %s\n", message.Content)
			go writeToClient(conn, message)
		case disconnectClient := <-disconnect:
			if disconnectClient == "exit()" {
				log.Println("Disconnecting the user: ", user.Email)
				return
			}
		}

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

func readClientMessage(clientMsg chan<- models.Message, reader *bufio.Reader, disconnect chan<- string) {
	for {
		var msg models.Message

		receivedData, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			disconnect <- "exit()"
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

func readClientUser(reader *bufio.Reader) *models.User {
	var user models.User

	receivedData, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading from connection:", err)
		os.Exit(1)
	}

	err = json.Unmarshal([]byte(receivedData), &user)
	if err != nil {
		log.Println("Error unmarshalling:", err)
		os.Exit(1)
	}

	return &user
}
