package command

import (
	"context"
	"net"
	"os"
	"strings"

	"connectrpc.com/connect"

	api "go.octolab.org/ecosystem/sparkle/internal/api/service/v1"
	"go.octolab.org/ecosystem/sparkle/internal/api/service/v1/servicev1connect"
)

const addr = "localhost:52346"

type service struct {
	servicev1connect.UnimplementedServiceHandler
}

func (srv *service) WhoAmI(
	_ context.Context,
	req *connect.Request[api.PingRequest],
) (*connect.Response[api.PongResponse], error) {
	var env map[string]string
	if req.Msg.Env {
		env = getEnv()
	}

	return connect.NewResponse(&api.PongResponse{
		Name:       "sparkle",
		Hostname:   addr,
		Ip:         getIPs(),
		RemoteAddr: req.Peer().Addr,
		Env:        env,
	}), nil
}

func getEnv() map[string]string {
	env := make(map[string]string)

	for _, kv := range os.Environ() {
		key, value, _ := strings.Cut(kv, "=")
		env[key] = value
	}

	return env
}

func getIPs() []string {
	var ips []string

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				ips = append(ips, ip.String())
			}
		}
	}

	return ips
}
