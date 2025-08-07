package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func Start(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to:", conn.RemoteAddr())
	fmt.Println("Write `exit()` to disconnect.")

	clientInput := make(chan string)
	serverResponse := make(chan string)

	go readFromServer(conn, serverResponse)
	go takeUserInput(clientInput)

	for {
		select {
		case sentence := <-clientInput:
			if sentence == "exit()" {
				writeToServer(conn, sentence)
				fmt.Println("Disconnecting...")
				return
			}
			go writeToServer(conn, sentence)
		case serverResponse := <-serverResponse:
			fmt.Println(serverResponse)
		}
	}
}

func readFromServer(conn net.Conn, serverResponse chan<- string) {
	responseReader := bufio.NewReader(conn)
	for {
		receivedData, err := responseReader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from server:", err)
			os.Exit(1) // Terminate on error
		}
		serverResponse <- strings.TrimSpace(receivedData)
	}
}

func takeUserInput(clientInput chan<- string) {
	stdinReader := bufio.NewReader(os.Stdin)

	for {
		sentence, _ := stdinReader.ReadString('\n')
		clientInput <- strings.TrimSpace(sentence)
	}
}

func writeToServer(conn net.Conn, clientInput string) {
	_, err := conn.Write([]byte(clientInput + "\n"))
	if err != nil {
		log.Println("Error writing to connection:", err)
		return
	}
}
