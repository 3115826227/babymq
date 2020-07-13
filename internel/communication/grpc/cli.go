package grpc

import (
	"github.com/3115826227/babymq/core/register"
	Register "github.com/3115826227/babymq/internel/register"
	"google.golang.org/grpc"
	"sync"
)

type GrpcClient struct {
	servers []register.ServerRegisterMeta
	lock    sync.Mutex
	conn    map[register.ServerRegisterMeta]*grpc.ClientConn
}

var grpcClient *GrpcClient

func init() {
	servers, err := Register.GetEtcdRegisterClient().GetServers()
	if err != nil {
		return
	}
	grpcClient = &GrpcClient{
		servers: servers,
		lock:    sync.Mutex{},
		conn:    make(map[register.ServerRegisterMeta]*grpc.ClientConn),
	}
}

func GetGrpcClient() *GrpcClient {
	return grpcClient
}

func (grpcClient *GrpcClient) GetConn(server register.ServerRegisterMeta) *grpc.ClientConn {
	return grpcClient.conn[server]
}

func (grpcClient *GrpcClient) Connect() (err error) {
	for _, server := range grpcClient.servers {
		if err = grpcClient.connect(server); err != nil {
			return
		}
	}
	return
}

func (grpcClient *GrpcClient) connect(server register.ServerRegisterMeta) (err error) {
	grpcClient.lock.Lock()
	defer grpcClient.lock.Unlock()
	conn, err := grpc.Dial(server.Address, grpc.WithInsecure())
	if err != nil {
		return
	}
	grpcClient.conn[server] = conn
	return
}

func (grpcClient *GrpcClient) Close() (err error) {
	for server := range grpcClient.conn {
		if err = grpcClient.close(server); err != nil {
			return
		}
	}
	return
}

func (grpcClient *GrpcClient) close(server register.ServerRegisterMeta) (err error) {
	grpcClient.lock.Lock()
	defer grpcClient.lock.Unlock()
	err = grpcClient.conn[server].Close()
	if err != nil {
		return
	}
	delete(grpcClient.conn, server)
	return
}
