package main

import (
	"fmt"
	"os"
	"net"
	"bufio"
)

func sendMsgs(msgarr *[]string, conn_list *[]net.Conn) {
	for {
		for _, msg := range *msgarr {
			for _, inst := range *conn_list {
				inst.Write([]byte(msg))
				fmt.Printf("sent message\n")
			}
			ret := *msgarr
			*msgarr = ret[1:]
		}
	}
}

func handleConn(conn net.Conn, msgarr *[]string) {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if (err != nil) {
			fmt.Printf("user left\n")
			return
		}
		*msgarr = append(*msgarr, string(data))
	}
}

func listen(port string) {
	conn := []net.Conn{}
	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Printf("error: unable to listen on given port " + port)
		return
	}
	defer ln.Close()
	for {
		test := []string{}
		conn_inst, err := ln.Accept()
		if err != nil { fmt.Printf("error\n") }
		conn = append(conn, conn_inst)
		go handleConn(conn_inst, &test)
		go sendMsgs(&test, &conn)
	}
}

func main() {
	arg := os.Args[1:]
	listen(arg[0])
}
