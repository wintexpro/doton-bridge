package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	TONTokenWalletAbi = "{\"ABI version\":2,\"data\":[{\"key\":1,\"name\":\"root_address\",\"type\":\"address\"},{\"key\":2,\"name\":\"code\",\"type\":\"cell\"},{\"key\":3,\"name\":\"wallet_public_key\",\"type\":\"uint256\"},{\"key\":4,\"name\":\"owner_address\",\"type\":\"address\"}],\"events\":[],\"functions\":[{\"inputs\":[],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getDetails\",\"outputs\":[{\"components\":[{\"name\":\"root_address\",\"type\":\"address\"},{\"name\":\"code\",\"type\":\"cell\"},{\"name\":\"wallet_public_key\",\"type\":\"uint256\"},{\"name\":\"owner_address\",\"type\":\"address\"},{\"name\":\"balance\",\"type\":\"uint128\"}],\"name\":\"value0\",\"type\":\"tuple\"}]},{\"inputs\":[],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"getWalletKey\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint256\"}]},{\"inputs\":[],\"name\":\"getRootAddress\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]},{\"inputs\":[],\"name\":\"getOwnerAddress\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"}],\"name\":\"accept\",\"outputs\":[]},{\"inputs\":[],\"name\":\"allowance\",\"outputs\":[{\"components\":[{\"name\":\"remaining_tokens\",\"type\":\"uint128\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"value0\",\"type\":\"tuple\"}]},{\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"remaining_tokens\",\"type\":\"uint128\"},{\"name\":\"tokens\",\"type\":\"uint128\"}],\"name\":\"approve\",\"outputs\":[]},{\"inputs\":[],\"name\":\"disapprove\",\"outputs\":[]},{\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"}],\"name\":\"transfer\",\"outputs\":[]},{\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"transferWithNotify\",\"outputs\":[]},{\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"}],\"name\":\"transferFrom\",\"outputs\":[]},{\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"transferFromWithNotify\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"sender_public_key\",\"type\":\"uint256\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"send_gas_to\",\"type\":\"address\"},{\"name\":\"notify_receiver\",\"type\":\"bool\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"internalTransfer\",\"outputs\":[]},{\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"send_gas_to\",\"type\":\"address\"},{\"name\":\"notify_receiver\",\"type\":\"bool\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"internalTransferFrom\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"burnByOwner\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"burnByRoot\",\"outputs\":[]},{\"inputs\":[{\"name\":\"receive_callback_\",\"type\":\"address\"}],\"name\":\"setReceiveCallback\",\"outputs\":[]},{\"inputs\":[{\"name\":\"gas_dest\",\"type\":\"address\"}],\"name\":\"destroy\",\"outputs\":[]},{\"inputs\":[],\"name\":\"balance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"receive_callback\",\"outputs\":[{\"name\":\"receive_callback\",\"type\":\"address\"}]},{\"inputs\":[],\"name\":\"target_gas_balance\",\"outputs\":[{\"name\":\"target_gas_balance\",\"type\":\"uint128\"}]}],\"header\":[\"time\",\"expire\"]}"
	TONTokenWalletTvc = "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"
)

type TONTokenWalletContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type TONTokenWallet struct {
	Ctx ContractContext
}
type TONTokenWalletInitVars struct {
	Rootaddress     string
	Code            string
	Walletpublickey string
	Owneraddress    string
}

