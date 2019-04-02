package main

import (
	"fmt"
	"os"
	"net"
	"bufio"
	"strings"
)

func privateMessage(msg string, usr user, usrlist *[]user) int {
	tmp := strings.Split(msg, " ")
	lst := *usrlist
	if tmp[2] == "/whisper" {
		for i := range lst {
			if lst[i].Name == tmp[3] {
				lst[i].Conn.Write([]byte(usr.Name + " whispers: " + strings.Join(tmp[4:], " ") + "\n"))
				usr.Conn.Write([]byte("to " + lst[i].Name +": " + strings.Join(tmp[4:], " ") + "\n"))
				return 1
			}
		}
		usr.Conn.Write([]byte("invalid user\n"))
		return -1
	}
	return 0
}

func addMsg(msg string, msgarr *[]string) {
	tmp := *msgarr
	tmp = append(tmp, msg)
	*msgarr = tmp
}

func endConn(usr user, usrlist *[]user, msgarr *[]string) {
	ret := *usrlist
	for i := range ret {
		if ret[i].Conn == usr.Conn {
			*usrlist = append(ret[0:i], ret[i + 1:len(ret)]...)
			fmt.Printf("user %s left\n", usr.Name)
			addMsg("user " + usr.Name + " left\n", msgarr)
			usr.Conn.Close()
			return
		}
	}
}

func sendMsgs(msgarr *[]string, conn_list *[]net.Conn) {
	for {
		fmt.Printf("")
		for _, msg := range *msgarr {
			for _, inst := range *conn_list { inst.Write([]byte(msg)) }
			ret := *msgarr
			*msgarr = ret[1:]
		}
	}
}

func waitForHandshake(conn net.Conn, msgarr *[]string, usrlist *[]user) {
	usr := findUser(conn, usrlist)
	for {
		data, _ := bufio.NewReader(usr.Conn).ReadString('\n')
		tmp := strings.TrimSuffix(string(data), "\n")
		tmpVals := strings.Split(tmp, "|")
		usr.Name = tmpVals[0]
		usr.Pass = tmpVals[1]
		updateUser(usr, usrlist)
		ret := validateUser(usr)
		if ret == 0 {
			fmt.Printf("user %s failed validation\n", usr.Name)
			endConn(usr, usrlist, msgarr)
			return
		}
		if ret == 2 {
			fmt.Printf("user %s registered a new account!\n", usr.Name)
			break
		}
		fmt.Printf("user %s joined\n", usr.Name)
		usrlist_tmp := *usrlist
		usrs := ""
		for i := range usrlist_tmp {
			usrs += usrlist_tmp[i].Name + "|"
		}
		conn.Write([]byte(usrs + "\n"))
		conn.Write([]byte("Message Of The Day: Tell Melee players that mozzerella sticks and deoderant sticks are not interchangeable plz thanks\n"))
		addMsg("user " + usr.Name + " joined\n", msgarr)
		break
	}
	go handleConn(conn, msgarr, usrlist)
	return
}

func handleConn(conn net.Conn, msgarr *[]string, usrs *[]user) {
	usr := findUser(conn, usrs)
	for {
		data, err := bufio.NewReader(usr.Conn).ReadString('\n')
		if (err != nil) {
			endConn(usr, usrs, msgarr)
			return
		}
		if privateMessage(string(data), usr, usrs) == 0 {
			*msgarr = append(*msgarr, string(data))
		}
	}
}

func listen(port string) {
	conn := []net.Conn{}
	usrs := []user{}
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
		usrs = append(usrs, user{conn_inst, "", "", ""})
		go waitForHandshake(conn_inst, &test, &usrs)
		go sendMsgs(&test, &conn)
	}
}

func main() {
	listen(os.Args[1])
}
