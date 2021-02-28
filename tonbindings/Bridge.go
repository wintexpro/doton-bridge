package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	BridgeAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_relayerInitState\",\"type\":\"cell\"},{\"name\":\"_accessControllerAddress\",\"type\":\"address\"},{\"name\":\"_voteControllerAddress\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"relayerPubKey\",\"type\":\"uint256\"}],\"name\":\"adminSetHandler\",\"outputs\":[]},{\"inputs\":[{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"},{\"name\":\"relayerPubKey\",\"type\":\"uint256\"}],\"name\":\"relayerVoteForProposal\",\"outputs\":[]},{\"inputs\":[{\"name\":\"messageType\",\"type\":\"uint256\"}],\"name\":\"getHandlerAddressByMessageType\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]}],\"header\":[\"time\"]}"
	BridgeTvc = "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"
)

type BridgeContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type Bridge struct {
	Ctx ContractContext
}
type BridgeDeployParams struct {
	RelayerInitState        string
	AccessControllerAddress string
	VoteControllerAddress   string
}

func (c *Bridge) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{BridgeTvc})
}
func (c *Bridge) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(BridgeAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *Bridge) DecodeMessageBody(body string, isInternal bool) (*client.DecodedMessageBody, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}

	params := client.ParamsOfDecodeMessageBody{
		Abi:        *abi,
		Body:       body,
		IsInternal: isInternal,
	}

	return c.Ctx.Conn.AbiDecodeMessageBody(&params)
}
func (c *Bridge) Address() (string, error) {
	bridgeDeployParams := BridgeDeployParams{
		AccessControllerAddress: "0:0000000000000000000000000000000000000000000000000000000000000000",
		RelayerInitState:        "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		VoteControllerAddress:   "0:0000000000000000000000000000000000000000000000000000000000000000",
	}
	encodeMessage, err := c.DeployEncodeMessage(&bridgeDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *Bridge) New(address string) (*BridgeContract, error) {
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
	contract := BridgeContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *Bridge) DeployEncodeMessage(bridgeDeployParams *BridgeDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         BridgeTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_relayerInitState\": \"%s\" ,\"_accessControllerAddress\": \"%s\" ,\"_voteControllerAddress\": \"%s\" }", bridgeDeployParams.RelayerInitState, bridgeDeployParams.AccessControllerAddress, bridgeDeployParams.VoteControllerAddress))
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
func (c *Bridge) Deploy(bridgeDeployParams *BridgeDeployParams, messageCallback func(event *client.ProcessingEvent)) (*BridgeContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(bridgeDeployParams)
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
	contract := BridgeContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *BridgeContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *BridgeContract) AdminSetHandler(messageType string, handlerAddress string, relayerPubKey string) *ContractMethod {
	input := fmt.Sprintf("{\"messageType\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"relayerPubKey\": \"%s\" }", messageType, handlerAddress, relayerPubKey)
	return contract.CallContractMethod("adminSetHandler", input)
}
func (contract *BridgeContract) RelayerVoteForProposal(choice string, chainId string, messageType string, nonce string, data string, relayerPubKey string) *ContractMethod {
	input := fmt.Sprintf("{\"choice\": \"%s\" ,\"chainId\": \"%s\" ,\"messageType\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" ,\"relayerPubKey\": \"%s\" }", choice, chainId, messageType, nonce, data, relayerPubKey)
	return contract.CallContractMethod("relayerVoteForProposal", input)
}
func (contract *BridgeContract) GetHandlerAddressByMessageType(messageType string) *ContractMethod {
	input := fmt.Sprintf("{\"messageType\": \"%s\" }", messageType)
	return contract.CallContractMethod("getHandlerAddressByMessageType", input)
}
func (contract *BridgeContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *BridgeContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *BridgeContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
