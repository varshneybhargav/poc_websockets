package main

import (
	"encoding/json"
	"flag"
	"github.com/aniruddha/grpc-websocket-proxy/echoserver"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
)

var addr = flag.String("addr", "localhost:8000", "http service address")

// type Result struct {
// 	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
// }

type Response struct {
	Result echoserver.EchoResponse
}

type Audio struct {
	Data []AudioData
}

type AudioData struct {
	AudioBytes string `json:"audioData"`
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Attempt to establish a WebSocket connection.
	err := establishWebSocketConnection(interrupt)
	if err != nil {
		log.Println("Failed to establish WebSocket connection:", err)
		// You can add additional error handling logic here if needed.
	}

	// Sleep before attempting to reconnect (you can adjust the delay).

}

func establishWebSocketConnection(interrupt chan os.Signal) error {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	// http.Header{"suki_jwt": []string{"suki_jwt"}, "suki_user_id": []string{"1ce620da-a218-4dc2-a91b-9531952805fa"}, "suki_organization_id": []string{"11111111-1111-1111-1111-111111111111"}}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"suki_jwt": []string{"suki_jwt"}, "suki_user_id": []string{"1ce620da-a218-4dc2-a91b-9531952805fa"}, "suki_organization_id": []string{"11111111-1111-1111-1111-111111111111"}})
	if err != nil {
		log.Println("dial:", err)
		return err
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var response Response
			json.Unmarshal(message, &response)
			// log.Printf("recv: %s %T", response.Result, response.Result)

			log.Println("recv:", response.Result, reflect.TypeOf(response.Result))
			// log.Println("recv:", message)
			// log.Printf("recv: %s", message)
		}
	}()

	data, err := os.ReadFile("/Users/bhargavvarshney/ms-ext-gateway/source/internal/test_wsclient/1min_audio.json")
	if err != nil {
		log.Println("error reading file", err)
		return err
	}

	var payload Audio
	err = json.Unmarshal(data, &payload)
	if err != nil {
		log.Println("cannot unmarshal json", err.Error())
		return err
	}

	for i := 0; i < 5; i++ {
		message := &echoserver.EchoRequest{
			Message: []byte(payload.Data[i].AudioBytes),
		}

		// message := &echoserver.EchoRequest{
		// 	Message: []byte("hello! there"),
		// }

		byteInput, _ := json.Marshal(message)
		err = c.WriteMessage(websocket.TextMessage, byteInput)
		// err := c.WriteMessage(websocket.TextMessage, []byte("{\"message\":\"hey hi hello si\"}"))
		if err != nil {
			log.Println("write:", err)
			return err
		}
	}

	return nil
}
