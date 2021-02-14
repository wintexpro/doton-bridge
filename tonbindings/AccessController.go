package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	AccessControllerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_accessCardInitState\",\"type\":\"cell\"},{\"name\":\"_initialValue\",\"type\":\"uint128\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"newInitialValue\",\"type\":\"uint128\"}],\"name\":\"updateInitialValue\",\"outputs\":[]},{\"inputs\":[{\"name\":\"accessCardAddress\",\"type\":\"address\"}],\"name\":\"grantSuperAdminRole\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getSuperAdminAddress\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]},{\"inputs\":[{\"name\":\"newSuperAdminAddress\",\"type\":\"address\"},{\"name\":\"touchingPublicKey\",\"type\":\"uint256\"}],\"name\":\"changeSuperAdmin\",\"outputs\":[]}],\"header\":[\"time\"]}"
	AccessControllerTvc = "te6ccgECIgEABUcAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAib/APSkICLAAZL0oOGK7VNYMPShCAQBCvSkIPShBQIJnwAAAAsHBgA3O1E0NP/0z/TANTTf/hs+Gv4an/4Yfhm+GP4YoAA9PhCyMv/+EPPCz/4Rs8LAPhK+Ev4TF4gzMt/zsntVIAIBIAwJAZb/f40IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIe1E0CDXScIBjhjT/9M/0wDU03/4bPhr+Gp/+GH4Zvhj+GIKAcaORfQFyMn4anD4a40IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhscAGAQPQO8r3XC//4YnD4Y3D4Zn/4YeLTAAGOEoECANcYIPkBWPhCIPhl+RDyqN7TPwELAHyOHfhDIbkgnzAg+COBA+iogggbd0Cgud6S+GPgMPI02NMfIcEDIoIQ/////byxk1vyPOAB8AH4R26TMPI83gIBIBQNAgFYEQ4BCbhA9ZQwDwH++EFukvAL3vpA1w3/ldTR0NP/39Eg+EohIdDIIdMAM8AAk3HPQJpxz0Eh0x8zzwsf4iHTADPAAJNxz0Cacc9BIdMBM88LAeIh0wAzwACTcc9AmHHPQSHUM88U4iHTADPDAfJ3cc9ByCPPC/8i1DTQ9AQBInAigED0QzEgyPQAIBAAeMklzDUl0wA3wACVJHHPQDWbJHHPQTUl1DclzDXiJMkIXwj4SfpCbxPXC/8h+QC68uBvI/hsXwTwCn/4ZwEJuEix8RASAfz4QW6Oau1E0CDXScIBjhjT/9M/0wDU03/4bPhr+Gp/+GH4Zvhj+GKORfQFyMn4anD4a40IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhscAGAQPQO8r3XC//4YnD4Y3D4Zn/4YeLe+Ebyc3H4ZtTTf9H4ACETAB74aiD4a/go+Gxb8Ap/+GcCASAbFQIBZhkWAQm0EZy1wBcB/PhBbpLwC976QNGL9Pbmx5IGZvciBvd25lcnOMjOyfhC+EUgbpIwcN668uhmjQhU3VwZXJhZG1pbiBhbHJlYWR5IGNyZWF0ZWQgZWFybGVygyM7J+Ez4KMcF8uhl+AAg+GwgyM+FiM6NBA5iWgAAAAAAAAAAAAAAAAABzxbPgRgAJM+Bz5HEJwv6yXH7ADDwCn/4ZwEJtKGmHcAaAPz4QW6S8AveIZnTH/hEWG91+GTf0fhMIcD/jiIj0NMB+kAwMcjPhyDOgGDPQM+Bz4HPkqUNMO4hzxbJcfsAjjb4RCBvEyFvEvhJVQJvEchyz0DKAHPPQM4B+gL0AIBoz0DPgc+B+ERvFc8LHyHPFsn4RG8U+wDiMJLwCt5/+GcCASAdHAB3uPfkef8ILdJeAXvab/oxfp7c2PJAzN7kQN7u3Mrk5xkZ2T8IXwikDdJGDhvXXl0M3wAEHw1mHgFP/wzwAgEgIR4BCbdkGHJgHwH8+EFukvAL3iGZ0x/4RFhvdfhk39H4SyHA/44jI9DTAfpAMDHIz4cgzoBgz0DPgc+Bz5I2QYcmIc8Lf8lx+wCON/hEIG8TIW8S+ElVAm8RyHLPQMoAc89AzgH6AvQAgGjPQM+Bz4H4RG8VzwsfIc8Lf8n4RG8U+wDiMJLwCt5/IAAE+GcAkNtwItDTA/pAMPhpqTgA+ER/b3GCCJiWgG9ybW9zcW90+GTcIccA3CHTHyHdIcEDIoIQ/////byxk1vyPOAB8AH4R26TMPI83g=="
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
