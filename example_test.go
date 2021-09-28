package wsstomp_test

import (
	"context"
	"log"
	"time"

	wsstomp "github.com/SoMuchForSubtlety/ws-stomp"
	"github.com/go-stomp/stomp/v3"
)

// ExampleEstablishConnection demonstrates how to use this library to talk STOMP over a websocket connection
func Example_establishConnection() {
	// timeout if the connection isn't established after ten seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// connect to websocket
	conn, err := wsstomp.Connect(ctx, "wss://test.com/ws", nil)
	cancel()
	if err != nil {
		log.Printf("error during WS connect: %v", err)
		return
	}

	// init STOMP connection using the websocket connections
	stompConn, err := stomp.Connect(conn)
	if err != nil {
		conn.Close()
		log.Printf("error during STOMP connect: %v", err)
		return
	}
	defer func() {
		err = stompConn.Disconnect()
		if err != nil {
			log.Printf("error during STOMP disconnect: %v", err)
		}
	}()

	// send a message
	err = stompConn.Send("/queue/a", "text/plain", []byte("hello world!"))
	if err != nil {
		log.Println(err)
		return
	}
}
