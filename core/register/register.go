package register

import "encoding/json"

type ServerRegisterClientInterface interface {
	Register(meta ServerRegisterMeta) error
}

type ServerRegisterMeta struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

func (meta *ServerRegisterMeta) ToString() string {
	data, _ := json.Marshal(meta)
	return string(data)
}
