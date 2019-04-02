
package main

import (
	"fmt"
	"log"
	"time"
	"net"
	"bufio"
	"strings"

	"github.com/marcusolsson/tui-go"
)

type post struct {
	username string
	message  string
	time     string
}

type sidebar struct {
	box		*tui.Box
	labels	int
	users	[]string
}

func listenForMessages(conn net.Conn, userName string, box **tui.Box, ui *tui.UI, usersBox **tui.Box) {
	test := 0
	var users []string
	users = append(users, userName)
	for {
		i := *box
		reader, _ := bufio.NewReader(conn).ReadString('\n')
		str := strings.Split(reader, " ")
		if (test == 0 || str[0] == "user") {
			if test == 0 { users = strings.Split(reader, "|") }
			if test != 0 { users = append(users, str[1]) }
			x := *usersBox
			if (test == 0) {
				for e := range users {
					if !strings.Contains(users[e], "\n") {
						x.Append(tui.NewLabel(strings.TrimSuffix(users[e], "\n")))
					}
					if (strings.Contains(users[e], "\n")) {
						users = append(users[0:e], users[e + 1:]...)
					}
				}
			}
			if test != 0 {
				if (str[2] == "joined\n" && str[1] != userName) {
					x.Append(tui.NewLabel(str[1]))
				}
				if (str[2] == "left\n") {
					i := []string{}
					for e := range users {
						exists := 0
						for n := range i {
							if i[n] == users[e] { exists++ }
						}
						if users[e] == str[1] { x.Remove(e + 5) }
						if exists <= 0 && users[e] != str[1] { i = append(i, users[e]) }
					}
					users = i
				}
			}
			*usersBox = x
		}
		if test != 0 {
			i.Append(tui.NewHBox(
				tui.NewLabel(strings.TrimSuffix(reader, "\n")),
			))
			if (test > 100) { i.Remove(0) }
			*box = i
		}
		test++
		tmp := *ui
		go tmp.Update(func() {
			//...
		})
	}
}

func initGui(conn net.Conn, userName string, servers []string) {
	sidebar := tui.NewVBox(
		tui.NewLabel("SERVERS"),
		tui.NewLabel(servers[0]),
		tui.NewLabel(servers[1]),
		tui.NewSpacer(),
		tui.NewLabel("USERS"),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		fmt.Fprintf(conn, "%s <%s> %s\n", time.Now().Format("15:04"), userName, e.Text())
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	go listenForMessages(conn, userName, &history, &ui, &sidebar)
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
