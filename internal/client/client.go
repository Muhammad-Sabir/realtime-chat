package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Muhammad-Sabir/realtime-chat/internal/models"
)

var user models.User

func Start(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to:", conn.RemoteAddr())
	fmt.Print("Write `exit()` to disconnect. \n\n")

	stdinReader := bufio.NewReader(os.Stdin)
	responseReader := bufio.NewReader(conn)

	clientInput := make(chan string)
	serverResponse := make(chan models.Message)

	getUserData(stdinReader)
	writeUserToServer(conn)

	go readFromServer(serverResponse, responseReader)
	go takeUserInput(clientInput, stdinReader)

	for {
		select {
		case text := <-clientInput:
			if text == "exit()" {
				fmt.Println("Disconnecting...")
				return
			}

			msg := models.NewMessage(user, text)
			go writeMessageToServer(conn, msg)
		case serverMsg := <-serverResponse:
			fmt.Println(serverMsg.String())
		}
	}
}

func getUserData(reader *bufio.Reader) {
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter your email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	user = models.NewUser(name, email)
}

func readFromServer(serverResponse chan<- models.Message, reader *bufio.Reader) {
	for {
		var msg models.Message

		receivedData, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from server:", err)
			os.Exit(1)
		}

		err = json.Unmarshal([]byte(receivedData), &msg)
		if err != nil {
			log.Println("Error unmarshalling:", err)
			continue
		}

		serverResponse <- msg
	}
}

func takeUserInput(clientInput chan<- string, reader *bufio.Reader) {
	for {
		sentence, _ := reader.ReadString('\n')
		clientInput <- strings.TrimSpace(sentence)
	}
}

func writeMessageToServer(conn net.Conn, msg models.Message) {
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

func writeUserToServer(conn net.Conn) {
	data, err := json.Marshal(user)
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
