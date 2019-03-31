
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

func listenForMessages(conn net.Conn, userName string, box **tui.Box, ui *tui.UI) {
	for {
		i := *box
		reader, _ := bufio.NewReader(conn).ReadString('\n')
		i.Append(tui.NewHBox(
			tui.NewLabel(strings.TrimSuffix(reader, "\n")),
		))
		*box = i
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
		tui.NewLabel(""),
		tui.NewLabel("DIRECT MESSAGES"),
		tui.NewLabel("slackbot"),
		tui.NewSpacer(),
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

	go listenForMessages(conn, userName, &history, &ui)
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
