package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	RootTokenContractAbi = "{\"ABI version\":2,\"data\":[{\"key\":1,\"name\":\"_randomNonce\",\"type\":\"uint256\"},{\"key\":2,\"name\":\"name\",\"type\":\"bytes\"},{\"key\":3,\"name\":\"symbol\",\"type\":\"bytes\"},{\"key\":4,\"name\":\"decimals\",\"type\":\"uint8\"},{\"key\":5,\"name\":\"wallet_code\",\"type\":\"cell\"},{\"key\":6,\"name\":\"root_public_key\",\"type\":\"uint256\"},{\"key\":7,\"name\":\"root_owner_address\",\"type\":\"address\"}],\"events\":[],\"functions\":[{\"inputs\":[],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getDetails\",\"outputs\":[{\"components\":[{\"name\":\"name\",\"type\":\"bytes\"},{\"name\":\"symbol\",\"type\":\"bytes\"},{\"name\":\"decimals\",\"type\":\"uint8\"},{\"name\":\"wallet_code\",\"type\":\"cell\"},{\"name\":\"root_public_key\",\"type\":\"uint256\"},{\"name\":\"root_owner_address\",\"type\":\"address\"},{\"name\":\"total_supply\",\"type\":\"uint128\"},{\"name\":\"start_gas_balance\",\"type\":\"uint128\"}],\"name\":\"value0\",\"type\":\"tuple\"}]},{\"inputs\":[],\"name\":\"getName\",\"outputs\":[{\"name\":\"value0\",\"type\":\"bytes\"}]},{\"inputs\":[],\"name\":\"getSymbol\",\"outputs\":[{\"name\":\"value0\",\"type\":\"bytes\"}]},{\"inputs\":[],\"name\":\"getDecimals\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint8\"}]},{\"inputs\":[],\"name\":\"getRootPublicKey\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint256\"}]},{\"inputs\":[],\"name\":\"getTotalSupply\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"getTotalGranted\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"}],\"name\":\"getWalletAddress\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"},{\"name\":\"gas_back_address\",\"type\":\"address\"}],\"name\":\"deployWallet\",\"outputs\":[]},{\"inputs\":[{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"},{\"name\":\"gas_back_address\",\"type\":\"address\"}],\"name\":\"deployEmptyWallet\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"}],\"name\":\"mint\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"to\",\"type\":\"address\"}],\"name\":\"grant\",\"outputs\":[]},{\"inputs\":[],\"name\":\"forbidRootControl\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"burnTokensOnWallet\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"sender_public_key\",\"type\":\"uint256\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"tokensBurned\",\"outputs\":[]},{\"inputs\":[],\"name\":\"withdrawExtraGas\",\"outputs\":[]},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"name\",\"type\":\"bytes\"}]},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"symbol\",\"type\":\"bytes\"}]},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\"}]},{\"inputs\":[],\"name\":\"wallet_code\",\"outputs\":[{\"name\":\"wallet_code\",\"type\":\"cell\"}]},{\"inputs\":[],\"name\":\"total_supply\",\"outputs\":[{\"name\":\"total_supply\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"total_granted\",\"outputs\":[{\"name\":\"total_granted\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"start_gas_balance\",\"outputs\":[{\"name\":\"start_gas_balance\",\"type\":\"uint128\"}]}],\"header\":[\"time\",\"expire\"]}"
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
type RootTokenContractInitVars struct {
	RandomNonce      string
	Name             string
	Symbol           string
	Decimals         string
	Walletcode       string
	Rootpublickey    string
	Rootowneraddress string
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
	encodeMessage, err := c.DeployEncodeMessage(rootTokenContractInitVars)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
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
func (c *RootTokenContract) DeployEncodeMessage(rootTokenContractInitVars *RootTokenContractInitVars) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	initialVars := json.RawMessage(fmt.Sprintf("{\"_randomNonce\": \"%s\" ,\"name\": \"%s\" ,\"symbol\": \"%s\" ,\"decimals\": \"%s\" ,\"wallet_code\": \"%s\" ,\"root_public_key\": \"%s\" ,\"root_owner_address\": \"%s\" }", rootTokenContractInitVars.RandomNonce, rootTokenContractInitVars.Name, rootTokenContractInitVars.Symbol, rootTokenContractInitVars.Decimals, rootTokenContractInitVars.Walletcode, rootTokenContractInitVars.Rootpublickey, rootTokenContractInitVars.Rootowneraddress))
	deploySet := client.DeploySet{
		InitialData: initialVars,
		Tvc:         RootTokenContractTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{}"))
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
func (c *RootTokenContract) Deploy(rootTokenContractInitVars *RootTokenContractInitVars, messageCallback func(event *client.ProcessingEvent)) (*RootTokenContractContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(rootTokenContractInitVars)
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
func (contract *RootTokenContractContract) GetName() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getName", input)
}
func (contract *RootTokenContractContract) GetSymbol() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getSymbol", input)
}
func (contract *RootTokenContractContract) GetDecimals() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getDecimals", input)
}
func (contract *RootTokenContractContract) GetRootPublicKey() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getRootPublicKey", input)
}
func (contract *RootTokenContractContract) GetTotalSupply() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getTotalSupply", input)
}
func (contract *RootTokenContractContract) GetTotalGranted() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getTotalGranted", input)
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
func (contract *RootTokenContractContract) Mint(tokens string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" }", tokens)
	return contract.CallContractMethod("mint", input)
}
func (contract *RootTokenContractContract) Grant(tokens string, to string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"to\": \"%s\" }", tokens, to)
	return contract.CallContractMethod("grant", input)
}
func (contract *RootTokenContractContract) ForbidRootControl() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("forbidRootControl", input)
}
func (contract *RootTokenContractContract) BurnTokensOnWallet(tokens string, sender_address string, callback_address string, callback_payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"sender_address\": \"%s\" ,\"callback_address\": \"%s\" ,\"callback_payload\": \"%s\" }", tokens, sender_address, callback_address, callback_payload)
	return contract.CallContractMethod("burnTokensOnWallet", input)
}
func (contract *RootTokenContractContract) TokensBurned(tokens string, sender_public_key string, sender_address string, callback_address string, callback_payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"sender_public_key\": \"%s\" ,\"sender_address\": \"%s\" ,\"callback_address\": \"%s\" ,\"callback_payload\": \"%s\" }", tokens, sender_public_key, sender_address, callback_address, callback_payload)
	return contract.CallContractMethod("tokensBurned", input)
}
func (contract *RootTokenContractContract) WithdrawExtraGas() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("withdrawExtraGas", input)
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
func (contract *RootTokenContractContract) Total_granted() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("total_granted", input)
}
func (contract *RootTokenContractContract) Start_gas_balance() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("start_gas_balance", input)
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
