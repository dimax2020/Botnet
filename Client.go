package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	for {
		
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Отправить: ")
		text, _ := reader.ReadString('\n')
		
		// Отправляем в socket
		conn.Write([]byte(text + "\n"))
		
		// Прослушиваем ответ
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Сообщение от сервера: " + message)
	}
}
