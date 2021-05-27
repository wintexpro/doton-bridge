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
