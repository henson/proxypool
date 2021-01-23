package rpc

import (
	"github.com/henson/proxypool/pkg/storage"
	"github.com/henson/proxyppol/pkg/rpc/grpc-proxypool"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type ProxyPoolGRPCService struct{}

func (ProxyPoolGRPCService) Get(ctx context.Context, request *grpc_proxypool.Request) (response *grpc_proxypool.Response, err error) {
	var result string
	switch request.Type {
	case "http":
		result = storage.ProxyRandom().Data
	case "https":
		result = storage.ProxyFind("https").Data
	default:
		result = storage.ProxyRandom().Data
	}
	return &grpc_proxypool.Response{Result: result}, nil
}
