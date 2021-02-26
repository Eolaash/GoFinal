// MAIN logic is:
// We have multiple datasources with common interface called Session (Connect, Disconnect, Subscribe, Unsubscribe?, IsConnected?(really need?), ReadData(reading channel and parse data to common datatype))
// Main program AS BASE variant will be CYCLED console APP with waiting some KEYBINDING or INPUT(from keyboard)
// In main (waiting to quit cycle) part callin goroutines which read DATASOURCE channels (wss, http) through common interface and WRITING it to common DATASTORAGE interface (DB, File, Memory ect) through CHANNELS(buffered?)
// Many aspects of flowcontrol - can't predict speed.. write interfaces to DATASTORAGE? never worked with it; OVERLOAD, OVERFLOW, DROPPED CONNECTIONS and ect;
// SYNC and DATA protection as main - buffering to prevent dataloss? but memory depends (leaks? progession?)? SYNC speed or READ WRITE? Selfcontrol analyzer?
// To predict READ - need more research in WEBSOCKET works (datasheets of API)
package main

import (
	"FinalTask/binance"
	"FinalTask/poloniex"
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Session is common interface for websocket multidata sources
type Session interface {
	Connect(inURL url.URL) // replace with Init func >> URL - Reconnections attemts limit - Subscribe strings and ect - Timeouts and ect -> so better to use INIT func (once) and RECONNECT func as part of cycled algorythm of reading
	Disconnect()
	IsConnected() bool
	Subscribe() bool
	ReadData() string
}

func main() {
	var wg sync.WaitGroup
	tSessions := []Session{new(poloniex.Session), new(binance.Session)}
	tConsoleRead := bufio.NewReader(os.Stdin)
	tFlag := false

	tWorkURLs := []url.URL{{
		Scheme: "wss",
		Host:   "api2.poloniex.com",
		Path:   "/",
	}, {
		Scheme: "wss",
		Host:   "stream.binance.com:9443",
		Path:   "/ws",
	}}

	fmt.Println("Starting. Print <kill> to stop app...")

	// CONSOLE SCAN
	for {
		defer wg.Done()
		for tIndex := range tSessions {
			defer tSessions[tIndex].Disconnect()
		}

		// tFlag for JUST single start goroutines
		if !tFlag {
			tFlag = true
			wg.Add(1)

			// FOR EACH Session should start goroutines (data exchange channel is single? DATASTORAGE is single mean - is multiple channels? SYNC?)
			for tIndex := range tSessions {
				go func(inIndex int) {
					for {
						// Connection or reconnection
						if !tSessions[inIndex].IsConnected() {
							tSessions[inIndex].Connect(tWorkURLs[inIndex]) // subscribe should check CONNECTION
							tSessions[inIndex].Subscribe()                 // test subs (no checks of ALREADY SUBSCRIBED (parse data neeed to complete IsSubscribed func))
						}

						// Reader
						if tSessions[inIndex].IsConnected() {
							fmt.Println(tSessions[inIndex].ReadData())
						}

						// Saver
						// ...

						time.Sleep(time.Second * 2) // for testing reason
					}
				}(tIndex)
			}
		}

		//fmt.Print("-> ")
		tText, _ := tConsoleRead.ReadString('\n')

		// CRLF kill
		tText = strings.Replace(tText, "\r\n", "", -1)

		if "kill" == tText {
			fmt.Println("Quiting app!")
			break
		}
	}
}
