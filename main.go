package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"

	_ "golang.org/x/net/trace"

	"github.com/aniruddha/grpc-websocket-proxy/echoserver"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	grpcAddr  = flag.String("grpcaddr", ":8001", "listen grpc addr")
	httpAddr  = flag.String("addr", ":8000", "listen http addr")
	debugAddr = flag.String("debugaddr", ":8002", "listen debug addr")
)

func CustomMatcher(key string) (string, bool) {
	key = strings.ToLower(key)
	fmt.Println("in custome matcher:", key)

	switch key {
	case "suki_jwt":
		return key, true
	case "suki_organization_id":
		return key, true
	case "suki_user_id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func HeaderForwader(header string) bool {
	header = strings.ToLower(header)
	fmt.Println("in header forwader:", header)

	switch header {
	case "suki_jwt":
		return true
	case "suki_organization_id":
		return true
	case "suki_user_id":
		return true
	default:
		return defaultHeaderForwarder(header)
	}
}

var defaultHeadersToForward = map[string]bool{
	"Origin":  true,
	"origin":  true,
	"Referer": true,
	"referer": true,
}

func defaultHeaderForwarder(header string) bool {
	return defaultHeadersToForward[header]
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := listenGRPC(*grpcAddr); err != nil {
		return err
	}

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(CustomMatcher))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := echoserver.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, *grpcAddr, opts)
	if err != nil {
		return err
	}
	go http.ListenAndServe(*debugAddr, nil)
	fmt.Println("listening")
	fmt.Println("listening and serving to to gateway http on ", *httpAddr)

	http.ListenAndServe(*httpAddr, wsproxy.WebsocketProxy(mux, wsproxy.WithForwardedHeaders(HeaderForwader)))

	return nil
}

func listenGRPC(listenAddr string) error {
	lis, err := net.Listen("tcp", listenAddr)
	fmt.Println("listening to grpc on ", listenAddr)

	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	echoserver.RegisterEchoServiceServer(grpcServer, &Server{})
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("serveGRPC err:", err)
		}
	}()
	return nil

}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
