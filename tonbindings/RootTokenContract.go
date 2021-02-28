package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	RootTokenContractAbi = "{\"ABI version\":2,\"data\":[{\"key\":1,\"name\":\"_randomNonce\",\"type\":\"uint256\"},{\"key\":2,\"name\":\"name\",\"type\":\"bytes\"},{\"key\":3,\"name\":\"symbol\",\"type\":\"bytes\"},{\"key\":4,\"name\":\"decimals\",\"type\":\"uint8\"},{\"key\":5,\"name\":\"wallet_code\",\"type\":\"cell\"}],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"root_public_key_\",\"type\":\"uint256\"},{\"name\":\"root_owner_address_\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getDetails\",\"outputs\":[{\"components\":[{\"name\":\"name\",\"type\":\"bytes\"},{\"name\":\"symbol\",\"type\":\"bytes\"},{\"name\":\"decimals\",\"type\":\"uint8\"},{\"name\":\"wallet_code\",\"type\":\"cell\"},{\"name\":\"root_public_key\",\"type\":\"uint256\"},{\"name\":\"root_owner_address\",\"type\":\"address\"},{\"name\":\"total_supply\",\"type\":\"uint128\"},{\"name\":\"start_gas_balance\",\"type\":\"uint128\"},{\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"value0\",\"type\":\"tuple\"}]},{\"inputs\":[{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"}],\"name\":\"getWalletAddress\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"},{\"name\":\"gas_back_address\",\"type\":\"address\"}],\"name\":\"deployWallet\",\"outputs\":[]},{\"inputs\":[{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"},{\"name\":\"gas_back_address\",\"type\":\"address\"}],\"name\":\"deployEmptyWallet\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"proxyBurn\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"sender_public_key\",\"type\":\"uint256\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"tokensBurned\",\"outputs\":[]},{\"inputs\":[],\"name\":\"withdrawExtraGas\",\"outputs\":[]},{\"inputs\":[{\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setPaused\",\"outputs\":[]},{\"inputs\":[{\"name\":\"callback_id\",\"type\":\"uint64\"},{\"name\":\"callback_addr\",\"type\":\"address\"}],\"name\":\"sendPausedCallbackTo\",\"outputs\":[]},{\"inputs\":[{\"name\":\"root_public_key_\",\"type\":\"uint256\"},{\"name\":\"root_owner_address_\",\"type\":\"address\"}],\"name\":\"transferOwner\",\"outputs\":[]},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"name\",\"type\":\"bytes\"}]},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"symbol\",\"type\":\"bytes\"}]},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\"}]},{\"inputs\":[],\"name\":\"wallet_code\",\"outputs\":[{\"name\":\"wallet_code\",\"type\":\"cell\"}]},{\"inputs\":[],\"name\":\"total_supply\",\"outputs\":[{\"name\":\"total_supply\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"paused\",\"type\":\"bool\"}]}],\"header\":[\"pubkey\",\"time\",\"expire\"]}"
	RootTokenContractTvc = "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"
)

type RootTokenContractContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type RootTokenContract struct {
	Ctx ContractContext
}
type RootTokenContractDeployParams struct {
	Rootpublickey    string
	Rootowneraddress string
}
type RootTokenContractInitVars struct {
	RandomNonce string
	Name        string
	Symbol      string
	Decimals    string
	Walletcode  string
}

