package rpc

/*
  rpc客户端接口
*/
type RPCClientInterface interface {
	//rpc连接
	Connect() error
	//rpc断开
	Close() error
}
