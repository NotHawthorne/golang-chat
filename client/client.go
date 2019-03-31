package main

import (
	"fmt"
	"os"
	"net"
	"net/http"
	"io/ioutil"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"crypto/md5"
	"encoding/hex"
)

func fetchServers() []string {
	resp, _ := http.Get("https://raw.githubusercontent.com/NotHawthorne/golang-chat/master/server/index.txt")
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	ret := strings.Split(string(bodyBytes), "\n")
	return ret
}

func clearTerm() {
	for i := 0; i < 90; i++ { fmt.Printf("\033[A\033[2K\r") }
	for i := 0; i < 90; i++ { fmt.Printf("\n") }
}

func dialService(connString string, userName string, pass string) {
	conn, err := net.Dial("tcp", connString)
	for err != nil {
		fmt.Printf("\033[2K\rerror connecting to server, reconnecting...")
		conn, err = net.Dial("tcp", connString)
	}
	fmt.Fprintf(conn, userName + "|" + pass + "\n")
	defer conn.Close()
	initGui(conn, userName, fetchServers())
}

func main() {
	if (len(os.Args) != 2) {
		fmt.Printf("usage: ./program <username>\n")
		os.Exit(1)
	}
	fmt.Printf("Password: ")
	bytePass, _ := terminal.ReadPassword(int(syscall.Stdin))
	hasher := md5.New()
	hasher.Write(bytePass)
	passHash := hex.EncodeToString(hasher.Sum(nil))
	clearTerm()
	servers := fetchServers()
	mainServer := strings.Split(servers[0], "|")
	arg := os.Args[1:]
	fmt.Printf(arg[0] + ": ")
	dialService(mainServer[0], arg[0], passHash)
}
