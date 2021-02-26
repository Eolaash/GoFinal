package binance

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

// Session - main control struct with defined interface
type Session struct {
	URL         url.URL
	webConn     *websocket.Conn
	Err         error
	isConnected bool
}

// Connect - connect by websocket to inURL and writing Session struct
func (inSession *Session) Connect(inURL url.URL) {
	inSession.URL = inURL
	inSession.webConn, _, inSession.Err = websocket.DefaultDialer.Dial(inSession.URL.String(), nil)

	// On error
	if inSession.Err != nil {
		fmt.Println("Connection TRY failed >> ", inSession.Err.Error())
		return
	}

	// Connection success
	fmt.Println("Connection OK")
	inSession.isConnected = true
}

// Disconnect - disconnect from active websocket
func (inSession *Session) Disconnect() {

	// Qiuck return (no action needed - already no active connection)
	if inSession.isConnected != true || inSession.webConn == nil {
		inSession.isConnected = false
		fmt.Println("No connection active.")
		return
	}

	// Disconnect
	inSession.Err = inSession.webConn.Close()

	// On error
	if inSession.Err != nil {
		fmt.Println("Disconnect failed >> ", inSession.Err.Error())
		return
	}

	// Disconnect success
	inSession.isConnected = false
	fmt.Println("Connection DROPPED")
}

// IsConnected - quick check connection (need to be updated to ping-pong coat - check is it alive or dead)
func (inSession *Session) IsConnected() bool {

	// quick check
	if inSession.webConn == nil {
		inSession.isConnected = false
	}

	if inSession.isConnected == true {
		return true
	} else {
		return false
	}
}

// Subscribe - Subscribe on predefined info channel
func (inSession *Session) Subscribe() bool {

	// SubScribe to channel (PREDEFINED - bad move? should be more flexible as IN_VAR)
	sendString := []byte(`{"method": "SUBSCRIBE", "params": ["btcusdt@aggTrade"], "id": 1}`)

	// Quick check connection
	if !inSession.IsConnected() {
		return false
	}

	// send message
	inSession.Err = inSession.webConn.WriteMessage(websocket.TextMessage, sendString)

	// On error
	if inSession.Err != nil {
		fmt.Println("Subscription error >> ", inSession.Err.Error())
		return false
	}

	// success return
	fmt.Println("Subscription OK")
	return true
}

// ReadData - for testing returns STRING
func (inSession *Session) ReadData() string {

	var tMessage []byte

	// receive message
	_, tMessage, inSession.Err = inSession.webConn.ReadMessage()
	if inSession.Err != nil {
		return "Error on READ! >> " + inSession.Err.Error()
	}

	return string(tMessage)
}
