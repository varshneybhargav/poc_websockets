package main

// import (
// 	"flag"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"os/signal"
// 	"time"

// 	"google.golang.org/protobuf/proto"

// 	"github.com/aniruddha/grpc-websocket-proxy/echoserver"
// 	"github.com/gorilla/websocket"
// )

// var addr = flag.String("addr", "localhost:8000", "http service address")

// func main() {
// 	flag.Parse()
// 	log.SetFlags(0)

// 	interrupt := make(chan os.Signal, 1)
// 	signal.Notify(interrupt, os.Interrupt)

// 	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo1"}
// 	log.Printf("connecting to %s", u.String())

// 	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"suki_jwt": []string{"suki_jwt"}, "suki_user_id": []string{"1ce620da-a218-4dc2-a91b-9531952805fa"}, "suki_organization_id": []string{"11111111-1111-1111-1111-111111111111"}})
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer c.Close()

// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done)
// 		for {
// 			_, message, err := c.ReadMessage()
// 			if err != nil {
// 				log.Println("read:", err)
// 				return
// 			}
// 			var response echoserver.EchoResponse
// 			proto.Unmarshal(message, &response)
// 			log.Println("recv:", response.Message)
// 		}
// 	}()

// 	ticker := time.NewTicker(time.Second)
// 	defer ticker.Stop()
// 	var input string
// 	for {
// 		select {
// 		case <-done:
// 			return
// 		case _ = <-ticker.C:
// 			fmt.Scan(&input)
// 			input = "->" + input
// 			err := c.WriteMessage(websocket.BinaryMessage, []byte(input))
// 			if err != nil {
// 				log.Println("write:", err)
// 				return
// 			}
// 		case <-interrupt:
// 			log.Println("interrupt")

// 			// Cleanly close the connection by sending a close message and then
// 			// waiting (with timeout) for the server to close the connection.
// 			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
// 			if err != nil {
// 				log.Println("write close:", err)
// 				return
// 			}
// 			select {
// 			case <-done:
// 			case <-time.After(time.Second):
// 			}
// 			return
// 		}
// 	}
// }