func (c *RootTokenContract) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{RootTokenContractTvc})
}
func (c *RootTokenContract) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(RootTokenContractAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *RootTokenContract) Address(rootTokenContractInitVars *RootTokenContractInitVars) (string, error) {
	rootTokenContractDeployParams := RootTokenContractDeployParams{
		Rootowneraddress: "0:0000000000000000000000000000000000000000000000000000000000000000",
		Rootpublickey:    "0x7",
	}
	encodeMessage, err := c.DeployEncodeMessage(&rootTokenContractDeployParams, rootTokenContractInitVars)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *RootTokenContract) DecodeMessageBody(body string, isInternal bool) (*client.DecodedMessageBody, error) {
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
func (c *RootTokenContract) New(address string, rootTokenContractInitVars *RootTokenContractInitVars) (*RootTokenContractContract, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	if address == "" {
		address, err = c.Address(rootTokenContractInitVars)
		if err != nil {
			return nil, err
		}
	}
	contract := RootTokenContractContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *RootTokenContract) DeployEncodeMessage(rootTokenContractDeployParams *RootTokenContractDeployParams, rootTokenContractInitVars *RootTokenContractInitVars) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	initialVars := json.RawMessage(fmt.Sprintf("{\"_randomNonce\": \"%s\" ,\"name\": \"%s\" ,\"symbol\": \"%s\" ,\"decimals\": \"%s\" ,\"wallet_code\": \"%s\" }", rootTokenContractInitVars.RandomNonce, rootTokenContractInitVars.Name, rootTokenContractInitVars.Symbol, rootTokenContractInitVars.Decimals, rootTokenContractInitVars.Walletcode))
	deploySet := client.DeploySet{
		InitialData: initialVars,
		Tvc:         RootTokenContractTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"root_public_key_\": \"%s\" ,\"root_owner_address_\": \"%s\" }", rootTokenContractDeployParams.Rootpublickey, rootTokenContractDeployParams.Rootowneraddress))
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
func (c *RootTokenContract) Deploy(rootTokenContractDeployParams *RootTokenContractDeployParams, rootTokenContractInitVars *RootTokenContractInitVars, messageCallback func(event *client.ProcessingEvent)) (*RootTokenContractContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(rootTokenContractDeployParams, rootTokenContractInitVars)
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
	contract := RootTokenContractContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *RootTokenContractContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *RootTokenContractContract) GetDetails() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getDetails", input)
}
func (contract *RootTokenContractContract) GetWalletAddress(wallet_public_key_ string, owner_address_ string) *ContractMethod {
	input := fmt.Sprintf("{\"wallet_public_key_\": \"%s\" ,\"owner_address_\": \"%s\" }", wallet_public_key_, owner_address_)
	return contract.CallContractMethod("getWalletAddress", input)
}
func (contract *RootTokenContractContract) DeployWallet(tokens string, grams string, wallet_public_key_ string, owner_address_ string, gas_back_address string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"grams\": \"%s\" ,\"wallet_public_key_\": \"%s\" ,\"owner_address_\": \"%s\" ,\"gas_back_address\": \"%s\" }", tokens, grams, wallet_public_key_, owner_address_, gas_back_address)
	return contract.CallContractMethod("deployWallet", input)
}
func (contract *RootTokenContractContract) DeployEmptyWallet(grams string, wallet_public_key_ string, owner_address_ string, gas_back_address string) *ContractMethod {
	input := fmt.Sprintf("{\"grams\": \"%s\" ,\"wallet_public_key_\": \"%s\" ,\"owner_address_\": \"%s\" ,\"gas_back_address\": \"%s\" }", grams, wallet_public_key_, owner_address_, gas_back_address)
	return contract.CallContractMethod("deployEmptyWallet", input)
}
func (contract *RootTokenContractContract) Mint(tokens string, to string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"to\": \"%s\" }", tokens, to)
	return contract.CallContractMethod("mint", input)
}
func (contract *RootTokenContractContract) ProxyBurn(tokens string, sender_address string, callback_address string, callback_payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"sender_address\": \"%s\" ,\"callback_address\": \"%s\" ,\"callback_payload\": \"%s\" }", tokens, sender_address, callback_address, callback_payload)
	return contract.CallContractMethod("proxyBurn", input)
}
func (contract *RootTokenContractContract) TokensBurned(tokens string, sender_public_key string, sender_address string, callback_address string, callback_payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"sender_public_key\": \"%s\" ,\"sender_address\": \"%s\" ,\"callback_address\": \"%s\" ,\"callback_payload\": \"%s\" }", tokens, sender_public_key, sender_address, callback_address, callback_payload)
	return contract.CallContractMethod("tokensBurned", input)
}
func (contract *RootTokenContractContract) WithdrawExtraGas() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("withdrawExtraGas", input)
}
func (contract *RootTokenContractContract) SetPaused(value string) *ContractMethod {
	input := fmt.Sprintf("{\"value\": \"%s\" }", value)
	return contract.CallContractMethod("setPaused", input)
}
func (contract *RootTokenContractContract) SendPausedCallbackTo(callback_id string, callback_addr string) *ContractMethod {
	input := fmt.Sprintf("{\"callback_id\": \"%s\" ,\"callback_addr\": \"%s\" }", callback_id, callback_addr)
	return contract.CallContractMethod("sendPausedCallbackTo", input)
}
func (contract *RootTokenContractContract) TransferOwner(root_public_key_ string, root_owner_address_ string) *ContractMethod {
	input := fmt.Sprintf("{\"root_public_key_\": \"%s\" ,\"root_owner_address_\": \"%s\" }", root_public_key_, root_owner_address_)
	return contract.CallContractMethod("transferOwner", input)
}
func (contract *RootTokenContractContract) Name() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("name", input)
}
func (contract *RootTokenContractContract) Symbol() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("symbol", input)
}
func (contract *RootTokenContractContract) Decimals() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("decimals", input)
}
func (contract *RootTokenContractContract) Wallet_code() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("wallet_code", input)
}
func (contract *RootTokenContractContract) Total_supply() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("total_supply", input)
}
func (contract *RootTokenContractContract) Paused() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("paused", input)
}
func (contract *RootTokenContractContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *RootTokenContractContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *RootTokenContractContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
