package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/3115826227/babymq/internel/communication/service"
	pbservices "github.com/3115826227/babymq/internel/communication/service/pbservice"
	"github.com/3115826227/babymq/internel/election/raft"
	"google.golang.org/grpc"
	"net"
)

type GrpcServerProvider struct {
	listenerPort int    `json:"listener_port"`
	tls          bool   `json:"tls"`
	tlsCertFile  string `json:"tls_cert_file"`
	tlsKeyFile   string `json:"tls_key_file"`

	grpcServer *grpc.Server
}

var grpcProvider *GrpcServerProvider

func init() {
	grpcProvider = &GrpcServerProvider{
		listenerPort: 5337,
		tls:          false,
		grpcServer:   grpc.NewServer(),
	}
}

// 获取本机IP
func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetGrpcProvider() *GrpcServerProvider {
	return grpcProvider
}

func (grpcServerProvider *GrpcServerProvider) GetServer() *grpc.Server {
	return grpcProvider.grpcServer
}

func (grpcProvider *GrpcServerProvider) Start() (err error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", localIP(), grpcProvider.listenerPort))
	if err != nil {
		return
	}
	if err = grpcProvider.grpcServer.Serve(listener); err != nil {
		return
	}
	grpcProvider.RegisterService()
	go grpcProvider.listener()
	return
}

func (grpcProvider *GrpcServerProvider) RegisterService() {
	pbservices.RegisterVoteServiceServer(grpcProvider.GetServer(), &service.VoteService{})
}

func (grpcProvider *GrpcServerProvider) Stop() {
	grpcProvider.grpcServer.Stop()
}

func (grpcServerProvider *GrpcServerProvider) listener() {
	for {
		select {
		case request := <-raft.GetElectionRaftProvider().SendRequest:
			fmt.Println(request)
			switch request.MessageType {
			case 1:
				var voteRequest = &pbservices.VoteRequest{}
				if err := json.Unmarshal(request.Content, voteRequest); err != nil {
					continue
				}
				client := pbservices.NewVoteServiceClient(grpcProvider.GetServer())
				response, err := client.Request(context.Background(), voteRequest)
				if err != nil {
					continue
				}
				raft.GetElectionRaftProvider().ReceiveResponse <- response
			}
		}
	}
}
