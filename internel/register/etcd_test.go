package register

import (
	"fmt"
	"github.com/3115826227/babymq/core/register"
	"testing"
)

func TestGetEtcdRegisterClient(t *testing.T) {
	client := GetEtcdRegisterClient()
	err := client.Connect()
	fmt.Println(err)
	meta := register.ServerRegisterMeta{
		ID:      "2",
		Address: "http://127.0.0.1:2332",
		Name:    "broker-2",
	}
	err = client.Register(meta)
	fmt.Println(err)
	resp, err := client.list(EtcdPrefix)
	fmt.Println(resp.Data, err)
	client.Close()
}
