package main

import (
	"fmt"
	"os"
	"net"
	"net/http"
	"io/ioutil"
	//"bufio"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"crypto/md5"
	"encoding/hex"
	//"github.com/jroimartin/gocui"
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

/*func listenForMessages(conn net.Conn, userName string) {
	for {
		reader, _ := bufio.NewReader(conn).ReadString('\n')
		posts = append(posts, post{username: userName, message: reader, time: "0:00"})
		//fmt.Printf("\033[2K\r%s%s: ", reader, userName)
	}
}*/

func dialService(connString string, userName string, pass string) {
	conn, err := net.Dial("tcp", connString)
	//scanner := bufio.NewScanner(os.Stdin)
	for err != nil {
		fmt.Printf("\033[2K\rerror connecting to server, reconnecting...")
		conn, err = net.Dial("tcp", connString)
	}
	fmt.Fprintf(conn, userName + "|" + pass + "\n")
	//go listenForMessages(conn, userName)
	defer conn.Close()
	initGui(conn, userName, fetchServers())
}

/*func initGui() {
	g, _ := gocui.NewGui(gocui.OutputNormal)
	g.MainLoop()
	maxX, maxY := g.Size()
	g.SetView("viewname", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2)
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("viewname")
		if err != nil {
			// handle error
		}
		v.Clear()
		fmt.Fprintln(v, "Writing from different goroutines")
		return nil
	})
}*/

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
