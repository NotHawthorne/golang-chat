package main

import (
	"fmt"
	"os"
	"net"
	"bufio"
)

func clearTerm() {
	for i := 0; i < 90; i++ { fmt.Printf("\033[A\033[2K\r") }
	for i := 0; i < 90; i++ { fmt.Printf("\n") }
}

func listenForMessages(conn net.Conn, userName string) {
	for {
		reader, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("\033[2K\r%s%s: ", reader, userName)
	}
}

func dialService(connString string, userName string) {
	conn, err := net.Dial("tcp", connString)
	scanner := bufio.NewScanner(os.Stdin)
	for err != nil {
		conn, err = net.Dial("tcp", connString)
	}
	fmt.Printf(conn, userName)
	go listenForMessages(conn, userName)
	defer conn.Close()
	for {
		for scanner.Scan() {
			text := scanner.Text()
			if (text != "") {
				fmt.Fprintf(conn, userName + ": " + text + "\n")
				fmt.Printf("\033[A\033[2K")
			}
		}
	}
}

func main() {
	if (len(os.Args) != 3) {
		fmt.Printf("usage: <ip address and port of desired party> <username>\n")
		os.Exit(1)
	}
	clearTerm()
	arg := os.Args[1:]
	go dialService(arg[0], arg[1])
	fmt.Printf(arg[1] + ": ")
	for { }
}
