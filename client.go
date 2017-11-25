package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

func sendMessage(message string) error {
	serverIP := "127.0.0.1"
	serverPort := "45455"

	serverAddr, err := net.ResolveTCPAddr("tcp", serverIP+":"+serverPort)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.Write([]byte(message))

	readBuf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	readlen, err := conn.Read(readBuf)
	if err != nil {
		return err
	}

	fmt.Println("server: " + string(readBuf[:readlen]))
	return nil
}

func main() {
	// fmt.Println(os.Args)
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "at least 1 argument")
		os.Exit(1)
	}

	message, err := json.Marshal(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
	}

	err = sendMessage(string(message))
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
	}
}
