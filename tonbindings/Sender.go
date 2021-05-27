package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	SenderAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"destination\",\"type\":\"address\"},{\"name\":\"bounce\",\"type\":\"bool\"},{\"name\":\"value\",\"type\":\"uint128\"},{\"name\":\"data\",\"type\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\"}],\"name\":\"sendData\",\"outputs\":[]},{\"inputs\":[],\"name\":\"constructor\",\"outputs\":[]}],\"header\":[\"time\"]}"
	SenderTvc = "te6ccgECDgEAAbwAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCwwFBA0CkiHbPNMAAY4SgQIA1xgg+QFY+EIg+GX5EPKo3tM/AY4d+EMhuSCfMCD4I4ED6KiCCBt3QKC53pMg+GPg8jTYMNMfAds8+Edu8nwIBgE0ItDXCwOpOADcIccA3CHTHyHdAds8+Edu8nwGAiggghA/VJiouuMCIIIQaLVfP7rjAgkHAiow+EFu4wD4RvJzcfhm0fgA2zx/+GcICgBq7UTQINdJwgGOENP/0z/TANF/+GH4Zvhj+GKOGPQFcAGAQPQO8r3XC//4YnD4Y3D4Zn/4YeIC2DD4QW7jAPpA0gDXDX+V1NHQ03/f1w3/ldTR0NP/39cN/5XU0dDT/9/R+EL4RSBukjBw3rry4Gz4ACIlJcjPhYDKAHPPQM4B+gKAac9Az4HPg8jPkYZ8I24jzwv/Is8L/83JcPsAXwXjAH/4ZwsKACT4QsjL//hDzws/+EbPCwDJ7VQAJu1E0NP/0z/TANF/+GH4Zvhj+GIBCvSkIPShDQAA"
)

type SenderContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type Sender struct {
	Ctx ContractContext
}

func (contract *SenderContract) SendData(destination string, bounce string, value string, data string, destinationChainId string) *ContractMethod {
	input := fmt.Sprintf("{\"destination\": \"%s\" ,\"bounce\": \"%s\" ,\"value\": \"%s\" ,\"data\": \"%s\" ,\"destinationChainId\": \"%s\" }", destination, bounce, value, data, destinationChainId)
	return contract.CallContractMethod("sendData", input)
}
func (c *Sender) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{SenderTvc})
}
func (c *Sender) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(SenderAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *Sender) Address() (string, error) {
	encodeMessage, err := c.DeployEncodeMessage()
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *Sender) New(address string) (*SenderContract, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	if address == "" {
		address, err = c.Address()
		if err != nil {
			return nil, err
		}
	}
	contract := SenderContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *Sender) DeployEncodeMessage() (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         SenderTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{}"))
	callSet := client.CallSet{
		FunctionName: "constructor",
		Input:        params,
	}
	paramsOfEncodeMessage := client.ParamsOfEncodeMessage{
		Abi:       *abi,
		CallSet:   &callSet,
		DeploySet: &deploySet,
		Signer:    *c.Ctx.Signer,
	}
	return c.Ctx.Conn.AbiEncodeMessage(&paramsOfEncodeMessage)
}
func (c *Sender) Deploy(messageCallback func(event *client.ProcessingEvent)) (*SenderContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage()
	if err != nil {
		return nil, err
	}
	paramsOfSendMessage := client.ParamsOfSendMessage{
		Abi:        abi,
		Message:    encodeMessage.Message,
		SendEvents: true,
	}
	_, err = c.Ctx.Conn.ProcessingSendMessage(&paramsOfSendMessage, messageCallback)
	if err != nil {
		return nil, err
	}
	contract := SenderContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *SenderContract) CallContractMethod(methodName string, input string) *ContractMethod {
	return &ContractMethod{
		Call: func() (interface{}, error) {
			result, err := contract.call(methodName, input)
			if err != nil {
				return nil, err
			}
			var output interface{}
			if err := json.Unmarshal(result.Output, &output); err != nil {
				return nil, err
			}
			return output, nil
		},
		Send: func(messageCallback func(event *client.ProcessingEvent)) (interface{}, error) {
			return contract.send(methodName, input, messageCallback)
		},
	}
}
func (contract *SenderContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
	callSet := client.CallSet{
		FunctionName: functionName,
		Input:        json.RawMessage(input),
	}
	paramsAbiEncodeMessage := client.ParamsOfEncodeMessage{
		Abi:     contract.Abi,
		Address: null.StringFrom(contract.Address),
		CallSet: &callSet,
		Signer:  *contract.Ctx.Signer,
	}
	return contract.Ctx.Conn.AbiEncodeMessage(&paramsAbiEncodeMessage)
}
func (contract *SenderContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
	message, err := contract.abiEncodeMessage(functionName, input)
	if err != nil {
		return "", err
	}
	paramsOfSendMessage := client.ParamsOfSendMessage{
		Abi:     &contract.Abi,
		Message: message.Message,
	}
	resultProcessingSendMessage, err := contract.Ctx.Conn.ProcessingSendMessage(&paramsOfSendMessage, messageCallback)
	if err != nil {
		return "", err
	}
	return resultProcessingSendMessage.ShardBlockID, nil
}
func (contract *SenderContract) call(functionName string, input string) (*client.DecodedOutput, error) {
	message, err := contract.abiEncodeMessage(functionName, input)
	if err != nil {
		return nil, err
	}
	filter := json.RawMessage(fmt.Sprintf("{\"id\":{\"eq\": \"%s\"}}", contract.Address))
	params := client.ParamsOfQueryCollection{
		Collection: "accounts",
		Filter:     filter,
		Limit:      null.Uint32From(1),
		Result:     "boc",
	}
	res, err := contract.Ctx.Conn.NetQueryCollection(&params)
	if err != nil {
		return nil, err
	}
	var account ContractAccount
	if err := json.Unmarshal(res.Result[0], &account); err != nil {
		return nil, err
	}
	paramsOfRunTvm := client.ParamsOfRunTvm{
		Abi:     &contract.Abi,
		Account: account.Boc,
		Message: message.Message,
	}
	resultOfRunTvm, err := contract.Ctx.Conn.TvmRunTvm(&paramsOfRunTvm)
	if err != nil {
		return nil, err
	}
	return resultOfRunTvm.Decoded, nil
}
