package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

func main() {

	flag.StringVar(&directory, "directory", "", "directory path")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go server(conn)

	}
}

func server(conn net.Conn) {
	buffer := make([]byte, 1024)
	charRead, _ := conn.Read(buffer)

	fields := strings.Fields(string(buffer[:charRead]))

	fmt.Println(string(buffer[:charRead]))
	fmt.Println(fields)

	if fields[1] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.Contains(fields[1], "/echo") {
		path := fields[1]
		firstIndex := strings.Index(path, "/")
		secIndex := strings.Index(path[firstIndex+1:], "/")

		urlPath := path[secIndex+2:]

		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(urlPath)) + "\r\n\r\n" + urlPath
		conn.Write([]byte(response))
	} else if strings.Contains(fields[1], "/user-agent") {
		header := fields[6]

		fmt.Println(header, len(header))

		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(header)) + "\r\n\r\n" + header
		conn.Write([]byte(response))
	} else if strings.Contains(fields[1], "/files") {

		fmt.Println(fields[1])
		extractedFileName := strings.TrimPrefix(fields[1], "/files/")
		// fmt.Println(fileNameIndex)
		// extractedFileName := string(fields[1][fileNameIndex:])

		fmt.Println(directory, extractedFileName)

		file, err := os.ReadFile(directory + extractedFileName)

		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"))

		} else {
			response := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + fmt.Sprint(len(file)) + "\r\n\r\n" + string(file)
			conn.Write([]byte(response))
		}

	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"))
	}

}
