package main

import (
	"encoding/json"
	"os/exec"
	"fmt"
	"net"
	"os"
	"time"
)

func startServer() error {
	port := ":45455"
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		return err
	}

	listner, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			continue
		}
		go handleClient(conn)
	}
}

var queue = []func(){}

func performCommand(commands []string) error {

	f := func() {
		fmt.Println("start: ", commands)
		out, err := exec.Command(commands[0], commands[1:]...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		fmt.Println("Result: " + string(out))
		fmt.Println("end: ", commands)
		if len(queue) > 0 {
			queue = queue[1:]
			if len(queue) > 0 {
				go queue[0]()
			}
		}
	}
	queue = append(queue, f)

	fmt.Println("queue size: ", len(queue))
	if len(queue) == 1 {
		go f()
	}


	return nil
}

func handleMessage(message string) error {
	fmt.Println(message)

	var commands []string

	err := json.Unmarshal([]byte(message), &commands);
	if err != nil {
		return err
	}

	err = performCommand(commands)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return err
	}

	return nil
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	fmt.Println()
	fmt.Println("client accept!")
	messageBuf := make([]byte, 1024)
	messageLen, err := conn.Read(messageBuf)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	message := string(messageBuf[:messageLen])
	err = handleMessage(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.Write([]byte("ok"))
}

func main() {
	err := startServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
