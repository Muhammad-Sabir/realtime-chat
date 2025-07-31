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
	fmt.Println("Write `!q` to disconnect.")

	go func() {
		responseReader := bufio.NewReader(conn)
		for {
			receivedData, err := responseReader.ReadString('\n')
			if err != nil {
				log.Println("Error reading from server:", err)
				os.Exit(1) // Terminate on error
			}
			receivedData = strings.TrimSpace(receivedData)
			fmt.Println("Received:", receivedData)
			fmt.Print("Send: ")
		}
	}()

	stdinReader := bufio.NewReader(os.Stdin)

	for {
		sentence, _ := stdinReader.ReadString('\n')
		sentence = strings.TrimSpace(sentence)

		if sentence == "!q" {
			fmt.Println("Disconnecting...")
			break
		}

		_, err = conn.Write([]byte(sentence + "\n"))
		if err != nil {
			log.Println("Error writing to connection:", err)
			return
		}
	}
}