func (c *TONTokenWallet) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{TONTokenWalletTvc})
}
func (c *TONTokenWallet) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(TONTokenWalletAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *TONTokenWallet) Address(tONTokenWalletInitVars *TONTokenWalletInitVars) (string, error) {
	encodeMessage, err := c.DeployEncodeMessage(tONTokenWalletInitVars)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *TONTokenWallet) New(address string, tONTokenWalletInitVars *TONTokenWalletInitVars) (*TONTokenWalletContract, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	if address == "" {
		address, err = c.Address(tONTokenWalletInitVars)
		if err != nil {
			return nil, err
		}
	}
	contract := TONTokenWalletContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *TONTokenWallet) DeployEncodeMessage(tONTokenWalletInitVars *TONTokenWalletInitVars) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	initialVars := json.RawMessage(fmt.Sprintf("{\"root_address\": \"%s\" ,\"code\": \"%s\" ,\"wallet_public_key\": \"%s\" ,\"owner_address\": \"%s\" }", tONTokenWalletInitVars.Rootaddress, tONTokenWalletInitVars.Code, tONTokenWalletInitVars.Walletpublickey, tONTokenWalletInitVars.Owneraddress))
	deploySet := client.DeploySet{
		InitialData: initialVars,
		Tvc:         TONTokenWalletTvc,
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
func (c *TONTokenWallet) Deploy(tONTokenWalletInitVars *TONTokenWalletInitVars, messageCallback func(event *client.ProcessingEvent)) (*TONTokenWalletContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(tONTokenWalletInitVars)
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
	contract := TONTokenWalletContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *TONTokenWalletContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *TONTokenWalletContract) GetDetails() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getDetails", input)
}
func (contract *TONTokenWalletContract) GetBalance() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getBalance", input)
}
func (contract *TONTokenWalletContract) GetWalletKey() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getWalletKey", input)
}
func (contract *TONTokenWalletContract) GetRootAddress() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getRootAddress", input)
}
func (contract *TONTokenWalletContract) GetOwnerAddress() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getOwnerAddress", input)
}
func (contract *TONTokenWalletContract) Accept(tokens string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" }", tokens)
	return contract.CallContractMethod("accept", input)
}
func (contract *TONTokenWalletContract) Allowance() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("allowance", input)
}
func (contract *TONTokenWalletContract) Approve(spender string, remaining_tokens string, tokens string) *ContractMethod {
	input := fmt.Sprintf("{\"spender\": \"%s\" ,\"remaining_tokens\": \"%s\" ,\"tokens\": \"%s\" }", spender, remaining_tokens, tokens)
	return contract.CallContractMethod("approve", input)
}
func (contract *TONTokenWalletContract) Disapprove() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("disapprove", input)
}
func (contract *TONTokenWalletContract) Transfer(to string, tokens string, grams string) *ContractMethod {
	input := fmt.Sprintf("{\"to\": \"%s\" ,\"tokens\": \"%s\" ,\"grams\": \"%s\" }", to, tokens, grams)
	return contract.CallContractMethod("transfer", input)
}
func (contract *TONTokenWalletContract) TransferWithNotify(to string, tokens string, grams string, payload string) *ContractMethod {
	input := fmt.Sprintf("{\"to\": \"%s\" ,\"tokens\": \"%s\" ,\"grams\": \"%s\" ,\"payload\": \"%s\" }", to, tokens, grams, payload)
	return contract.CallContractMethod("transferWithNotify", input)
}
func (contract *TONTokenWalletContract) TransferFrom(from string, to string, tokens string, grams string) *ContractMethod {
	input := fmt.Sprintf("{\"from\": \"%s\" ,\"to\": \"%s\" ,\"tokens\": \"%s\" ,\"grams\": \"%s\" }", from, to, tokens, grams)
	return contract.CallContractMethod("transferFrom", input)
}
func (contract *TONTokenWalletContract) TransferFromWithNotify(from string, to string, tokens string, grams string, payload string) *ContractMethod {
	input := fmt.Sprintf("{\"from\": \"%s\" ,\"to\": \"%s\" ,\"tokens\": \"%s\" ,\"grams\": \"%s\" ,\"payload\": \"%s\" }", from, to, tokens, grams, payload)
	return contract.CallContractMethod("transferFromWithNotify", input)
}
func (contract *TONTokenWalletContract) InternalTransfer(tokens string, sender_public_key string, sender_address string, send_gas_to string, notify_receiver string, payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"sender_public_key\": \"%s\" ,\"sender_address\": \"%s\" ,\"send_gas_to\": \"%s\" ,\"notify_receiver\": \"%s\" ,\"payload\": \"%s\" }", tokens, sender_public_key, sender_address, send_gas_to, notify_receiver, payload)
	return contract.CallContractMethod("internalTransfer", input)
}
func (contract *TONTokenWalletContract) InternalTransferFrom(to string, tokens string, send_gas_to string, notify_receiver string, payload string) *ContractMethod {
	input := fmt.Sprintf("{\"to\": \"%s\" ,\"tokens\": \"%s\" ,\"send_gas_to\": \"%s\" ,\"notify_receiver\": \"%s\" ,\"payload\": \"%s\" }", to, tokens, send_gas_to, notify_receiver, payload)
	return contract.CallContractMethod("internalTransferFrom", input)
}
func (contract *TONTokenWalletContract) BurnByOwner(tokens string, grams string, callback_address string, callback_payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"grams\": \"%s\" ,\"callback_address\": \"%s\" ,\"callback_payload\": \"%s\" }", tokens, grams, callback_address, callback_payload)
	return contract.CallContractMethod("burnByOwner", input)
}
func (contract *TONTokenWalletContract) BurnByRoot(tokens string, callback_address string, callback_payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"callback_address\": \"%s\" ,\"callback_payload\": \"%s\" }", tokens, callback_address, callback_payload)
	return contract.CallContractMethod("burnByRoot", input)
}
func (contract *TONTokenWalletContract) SetReceiveCallback(receive_callback_ string) *ContractMethod {
	input := fmt.Sprintf("{\"receive_callback_\": \"%s\" }", receive_callback_)
	return contract.CallContractMethod("setReceiveCallback", input)
}
func (contract *TONTokenWalletContract) Destroy(gas_dest string) *ContractMethod {
	input := fmt.Sprintf("{\"gas_dest\": \"%s\" }", gas_dest)
	return contract.CallContractMethod("destroy", input)
}
func (contract *TONTokenWalletContract) Balance() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("balance", input)
}
func (contract *TONTokenWalletContract) Receive_callback() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("receive_callback", input)
}
func (contract *TONTokenWalletContract) Target_gas_balance() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("target_gas_balance", input)
}
func (contract *TONTokenWalletContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *TONTokenWalletContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *TONTokenWalletContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
