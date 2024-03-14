package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aniruddha/grpc-websocket-proxy/echoserver"
	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	echoserver.EchoServiceServer
}

func (s *Server) Stream(stream echoserver.EchoService_StreamServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("failed to extract metadata")
	}

	for key, values := range md {
		log.Printf("Header: %s\n", key)
		for _, value := range values {
			log.Printf("  Value: %s\n", value)
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			req, err := stream.Recv()
			if err != nil {
				log.Println("error in stream recv:", err)
				return
			}
			log.Println("msg received on server is :", req.Message)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			if err := stream.Send(&echoserver.EchoResponse{
				Message: "server stream",
			}); err != nil {
				log.Println("error in stream send:", err)
				return
			}
		}
	}()
	return nil
}

func (s *Server) Echo(srv echoserver.EchoService_EchoServer) error {
	md, ok := metadata.FromIncomingContext(srv.Context())
	if !ok {
		return fmt.Errorf("failed to extract metadata")
	}

	for key, values := range md {
		log.Printf("Header: %s\n", key)
		for _, value := range values {
			log.Printf("  Value: %s\n", value)
		}
	}
	for {
		req, err := srv.Recv()
		if err != nil {
			return fmt.Errorf("error is due to %w", err)
		}

		fmt.Println("msg received on server is :", req.Message)
		if err := srv.Send(&echoserver.EchoResponse{
			Message: string(req.Message) + "!",
		}); err != nil {
			return err
		}

	}
}

func (s *Server) Heartbeats(srv echoserver.EchoService_HeartbeatsServer) error {

	md, ok := metadata.FromIncomingContext(srv.Context())
	if !ok {
		return fmt.Errorf("failed to extract metadata")
	}

	for key, values := range md {
		log.Printf("Header: %s\n", key)
		for _, value := range values {
			log.Printf("  Value: %s\n", value)
		}
	}

	go func() {
		for {
			_, err := srv.Recv()
			if err != nil {
				log.Println("Recv() err:", err)
				return
			}
			log.Println("got hb from client")
		}
	}()
	t := time.NewTicker(time.Second * 1)
	for {
		log.Println("sending hb")
		hb := &echoserver.Heartbeat{
			Status: echoserver.Heartbeat_OK,
		}
		b := new(bytes.Buffer)
		if err := (&jsonpb.Marshaler{}).Marshal(b, hb); err != nil {
			log.Println("marshal err:", err)
		}
		log.Println(string(b.Bytes()))
		if err := srv.Send(hb); err != nil {
			return err
		}
		<-t.C
	}
	return nil
}

func (s *Server) Sample(ctx context.Context, in *echoserver.EchoRequest) (*echoserver.Empty, error) {
	log.Printf("greet many times was invoked with %v\n", in)
	return &echoserver.Empty{}, nil

}
