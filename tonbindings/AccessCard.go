package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	AccessCardAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_accessControllerAddress\",\"type\":\"address\"},{\"name\":\"_myPublicKey\",\"type\":\"uint256\"},{\"name\":\"_myInitState\",\"type\":\"cell\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint128\"}],\"name\":\"updateValueForChangeRole\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getValueForChangeRole\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"newValue\",\"type\":\"uint128\"}],\"name\":\"updateValueForChangeSuperAdmin\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getValueForChangeSuperAdmin\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"grantSuperAdmin\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getRole\",\"outputs\":[{\"name\":\"my_role\",\"type\":\"uint8\"}]},{\"inputs\":[{\"name\":\"role\",\"type\":\"uint8\"},{\"name\":\"targetAddress\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[]},{\"inputs\":[{\"name\":\"initiatorRole\",\"type\":\"uint8\"},{\"name\":\"role\",\"type\":\"uint8\"},{\"name\":\"touchingPublicKey\",\"type\":\"uint256\"}],\"name\":\"changeRole\",\"outputs\":[]},{\"inputs\":[],\"name\":\"deactivateHimself\",\"outputs\":[]}],\"header\":[\"time\"]}"
	AccessCardTvc = "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"
)

type AccessCardContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type AccessCard struct {
	Ctx ContractContext
}
type AccessCardDeployParams struct {
	AccessControllerAddress string
	MyPublicKey             string
	MyInitState             string
}

func (c *AccessCard) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{AccessCardTvc})
}
func (c *AccessCard) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(AccessCardAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *AccessCard) Address() (string, error) {
	accessCardDeployParams := AccessCardDeployParams{
		AccessControllerAddress: "0:0000000000000000000000000000000000000000000000000000000000000000",
		MyInitState:             "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		MyPublicKey:             "0x7",
	}
	encodeMessage, err := c.DeployEncodeMessage(&accessCardDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *AccessCard) New(address string) (*AccessCardContract, error) {
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
	contract := AccessCardContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *AccessCard) DeployEncodeMessage(accessCardDeployParams *AccessCardDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         AccessCardTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_accessControllerAddress\": \"%s\" ,\"_myPublicKey\": \"%s\" ,\"_myInitState\": \"%s\" }", accessCardDeployParams.AccessControllerAddress, accessCardDeployParams.MyPublicKey, accessCardDeployParams.MyInitState))
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
func (c *AccessCard) Deploy(accessCardDeployParams *AccessCardDeployParams, messageCallback func(event *client.ProcessingEvent)) (*AccessCardContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(accessCardDeployParams)
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
	contract := AccessCardContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *AccessCardContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *AccessCardContract) UpdateValueForChangeRole(newValue string) *ContractMethod {
	input := fmt.Sprintf("{\"newValue\": \"%s\" }", newValue)
	return contract.CallContractMethod("updateValueForChangeRole", input)
}
func (contract *AccessCardContract) GetValueForChangeRole() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getValueForChangeRole", input)
}
func (contract *AccessCardContract) UpdateValueForChangeSuperAdmin(newValue string) *ContractMethod {
	input := fmt.Sprintf("{\"newValue\": \"%s\" }", newValue)
	return contract.CallContractMethod("updateValueForChangeSuperAdmin", input)
}
func (contract *AccessCardContract) GetValueForChangeSuperAdmin() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getValueForChangeSuperAdmin", input)
}
func (contract *AccessCardContract) GrantSuperAdmin() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("grantSuperAdmin", input)
}
func (contract *AccessCardContract) GetRole() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getRole", input)
}
func (contract *AccessCardContract) GrantRole(role string, targetAddress string) *ContractMethod {
	input := fmt.Sprintf("{\"role\": \"%s\" ,\"targetAddress\": \"%s\" }", role, targetAddress)
	return contract.CallContractMethod("grantRole", input)
}
func (contract *AccessCardContract) ChangeRole(initiatorRole string, role string, touchingPublicKey string) *ContractMethod {
	input := fmt.Sprintf("{\"initiatorRole\": \"%s\" ,\"role\": \"%s\" ,\"touchingPublicKey\": \"%s\" }", initiatorRole, role, touchingPublicKey)
	return contract.CallContractMethod("changeRole", input)
}
func (contract *AccessCardContract) DeactivateHimself() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("deactivateHimself", input)
}
func (contract *AccessCardContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *AccessCardContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *AccessCardContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
