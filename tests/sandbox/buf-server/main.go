package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	whoami "sparkle/sandbox/buf-server/internal/api/service/v1"
	server "sparkle/sandbox/buf-server/internal/api/service/v1/servicev1connect"
)

const address = "localhost:8080"

//	run buf curl \
//	 --schema api/protobuf \
//	 --data '{"env": true}' \
//	 http://localhost:8080/sparkle.service.v1.Service/WhoAmI
func main() {
	mux := http.NewServeMux()
	path, handler := server.NewServiceHandler(&service{})
	mux.Handle(path, handler)
	fmt.Println("listening on", address)
	_ = http.ListenAndServe(
		address,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

type service struct {
	server.UnimplementedServiceHandler
}

func (s *service) WhoAmI(
	_ context.Context,
	req *connect.Request[whoami.PingRequest],
) (*connect.Response[whoami.PongResponse], error) {
	var env map[string]string
	if req.Msg.Env {
		env = make(map[string]string)
		for _, kv := range os.Environ() {
			key, value, _ := strings.Cut(kv, "=")
			env[key] = value
		}
	}
	return connect.NewResponse(&whoami.PongResponse{
		Name:       "sparkle",
		Hostname:   "localhost",
		Ip:         []string{"127.0.0.1"},
		RemoteAddr: "127.0.0.1",
		Env:        env,
	}), nil
}
