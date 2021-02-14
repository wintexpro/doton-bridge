package tonbindings

import (
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

type ContractContext struct {
	Conn        *client.Client
	Signer      *client.Signer
	WorkchainID null.Int32
}
type ContractAccount struct {
	Boc string `json:"boc"`
}
type ContractMethod struct {
	Call func() (interface{}, error)
	Send func(messageCallback func(event *client.ProcessingEvent)) (interface{}, error)
}
type Contract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
