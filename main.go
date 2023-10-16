package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	println("wait for connection")
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	println("user was connected")
	defer conn.Close()
	for {
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		if err != nil {
			if err == io.EOF {
				fmt.Println("user was disconected")
				break
			}
			fmt.Println("can not get property input from user")
			os.Exit(0)
		}
		conn.Write([]byte("+OK\r\n"))
	}

}
