package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	RelayerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_accessControllerAddress\",\"type\":\"address\"},{\"name\":\"_myPublicKey\",\"type\":\"uint256\"},{\"name\":\"_myInitState\",\"type\":\"cell\"},{\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"}],\"name\":\"bridgeSetHandler\",\"outputs\":[]},{\"inputs\":[{\"name\":\"epochNumber\",\"type\":\"uint64\"},{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"voteThroughBridge\",\"outputs\":[]},{\"inputs\":[{\"name\":\"epochAddress\",\"type\":\"address\"},{\"name\":\"signHighPart\",\"type\":\"uint256\"},{\"name\":\"signLowPart\",\"type\":\"uint256\"},{\"name\":\"pubkey\",\"type\":\"uint256\"}],\"name\":\"signUpForEpoch\",\"outputs\":[]},{\"inputs\":[{\"name\":\"epochAddress\",\"type\":\"address\"},{\"name\":\"signHighPart\",\"type\":\"uint256\"},{\"name\":\"signLowPart\",\"type\":\"uint256\"},{\"name\":\"pubkey\",\"type\":\"uint256\"}],\"name\":\"forceEra\",\"outputs\":[]},{\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint128\"}],\"name\":\"updateValueForChangeRole\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getValueForChangeRole\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint128\"}],\"name\":\"updateValueForChangeSuperAdmin\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getValueForChangeSuperAdmin\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"grantSuperAdmin\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getRole\",\"outputs\":[{\"name\":\"my_role\",\"type\":\"uint8\"}]},{\"inputs\":[{\"name\":\"role\",\"type\":\"uint8\"},{\"name\":\"targetAddress\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[]},{\"inputs\":[{\"name\":\"initiatorRole\",\"type\":\"uint8\"},{\"name\":\"role\",\"type\":\"uint8\"},{\"name\":\"touchingPublicKey\",\"type\":\"uint256\"}],\"name\":\"changeRole\",\"outputs\":[]},{\"inputs\":[],\"name\":\"deactivateHimself\",\"outputs\":[]}],\"header\":[\"time\"]}"
	RelayerTvc = "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"
)

type RelayerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type Relayer struct {
	Ctx ContractContext
}
type RelayerDeployParams struct {
	AccessControllerAddress string
	MyPublicKey             string
	MyInitState             string
	BridgeAddress           string
}

func (c *Relayer) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{RelayerTvc})
}
func (c *Relayer) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(RelayerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *Relayer) Address() (string, error) {
	relayerDeployParams := RelayerDeployParams{
		AccessControllerAddress: "0:0000000000000000000000000000000000000000000000000000000000000000",
		BridgeAddress:           "0:0000000000000000000000000000000000000000000000000000000000000000",
		MyInitState:             "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		MyPublicKey:             "0x7",
	}
	encodeMessage, err := c.DeployEncodeMessage(&relayerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *Relayer) New(address string) (*RelayerContract, error) {
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
	contract := RelayerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *Relayer) DeployEncodeMessage(relayerDeployParams *RelayerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         RelayerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_accessControllerAddress\": \"%s\" ,\"_myPublicKey\": \"%s\" ,\"_myInitState\": \"%s\" ,\"_bridgeAddress\": \"%s\" }", relayerDeployParams.AccessControllerAddress, relayerDeployParams.MyPublicKey, relayerDeployParams.MyInitState, relayerDeployParams.BridgeAddress))
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
func (c *Relayer) Deploy(relayerDeployParams *RelayerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*RelayerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(relayerDeployParams)
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
	contract := RelayerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *RelayerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *RelayerContract) BridgeSetHandler(messageType string, handlerAddress string) *ContractMethod {
	input := fmt.Sprintf("{\"messageType\": \"%s\" ,\"handlerAddress\": \"%s\" }", messageType, handlerAddress)
	return contract.CallContractMethod("bridgeSetHandler", input)
}
func (contract *RelayerContract) VoteThroughBridge(epochNumber string, choice string, chainId string, messageType string, nonce string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"epochNumber\": \"%s\" ,\"choice\": \"%s\" ,\"chainId\": \"%s\" ,\"messageType\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" }", epochNumber, choice, chainId, messageType, nonce, data)
	return contract.CallContractMethod("voteThroughBridge", input)
}
func (contract *RelayerContract) SignUpForEpoch(epochAddress string, signHighPart string, signLowPart string, pubkey string) *ContractMethod {
	input := fmt.Sprintf("{\"epochAddress\": \"%s\" ,\"signHighPart\": \"%s\" ,\"signLowPart\": \"%s\" ,\"pubkey\": \"%s\" }", epochAddress, signHighPart, signLowPart, pubkey)
	return contract.CallContractMethod("signUpForEpoch", input)
}
func (contract *RelayerContract) ForceEra(epochAddress string, signHighPart string, signLowPart string, pubkey string) *ContractMethod {
	input := fmt.Sprintf("{\"epochAddress\": \"%s\" ,\"signHighPart\": \"%s\" ,\"signLowPart\": \"%s\" ,\"pubkey\": \"%s\" }", epochAddress, signHighPart, signLowPart, pubkey)
	return contract.CallContractMethod("forceEra", input)
}
func (contract *RelayerContract) UpdateValueForChangeRole(newValue string) *ContractMethod {
	input := fmt.Sprintf("{\"newValue\": \"%s\" }", newValue)
	return contract.CallContractMethod("updateValueForChangeRole", input)
}
func (contract *RelayerContract) GetValueForChangeRole() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getValueForChangeRole", input)
}
func (contract *RelayerContract) UpdateValueForChangeSuperAdmin(newValue string) *ContractMethod {
	input := fmt.Sprintf("{\"newValue\": \"%s\" }", newValue)
	return contract.CallContractMethod("updateValueForChangeSuperAdmin", input)
}
func (contract *RelayerContract) GetValueForChangeSuperAdmin() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getValueForChangeSuperAdmin", input)
}
func (contract *RelayerContract) GrantSuperAdmin() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("grantSuperAdmin", input)
}
func (contract *RelayerContract) GetRole() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getRole", input)
}
func (contract *RelayerContract) GrantRole(role string, targetAddress string) *ContractMethod {
	input := fmt.Sprintf("{\"role\": \"%s\" ,\"targetAddress\": \"%s\" }", role, targetAddress)
	return contract.CallContractMethod("grantRole", input)
}
func (contract *RelayerContract) ChangeRole(initiatorRole string, role string, touchingPublicKey string) *ContractMethod {
	input := fmt.Sprintf("{\"initiatorRole\": \"%s\" ,\"role\": \"%s\" ,\"touchingPublicKey\": \"%s\" }", initiatorRole, role, touchingPublicKey)
	return contract.CallContractMethod("changeRole", input)
}
func (contract *RelayerContract) DeactivateHimself() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("deactivateHimself", input)
}
func (contract *RelayerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *RelayerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *RelayerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
