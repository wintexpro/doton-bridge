package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	ReceiverAbi = "{\"ABI version\":2,\"data\":[],\"events\":[{\"inputs\":[{\"name\":\"data\",\"type\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"DataReceived\",\"outputs\":[]}],\"functions\":[{\"inputs\":[{\"name\":\"data\",\"type\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\"}],\"name\":\"receiveData\",\"outputs\":[]},{\"inputs\":[{\"name\":\"destinationChainId\",\"type\":\"uint256\"}],\"name\":\"getNonceByChainId\",\"outputs\":[{\"name\":\"nonce\",\"type\":\"uint256\"}]},{\"inputs\":[],\"name\":\"constructor\",\"outputs\":[]}],\"header\":[\"time\"]}"
	ReceiverTvc = "te6ccgECEAEAAikAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCw4FBA8CiiHbPNMAAZ+BAgDXGCD5AVj4QvkQ8qje0z8Bjh34QyG5IJ8wIPgjgQPoqIIIG3dAoLnekyD4Y+DyNNgw0x8B2zz4R27yfAgGATQi0NcLA6k4ANwhxwDcIdMfId0B2zz4R27yfAYDPCCCECynyLK64wIgghBhnwjbuuMCIIIQaLVfP7rjAgoJBwIqMPhBbuMA+Ebyc3H4ZtH4ANs8f/hnCAsAdu1E0CDXScIBjhPT/9M/0wD0Bfhqf/hh+Gb4Y/hijhv0BW34anABgED0DvK91wv/+GJw+GNw+GZ/+GHiAuIw+EFu4wDT/9cN/5XU0dDT/9/R+EohAVMQgQEA9A6T1wv/kXDipMjL/1mBAQD0Q/hqyIvcAAAAAAAAAAAAAAAAIM8Wz4HPgc+ROGILIiLPC/8hzwv/IfhKgQEA9A6T1wv/kXDizwv/yXD7AFvbPH/4Zw0LA3Qw+EFu4wDT/9HbPCHA/44jI9DTAfpAMDHIz4cgzoBgz0DPgc+Bz5KynyLKIc8L/8lw+wDeMOMAf/hnDQwLAC74QsjL//hDzws/+EbPCwD4SgH0AMntVAAkcCH4SoEBAPQOk9cL/5Fw4jExACztRNDT/9M/0wD0Bfhqf/hh+Gb4Y/hiAQr0pCD0oQ8AAA=="
)

type ReceiverContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type Receiver struct {
	Ctx ContractContext
}

func (contract *ReceiverContract) ReceiveData(data string, destinationChainId string) *ContractMethod {
	input := fmt.Sprintf("{\"data\": \"%s\" ,\"destinationChainId\": \"%s\" }", data, destinationChainId)
	return contract.CallContractMethod("receiveData", input)
}
func (contract *ReceiverContract) GetNonceByChainId(destinationChainId string) *ContractMethod {
	input := fmt.Sprintf("{\"destinationChainId\": \"%s\" }", destinationChainId)
	return contract.CallContractMethod("getNonceByChainId", input)
}
func (c *Receiver) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{ReceiverTvc})
}
func (c *Receiver) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(ReceiverAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *Receiver) Address() (string, error) {
	encodeMessage, err := c.DeployEncodeMessage()
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *Receiver) New(address string) (*ReceiverContract, error) {
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
	contract := ReceiverContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *Receiver) DeployEncodeMessage() (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         ReceiverTvc,
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
func (c *Receiver) Deploy(messageCallback func(event *client.ProcessingEvent)) (*ReceiverContract, error) {
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
	contract := ReceiverContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *ReceiverContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *ReceiverContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *ReceiverContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *ReceiverContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
	if len(res.Result) == 0 {
		return &client.DecodedOutput{}, nil
	}
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
