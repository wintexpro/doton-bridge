package tonbindings

import (
	"encoding/json"
	"math/big"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	GiverAbi = "{\"ABI version\":1,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"dest\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint64\"}],\"name\":\"sendGrams\",\"outputs\":[]}]}"
	GiverTvc = ""
)

type Giver struct {
	Abi         client.Abi
	Address     string
	Signer      client.Signer
	Conn        *client.Client
	WorkchainID null.Int32
}

func AbiGiver() client.Abi {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(GiverAbi), &abi.Value); err != nil {
		panic(err)
	}
	return abi
}
func NewGiver(conn *client.Client, signer client.Signer, workchainID null.Int32) (*Giver, error) {
	abi := AbiGiver()

	contract := Giver{
		Abi:         abi,
		Address:     "0:841288ed3b55d9cdafa806807f02a0ae0c169aa5edfe88a789a6482429756a94",
		Conn:        conn,
		Signer:      signer,
		WorkchainID: workchainID,
	}
	return &contract, nil
}

func (c *Giver) SendGrams(dest string, amount *big.Int, messageCallback func(event *client.ProcessingEvent)) (*client.ResultOfSendMessage, error) {
	functionName := "sendGrams"

	callSet := client.CallSet{
		FunctionName: functionName,
		Input: json.RawMessage(
			"{ \"dest\": \"" + dest + "\", \"amount\": " + amount.String() + " }",
		),
	}
	paramsAbiEncodeMessage := client.ParamsOfEncodeMessage{
		Abi:     c.Abi,
		Address: null.StringFrom(c.Address),
		CallSet: &callSet,
		Signer:  c.Signer,
	}
	message, err := c.Conn.AbiEncodeMessage(&paramsAbiEncodeMessage)
	if err != nil {
		panic(err)
	}

	paramsOfSendMessage := client.ParamsOfSendMessage{
		Message: message.Message,
		Abi:     &c.Abi,
	}

	return c.Conn.ProcessingSendMessage(&paramsOfSendMessage, messageCallback)
}
