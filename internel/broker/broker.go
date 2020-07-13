package broker

import (
	"github.com/3115826227/babymq/core/communication/rpc"
	"github.com/3115826227/babymq/core/register"
	"github.com/3115826227/babymq/internel/election/raft"
)

type RegisterType uint32

const (
	Etcd RegisterType = iota
)

type CommunicationType uint32

const (
	RPC CommunicationType = iota

	HTTP
)

type Broker struct {
	server   register.ServerRegisterMeta
	provider *raft.ElectionRaftProvider

	message chan interface{}

	communicationType CommunicationType
	rpcServer         rpc.RPCServerInterface
	rpcClient         rpc.RPCClientInterface

	registerType RegisterType
	register     register.ServerRegisterClientInterface
}

func (broker *Broker) Start() {

	broker.register.Register(broker.server)

	broker.rpcServer.Start()
	broker.rpcClient.Connect()

	broker.provider = raft.GetElectionRaftProvider()
	broker.provider.Run()

}

func (broker *Broker) CreateTopic(topic string) (err error) {
	return
}

func (broker *Broker) GetTopics() {
	return
}
