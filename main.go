package main

import (
	"FinalTask/testpublic"
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	var wg sync.WaitGroup

	fmt.Println(testpublic.TestPublic())

	tPoloURL := url.URL{
		Scheme: "wss",
		Host:   "api2.poloniex.com",
		Path:   "/",
	}

	//tPoloURL := url.URL{
	//	Scheme: "wss",
	//	Host:   "stream.binance.com:9443",
	//	Path:   "/ws",
	//}

	c, _, err := websocket.DefaultDialer.Dial(tPoloURL.String(), nil)
	if err != nil {
		fmt.Println("Oops! >> ", err.Error())
		return
	}

	fmt.Println(sendMessage(c, []byte(`{"command": "subscribe", "channel": 1002}`)))
	//fmt.Println(sendMessage(c, []byte(`{"method": "SUBSCRIBE", "params": ["btcusdt@aggTrade"], "id": 1}`)))
	fmt.Println("Read #INIT-SUB > ", readMessage(c))

	tConsoleRead := bufio.NewReader(os.Stdin)

	tFlag := false

	for {
		defer wg.Done()

		// tFlag for single start
		if !tFlag {
			fmt.Println("Listening subscribtion.. ")
			tReads := 0
			tFlag = true
			wg.Add(1)
			go func() {
				for {
					tReads++
					fmt.Println("Read #", tReads, " > ", readMessage(c))
					time.Sleep(time.Second * 2)
				}
			}()
		}

		//fmt.Print("-> ")
		tText, _ := tConsoleRead.ReadString('\n')

		// CRLF kill
		tText = strings.Replace(tText, "\r\n", "", -1)

		if "kill" == tText {
			fmt.Println("Read #LAST > ", readMessage(c))
			break
		}
	}

	fmt.Println(sendMessage(c, []byte(`{"command": "unsubscribe", "channel": 1002}`)))
	//fmt.Println(sendMessage(c, []byte(`{"method": "UNSUBSCRIBE", "params": ["btcusdt@aggTrade"], "id": 132}`)))
	fmt.Println("Read #UNSUB > ", readMessage(c))
	fmt.Println("Read #UNSUB > ", readMessage(c))
	fmt.Println("Read #UNSUB > ", readMessage(c))

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
