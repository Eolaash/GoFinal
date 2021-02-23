package main

import (
	"FinalTask/testpublic"
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println(testpublic.TestPublic())

	tPoloURL := url.URL{
		Scheme: "wss",
		Host:   "api2.poloniex.com",
		Path:   "/",
	}

	c, _, err := websocket.DefaultDialer.Dial(tPoloURL.String(), nil)
	if err != nil {
		// handle error
	}

	fmt.Println(sendMessage(c, []byte(`{"command": "subscribe", "channel": 1002}`)))

	tConsoleRead := bufio.NewReader(os.Stdin)

	tFlag := false

	for {

		// tFlag for single start
		if !tFlag {
			tFlag = true
			go func() {
				for {
					fmt.Println(readMessage(c))
					time.Sleep(time.Second * 2)
				}
			}()
		}

		//fmt.Print("-> ")
		tText, _ := tConsoleRead.ReadString('\n')

		// CRLF kill
		tText = strings.Replace(tText, "\r\n", "", -1)

		if "kill" == tText {
			break
		}
	}

	fmt.Println(sendMessage(c, []byte(`{"command": "unsubscribe", "channel": 1002}`)))

	//outMessage := []byte(`{"command": "subscribe", "channel": 1002}`)
	//outMessage = []byte(`{"command": "unsubscribe", "channel": 1002}`)
}

func sendMessage(webConn *websocket.Conn, sendString []byte) string {

	// send message
	err := webConn.WriteMessage(websocket.TextMessage, sendString)
	if err != nil {
		return "Error on WRITE!"
	}

	return "OK"
}

func readMessage(webConn *websocket.Conn) string {

	// receive message
	_, message, err := webConn.ReadMessage()
	if err != nil {
		return "Error on READ1!"
	}

	return string(message)
}
