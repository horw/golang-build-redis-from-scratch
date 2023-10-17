package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	aof, err := NewAof("hello.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	println("wait for connection")
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	println("user was connected")
	defer conn.Close()
	for {
		resp := NewResp(conn)
		val, err := resp.Read()

		if err != nil {
			if err == io.EOF {
				fmt.Println("user was disconected")
				break
			}
			fmt.Println("can not get property input from user")
			os.Exit(0)
		}
		fmt.Println(val)
		if val.typ != "array" {
			fmt.Println("Bad Type")
			continue
		}
		if len(val.array) == 0 {
			fmt.Println("Small array")
			continue
		}
		command := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]
		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("wrong command")
			writer.Write(Value{typ: "string", str: ""})
			continue
		}
		if command == "SET" || command == "HSET" {
			aof.Write(val)
		}
		result := handler(args)
		writer.Write(result)

	}

}
