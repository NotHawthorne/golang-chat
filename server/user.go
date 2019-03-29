package main

import (
	"net"
	"os"
	"bufio"
	"fmt"
	"strings"
)

type user struct {
	Conn net.Conn
	Name string
	Pass string
	IP string
}

func findUser(conn net.Conn, usrlist *[]user) user {
	tmp := *usrlist
	for i := range tmp {
		if conn == tmp[i].Conn {
			return tmp[i]
		}
	}
	return (user{nil, "", "", ""})
}

func updateUser(usr user, usrlist *[]user) {
	tmp := *usrlist
	for i := range tmp {
		if usr.Conn == tmp[i].Conn {
			tmp[i] = usr
			*usrlist = tmp
			return
		}
	}
}

func validateUser(usr user) int {
	file, err := os.Open("users.db")
	if err != nil { fmt.Printf("error\n") }
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		vals := strings.Split(line, "|")
		if (vals[0] == usr.Name) {
			if (vals[1] == usr.Pass) { return 1 }
			return 0
		}
	}
	file.Close()
	file, _ = os.OpenFile("users.db", os.O_WRONLY|os.O_APPEND, 0600)
	_, err = file.WriteString(usr.Name + "|" + usr.Pass + "\n")
	if (err != nil) { fmt.Printf("error\n") }
	usr.Conn.Write([]byte("registered new account!\n"))
	return 2
}
