package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handle_command(conn)
	}
}

func handle_command(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)
	for {
		err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return
		}

		_, err = c.Read(buf)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Println("Deadline exceeded")
			} else {
				fmt.Println("Error reading from connection: ", err.Error())
			}
			return
		}

		c.SetWriteDeadline(time.Now().Add(2 * time.Second))
		_, err = c.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			return
		}
	}
}
