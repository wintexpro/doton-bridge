package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	RelayerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_accessControllerAddress\",\"type\":\"address\"},{\"name\":\"_myPublicKey\",\"type\":\"uint256\"},{\"name\":\"_myInitState\",\"type\":\"cell\"},{\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"}],\"name\":\"bridgeSetHandler\",\"outputs\":[]},{\"inputs\":[{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"voteThroughBridge\",\"outputs\":[]},{\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint128\"}],\"name\":\"updateValueForChangeRole\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getValueForChangeRole\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint128\"}],\"name\":\"updateValueForChangeSuperAdmin\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getValueForChangeSuperAdmin\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"grantSuperAdmin\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getRole\",\"outputs\":[{\"name\":\"my_role\",\"type\":\"uint8\"}]},{\"inputs\":[{\"name\":\"role\",\"type\":\"uint8\"},{\"name\":\"targetAddress\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[]},{\"inputs\":[{\"name\":\"initiatorRole\",\"type\":\"uint8\"},{\"name\":\"role\",\"type\":\"uint8\"},{\"name\":\"touchingPublicKey\",\"type\":\"uint256\"}],\"name\":\"changeRole\",\"outputs\":[]},{\"inputs\":[],\"name\":\"deactivateHimself\",\"outputs\":[]}],\"header\":[\"time\"]}"
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
func (contract *RelayerContract) VoteThroughBridge(choice string, chainId string, messageType string, nonce string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"choice\": \"%s\" ,\"chainId\": \"%s\" ,\"messageType\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" }", choice, chainId, messageType, nonce, data)
	fmt.Print("\n ------------ \n")
	fmt.Print(input)
	fmt.Print("\n ------------ \n")
	return contract.CallContractMethod("voteThroughBridge", input)
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
