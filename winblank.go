package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func sleepWindows() {
	c := exec.Command("cmd", "/C", "shutdown", "/h")

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}

var interrupt chan os.Signal

func wsConnect() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: "dak.jpoles1.com", Path: "/"}
	//log.Printf("Connecting to %s\n", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("Error connecting to server - ", err)
		time.Sleep(1500 * time.Millisecond)
		return
	}
	log.Println("Connected to DAK server.")
	defer c.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("WS read error - ", err)
				return
			}
			log.Printf("Received: %s", message)
			if string(message) == "alexaevent:computer:power:off" {
				sleepWindows()
			}
		}
	}()
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Fatal("Received OS signal, closing goblank client.")
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func main() {
	for {
		wsConnect()
	}
}
