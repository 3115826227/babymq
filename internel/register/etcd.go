package register

import (
	"context"
	"encoding/json"
	"github.com/3115826227/babymq/core/register"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

const (
	EtcdPrefix       = "/babymq/etcd/"
	EtcdServerPrefix = "/babymq/etcd/server/"
)

type EtcdRegisterClient struct {
	Address       []string `json:"address"`
	client        *clientv3.Client
	kv            clientv3.KV
	ctx           context.Context
	lease         clientv3.Lease
	leaseID       clientv3.LeaseID
	leaseTTL      int64
	leaseRespChan <-chan *clientv3.LeaseKeepAliveResponse
	Username      string
	Password      string
}

type EtcdRegisterClientResp struct {
	lock sync.Mutex
	Data map[string]string
}

var etcdClient *EtcdRegisterClient

func init() {
	etcdClient = &EtcdRegisterClient{
		Address:  []string{"http://127.0.0.1:23791"},
		leaseTTL: 20,
	}
}

func GetEtcdRegisterClient() *EtcdRegisterClient {
	return etcdClient
}

//连接etcd
func (etcdClient *EtcdRegisterClient) Connect() (err error) {
	etcdClient.client, err = clientv3.New(clientv3.Config{
		Endpoints:   etcdClient.Address,
		DialTimeout: 5 * time.Second,
		TLS:         nil,
		Username:    etcdClient.Username,
		Password:    etcdClient.Password,
	})
	etcdClient.kv = clientv3.NewKV(etcdClient.client)
	etcdClient.ctx = context.Background()
	return
}

func (etcdClient *EtcdRegisterClient) Close() (err error) {
	return etcdClient.client.Close()
}

//注册服务
func (etcdClient *EtcdRegisterClient) Register(meta register.ServerRegisterMeta) (err error) {
	lease := clientv3.NewLease(etcdClient.client)
	leaseResp, err := lease.Grant(context.TODO(), etcdClient.leaseTTL)
	if err != nil {
		return
	}
	etcdClient.leaseID = leaseResp.ID
	go etcdClient.listenerLease()
	_, err = etcdClient.kv.Put(etcdClient.ctx, EtcdPrefix+meta.ID+"/"+meta.Name, meta.ToString(), clientv3.WithLease(leaseResp.ID))
	return
}

func (etcdClient *EtcdRegisterClient) GetServers() (servers []register.ServerRegisterMeta, err error) {
	etcdResp, err := etcdClient.list(EtcdServerPrefix)
	if err != nil {
		return
	}
	servers = make([]register.ServerRegisterMeta, 0)
	for _, value := range etcdResp.Data {
		var meta register.ServerRegisterMeta
		if err = json.Unmarshal([]byte(value), &meta); err != nil {
			return
		}
		servers = append(servers, meta)
	}
	return
}

func (etcdClient *EtcdRegisterClient) list(prefix string) (etcdResp EtcdRegisterClientResp, err error) {
	resp, err := etcdClient.kv.Get(etcdClient.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return
	}
	etcdResp.Data = make(map[string]string)
	for _, value := range resp.Kvs {
		if value != nil {
			etcdResp.addValue(string(value.Key), string(value.Value))
		}
	}
	return
}

func (etcdClient *EtcdRegisterClient) leaseKeepAlive() {
	ctx, _ := context.WithCancel(context.TODO())
	leaseChan, err := etcdClient.lease.KeepAlive(ctx, etcdClient.leaseID)
	if err != nil {
		return
	}
	etcdClient.leaseRespChan = leaseChan
}

func (etcdClient *EtcdRegisterClient) listenerLease() {
	for {
		select {
		case leaseKeepResp := <-etcdClient.leaseRespChan:
			if leaseKeepResp == nil {
				etcdClient.leaseKeepAlive()
				continue
			}
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func (etcdResp *EtcdRegisterClientResp) addValue(key, value string) {
	etcdResp.lock.Lock()
	defer etcdResp.lock.Unlock()
	etcdResp.Data[key] = value
}
