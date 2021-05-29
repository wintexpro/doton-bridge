package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	AccessControllerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_accessCardInitState\",\"type\":\"cell\"},{\"name\":\"_initialValue\",\"type\":\"uint128\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"newInitialValue\",\"type\":\"uint128\"}],\"name\":\"updateInitialValue\",\"outputs\":[]},{\"inputs\":[{\"name\":\"accessCardAddress\",\"type\":\"address\"}],\"name\":\"grantSuperAdminRole\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getSuperAdminAddress\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]},{\"inputs\":[{\"name\":\"newSuperAdminAddress\",\"type\":\"address\"},{\"name\":\"touchingPublicKey\",\"type\":\"uint256\"}],\"name\":\"changeSuperAdmin\",\"outputs\":[]}],\"header\":[\"time\"]}"
	AccessControllerTvc = "te6ccgECGwEABIAAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCxkFBBoC3o0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIds80wABjhKBAgDXGCD5AVj4QiD4ZfkQ8qje0z8Bjh34QyG5IJ8wIPgjgQPoqIIIG3dAoLnekyD4Y+DyNNgw0x8B2zz4R27yfAwGATwi0NMD+kAw+GmpOADcIccA3CHTHyHdAds8+Edu8nwGA0IgghAXvyPPu46A4CCCECwjOWu7joDgIIIQcgesobuOgOASDgcCKCCCEGJFj4i64wIgghByB6yhuuMCCwgDejD4QW7jAPpA1w3/ldTR0NP/39Eg+Eoh2zz4TyDBApMwgGTe+En6Qm8T1wv/IvkAuvL0I/hsMDBb2zx/+GcYCRcB7iHQyCHTADPAAJNxz0Cacc9BIdMfM88LH+Ih0wAzwACTcc9AmnHPQSHTATPPCwHiIdMAM8AAk3HPQJhxz0Eh1DPPFOIh0wAzwwHyd3HPQcgjzwv/ItQ00PQEASJwIoBA9EMxIMj0ACDJJcw1XwQh0wAzwACTcc9ACgAcmHHPQSHUM88U4iDJbEECRjD4QW7jAPhG8nNx+GbU03/R+AAh+Gog+Gv4KPhsW9s8f/hnDBcBaO1E0CDXScIBjifT/9M/0wDU03/6QNMH0wfXCwf4b/hu+G34bPhr+Gp/+GH4Zvhj+GKOgOINAbL0BYj4anD4a40IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhscPhtcPhucPhvcAGAQPQO8r3XC//4YnD4Y3D4Zn/4YYBl+G2AZvhugG/4bxoCKCCCEClDTDu64wIgghAsIzlruuMCEQ8E0DD4QW7jAPpA0Yj4TiDBApMwgGTe+EL4RSBukjBw3rry9Yj4TSDBApMwgGTe+Ez4KMcF8vX4ACD4bCDIz4WIzo0EDmJaAAAAAAAAAAAAAAAAAAHPFs+Bz4HPkcQnC/rJcPsAMNs8f/hnGBUQFwBCU3VwZXJhZG1pbiBhbHJlYWR5IGNyZWF0ZWQgZWFybGVyAm4w+EFu4wDR+EwhwP+OIiPQ0wH6QDAxyM+HIM6AYM9Az4HPgc+SpQ0w7iHPFslw+wDeMOMAf/hnGBcCKCCCEA2QYcm64wIgghAXvyPPuuMCFhMDIDD4QW7jANN/0ds82zx/+GcYFBcBPIj4TiDBApMwgGTe+EL4RSBukjBw3rry9fgAIPhrMBUAHk9ubHkgZm9yIG93bmVycwJwMPhBbuMA0fhLIcD/jiMj0NMB+kAwMcjPhyDOgGDPQM+Bz4HPkjZBhyYhzwt/yXD7AN4w4wB/+GcYFwBU+ELIy//4Q88LP/hGzwsA+Er4S/hM+E34TvhPXlDMy3/OywfLB8sHye1UAFTtRNDT/9M/0wDU03/6QNMH0wfXCwf4b/hu+G34bPhr+Gp/+GH4Zvhj+GIBCvSkIPShGgAA"
)

type AccessControllerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type AccessController struct {
	Ctx ContractContext
}
type AccessControllerDeployParams struct {
	AccessCardInitState string
	InitialValue        string
}

func (c *AccessController) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{AccessControllerTvc})
}
func (c *AccessController) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(AccessControllerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *AccessController) Address() (string, error) {
	accessControllerDeployParams := AccessControllerDeployParams{
		AccessCardInitState: "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		InitialValue:        "0x6",
	}
	encodeMessage, err := c.DeployEncodeMessage(&accessControllerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *AccessController) New(address string) (*AccessControllerContract, error) {
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
	contract := AccessControllerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *AccessController) DeployEncodeMessage(accessControllerDeployParams *AccessControllerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         AccessControllerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_accessCardInitState\": \"%s\" ,\"_initialValue\": \"%s\" }", accessControllerDeployParams.AccessCardInitState, accessControllerDeployParams.InitialValue))
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
func (c *AccessController) Deploy(accessControllerDeployParams *AccessControllerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*AccessControllerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(accessControllerDeployParams)
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
	contract := AccessControllerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *AccessControllerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *AccessControllerContract) GetInitialValue() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getInitialValue", input)
}
func (contract *AccessControllerContract) UpdateInitialValue(newInitialValue string) *ContractMethod {
	input := fmt.Sprintf("{\"newInitialValue\": \"%s\" }", newInitialValue)
	return contract.CallContractMethod("updateInitialValue", input)
}
func (contract *AccessControllerContract) GrantSuperAdminRole(accessCardAddress string) *ContractMethod {
	input := fmt.Sprintf("{\"accessCardAddress\": \"%s\" }", accessCardAddress)
	return contract.CallContractMethod("grantSuperAdminRole", input)
}
func (contract *AccessControllerContract) GetSuperAdminAddress() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getSuperAdminAddress", input)
}
func (contract *AccessControllerContract) ChangeSuperAdmin(newSuperAdminAddress string, touchingPublicKey string) *ContractMethod {
	input := fmt.Sprintf("{\"newSuperAdminAddress\": \"%s\" ,\"touchingPublicKey\": \"%s\" }", newSuperAdminAddress, touchingPublicKey)
	return contract.CallContractMethod("changeSuperAdmin", input)
}
func (contract *AccessControllerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *AccessControllerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *AccessControllerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
