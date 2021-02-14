package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	BurnedTokensHandlerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[{\"inputs\":[{\"name\":\"destinationChainID\",\"type\":\"uint8\"},{\"name\":\"resourceID\",\"type\":\"uint256\"},{\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint128\"},{\"name\":\"recipient\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"outputs\":[]}],\"functions\":[{\"inputs\":[{\"name\":\"_tip3RootAddress\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"payload\",\"type\":\"cell\"},{\"name\":\"sender_public_key\",\"type\":\"uint256\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"wallet_address\",\"type\":\"address\"}],\"name\":\"burnCallback\",\"outputs\":[]},{\"inputs\":[{\"name\":\"destinationChainID\",\"type\":\"uint8\"},{\"name\":\"resourceID\",\"type\":\"uint256\"},{\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"name\":\"amount\",\"type\":\"uint128\"},{\"name\":\"recipient\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[]}],\"header\":[\"time\"]}"
	BurnedTokensHandlerTvc = "te6ccgECEwEAAxYAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAib/APSkICLAAZL0oOGK7VNYMPShCAQBCvSkIPShBQIJngAAAAoHBgApTtRNDT/9M/0wD4an/4Yfhm+GP4YoAC1fhCyMv/+EPPCz/4Rs8LAPhKAc7J7VSAIBIAwJAYj/f40IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIe1E0CDXScIBjhHT/9M/0wD4an/4Yfhm+GP4YgoB/I4+9AWNCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAT4anABgED0DvK91wv/+GJw+GNw+GZ/+GHi0wABn4ECANcYIPkBWPhC+RDyqN7TPwGOHfhDIbkgnzAg+COBA+iogggbd0Cgud6S+GPgMPI02NMfIcEDIgsAMIIQ/////byxk1vyPOAB8AH4R26TMPI83gIBIBINAgEgDw4At7sd0ndPhBbpLwCt7Tf9TT//pBldTR0PpA3/pBldTR0PpA39H4SfhKxwXy4GT4KMjPhQjOjQRQF9eEAAAAAAAAAAAAAAAAAAHPFs+Bz4MkzxTJcPsAXwXwCX/4Z4AgEgERAA5blCYSdaYPp/+mf64a/yupo6Gm/7+uG/8rqaOhp/+/o/CT9ITeJ64X//BR9ITeJ64X/3XlwMmRF7gAAAAAAAAAAAAAAABBni2fA58DnyH50ISMS54WDkmeF/5HnhZ+RZ4W/kOeF/+S4/YAvgsl4BO8//DPAA97jqypufCC3Ry52omgQa6ThAMcI6f/pn+mAfDU//DD8M3wx/DFHH3oCxoQwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACfDU4AMAgegd5XuuF//wxOHwxuHwzP/ww8W98I3k5uPwzfSBo/AAQfDUYeAS//DPAAkN5wItDTA/pAMPhpqTgA+ER/b3GCCJiWgG9ybW9zcW90+GTcIccA3CHTHyHdIcEDIoIQ/////byxk1vyPOAB8AH4R26TMPI83g=="
)

type BurnedTokensHandlerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type BurnedTokensHandler struct {
	Ctx ContractContext
}
type BurnedTokensHandlerDeployParams struct {
	Tip3RootAddress string
}

func (c *BurnedTokensHandler) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{BurnedTokensHandlerTvc})
}
func (c *BurnedTokensHandler) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(BurnedTokensHandlerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *BurnedTokensHandler) Address() (string, error) {
	burnedTokensHandlerDeployParams := BurnedTokensHandlerDeployParams{Tip3RootAddress: "0:0000000000000000000000000000000000000000000000000000000000000000"}
	encodeMessage, err := c.DeployEncodeMessage(&burnedTokensHandlerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *BurnedTokensHandler) New(address string) (*BurnedTokensHandlerContract, error) {
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
	contract := BurnedTokensHandlerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *BurnedTokensHandler) DeployEncodeMessage(burnedTokensHandlerDeployParams *BurnedTokensHandlerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         BurnedTokensHandlerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_tip3RootAddress\": \"%s\" }", burnedTokensHandlerDeployParams.Tip3RootAddress))
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
func (c *BurnedTokensHandler) Deploy(burnedTokensHandlerDeployParams *BurnedTokensHandlerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*BurnedTokensHandlerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(burnedTokensHandlerDeployParams)
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
	contract := BurnedTokensHandlerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *BurnedTokensHandlerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *BurnedTokensHandlerContract) BurnCallback(tokens string, payload string, sender_public_key string, sender_address string, wallet_address string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"payload\": \"%s\" ,\"sender_public_key\": \"%s\" ,\"sender_address\": \"%s\" ,\"wallet_address\": \"%s\" }", tokens, payload, sender_public_key, sender_address, wallet_address)
	return contract.CallContractMethod("burnCallback", input)
}
func (contract *BurnedTokensHandlerContract) Deposit(destinationChainID string, resourceID string, depositNonce string, amount string, recipient string) *ContractMethod {
	input := fmt.Sprintf("{\"destinationChainID\": \"%s\" ,\"resourceID\": \"%s\" ,\"depositNonce\": \"%s\" ,\"amount\": \"%s\" ,\"recipient\": \"%s\" }", destinationChainID, resourceID, depositNonce, amount, recipient)
	return contract.CallContractMethod("deposit", input)
}
func (contract *BurnedTokensHandlerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *BurnedTokensHandlerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *BurnedTokensHandlerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
