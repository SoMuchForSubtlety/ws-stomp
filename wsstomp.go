package wsstomp

import (
	"context"
	"time"

	"nhooyr.io/websocket"
)

type WebsocketSTOMP struct {
	connection *websocket.Conn

	readerBuffer []byte
	writeBuffer  []byte
}

const (
	NullByte     = 0x00
	LineFeedByte = 0x0a
)

// Read messages from the websocket connection until the provided array is full.
// Any surplus data is preserved for the next Read call
func (w *WebsocketSTOMP) Read(p []byte) (int, error) {
	// if we have no more data, read the next message from the websocket
	if len(w.readerBuffer) == 0 {
		_, msg, err := w.connection.Read(context.Background())
		if err != nil {
			return 0, err
		}
		w.readerBuffer = msg
	}

	n := copy(p, w.readerBuffer)
	w.readerBuffer = w.readerBuffer[n:]
	return n, nil
}

// Write to the websocket.
//
// The written data is held back until a full STOMP frame has been written,
// then a WS message is sent.
func (w *WebsocketSTOMP) Write(p []byte) (int, error) {
	var err error
	w.writeBuffer = append(w.writeBuffer, p...)
	// if we reach a null byte or the entire message is a linefeed (heartbeat), send the message
	if p[len(p)-1] == NullByte || (len(w.writeBuffer) == 1 && len(p) == 1 && p[0] == LineFeedByte) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		err = w.connection.Write(ctx, websocket.MessageText, w.writeBuffer)
		cancel()
		// TODO: preserve write buffer if write fails?
		w.writeBuffer = []byte{}
	}
	return len(p), err
}

func (w *WebsocketSTOMP) Close() error {
	return w.connection.Close(websocket.StatusNormalClosure, "terminating connection")
}

// Establish a websocket connection with the provided URL.
// The context parameter will only be used for the connection handshake,
// and not for the full lifetime of the connection.
func Connect(ctx context.Context, url string, options *websocket.DialOptions) (*WebsocketSTOMP, error) {
	con, _, err := websocket.Dial(ctx, url, options)
	return &WebsocketSTOMP{
		connection: con,
	}, err
}
