package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	MessageHandlerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"ProposalExecuted\",\"outputs\":[]}],\"functions\":[{\"inputs\":[{\"name\":\"_proposalCode\",\"type\":\"cell\"},{\"name\":\"_bridgeVoteControllerAddress\",\"type\":\"address\"},{\"name\":\"_bridgeVoteControllerPubKey\",\"type\":\"uint256\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"proposalPubKey\",\"type\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"executeProposal\",\"outputs\":[]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"receiveMessage\",\"outputs\":[]}],\"header\":[\"time\"]}"
	MessageHandlerTvc = "te6ccgECFgEAA6MAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAib/APSkICLAAZL0oOGK7VNYMPShCAQBCvSkIPShBQIJngAAAAoHBgA9TtRNDT/9M/0wDU+kDXC//4bPhr+Gp/+GH4Zvhj+GKAA9X4QsjL//hDzws/+EbPCwD4SvhL+ExeIMzOy//J7VSAIBIAwJAZz/f40IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIe1E0CDXScIBjhvT/9M/0wDU+kDXC//4bPhr+Gp/+GH4Zvhj+GIKAf6ORfQFyMn4ao0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhrcPhscAGAQPQO8r3XC//4YnD4Y3D4Zn/4YeLTAAGfgQIA1xgg+QFY+EL5EPKo3tM/AY4d+EMhuSCfMCD4I4ED6KiCCBt3QKC53pL4Y+Aw8jTYCwA80x8hwQMighD////9vLGTW/I84AHwAfhHbpMw8jzeAgEgEA0BCbw2SA8cDgH6+EFujm3tRNAg10nCAY4b0//TP9MA1PpA1wv/+Gz4a/hqf/hh+Gb4Y/hijkX0BcjJ+GqNCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAT4a3D4bHABgED0DvK91wv/+GJw+GNw+GZ/+GHi3vhG8nNx+GbU+kAPAETXDf+V1NHQ0//f0SIiIvgAIvhqIfhrIPhsXwNfA/AJf/hnAgEgExEB8bqELXGPhBbpLwCt7T/9MH0z/XDf+V1NHQ0//f1NEkJCQjbfhMyMv/cFiAQPRDI8jLB3FYgED0QyLIyz9yWIBA9EP4S3NYgED0FiF0WIBA9BfI9ADJ+ErIz4SA9AD0AM+ByfhJ+kJvE9cL/yH5ALry4GX4KMjPhQjOgSAFSNBFAX14QAAAAAAAAAAAAAAAAAAc8Wz4HPgybPFMlw+wBfCpLwCd5/+GcCAUgVFADJt9kZIfTB9M/0//XDf+V1NHQ0//f0fhJ+kJvE9cL//go+kJvE9cL/7ry4GTIi9wAAAAAAAAAAAAAAAAgzxbPgc+Bz5CJijO6JM8LByPPCz8izwv/Ic8L/8lx+wBfBJLwCd5/+GeAAkNtwItDTA/pAMPhpqTgA+ER/b3GCCJiWgG9ybW9zcW90+GTcIccA3CHTHyHdIcEDIoIQ/////byxk1vyPOAB8AH4R26TMPI83g=="
)

type MessageHandlerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type MessageHandler struct {
	Ctx ContractContext
}
type MessageHandlerDeployParams struct {
	ProposalCode                string
	BridgeVoteControllerAddress string
	BridgeVoteControllerPubKey  string
}

func (c *MessageHandler) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{MessageHandlerTvc})
}
func (c *MessageHandler) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(MessageHandlerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *MessageHandler) Address() (string, error) {
	messageHandlerDeployParams := MessageHandlerDeployParams{
		BridgeVoteControllerAddress: "0:0000000000000000000000000000000000000000000000000000000000000000",
		BridgeVoteControllerPubKey:  "0x7",
		ProposalCode:                "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
	}
	encodeMessage, err := c.DeployEncodeMessage(&messageHandlerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *MessageHandler) New(address string) (*MessageHandlerContract, error) {
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
	contract := MessageHandlerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *MessageHandler) DeployEncodeMessage(messageHandlerDeployParams *MessageHandlerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         MessageHandlerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_proposalCode\": \"%s\" ,\"_bridgeVoteControllerAddress\": \"%s\" ,\"_bridgeVoteControllerPubKey\": \"%s\" }", messageHandlerDeployParams.ProposalCode, messageHandlerDeployParams.BridgeVoteControllerAddress, messageHandlerDeployParams.BridgeVoteControllerPubKey))
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
func (c *MessageHandler) Deploy(messageHandlerDeployParams *MessageHandlerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*MessageHandlerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(messageHandlerDeployParams)
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
	contract := MessageHandlerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *MessageHandlerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *MessageHandlerContract) ExecuteProposal(proposalPubKey string, chainId string, nonce string, messageType string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"proposalPubKey\": \"%s\" ,\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"messageType\": \"%s\" ,\"data\": \"%s\" }", proposalPubKey, chainId, nonce, messageType, data)
	return contract.CallContractMethod("executeProposal", input)
}
func (contract *MessageHandlerContract) ReceiveMessage(chainId string, nonce string, messageType string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"messageType\": \"%s\" ,\"data\": \"%s\" }", chainId, nonce, messageType, data)
	return contract.CallContractMethod("receiveMessage", input)
}
func (contract *MessageHandlerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *MessageHandlerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *MessageHandlerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
