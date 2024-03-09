package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buffer := make([]byte, 1024)
	charRead, _ := conn.Read(buffer)

	fields := strings.Fields(string(buffer[:charRead]))
	fmt.Println(fields[1])

	if fields[1] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Contains(fields[1], "/echo") {
		path := fields[1]
		firstIndex := strings.Index(path, "/")
		secIndex := strings.Index(path[firstIndex+1:], "/")

		urlPath := path[secIndex+2:]

		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(urlPath)) + "\r\n\r\n" + urlPath

		fmt.Println(response)
		conn.Write([]byte(response))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}
