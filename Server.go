package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Message from", conn.LocalAddr(), ":", message)
		newMessage := strings.ToUpper(message)
		conn.Write([]byte(newMessage + "\n"))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	fmt.Println("Starting server...")
	var protocol = "tcp"
	var address = "127.0.0.1:8081"

	ln, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Прослушивается: " + address + "\nПротокол: " + protocol)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Новое подключение:", conn.LocalAddr())
		go handleConnection(conn)
	}
}
