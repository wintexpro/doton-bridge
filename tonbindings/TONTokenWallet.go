package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	TONTokenWalletAbi = "{\"ABI version\":2,\"data\":[{\"key\":1,\"name\":\"root_address\",\"type\":\"address\"},{\"key\":2,\"name\":\"code\",\"type\":\"cell\"},{\"key\":3,\"name\":\"wallet_public_key\",\"type\":\"uint256\"},{\"key\":4,\"name\":\"owner_address\",\"type\":\"address\"}],\"events\":[],\"functions\":[{\"inputs\":[],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getDetails\",\"outputs\":[{\"components\":[{\"name\":\"root_address\",\"type\":\"address\"},{\"name\":\"code\",\"type\":\"cell\"},{\"name\":\"wallet_public_key\",\"type\":\"uint256\"},{\"name\":\"owner_address\",\"type\":\"address\"},{\"name\":\"balance\",\"type\":\"uint128\"}],\"name\":\"value0\",\"type\":\"tuple\"}]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"}],\"name\":\"accept\",\"outputs\":[]},{\"inputs\":[],\"name\":\"allowance\",\"outputs\":[{\"components\":[{\"name\":\"remaining_tokens\",\"type\":\"uint128\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"value0\",\"type\":\"tuple\"}]},{\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"remaining_tokens\",\"type\":\"uint128\"},{\"name\":\"tokens\",\"type\":\"uint128\"}],\"name\":\"approve\",\"outputs\":[]},{\"inputs\":[],\"name\":\"disapprove\",\"outputs\":[]},{\"inputs\":[{\"name\":\"recipient_public_key\",\"type\":\"uint256\"},{\"name\":\"recipient_address\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"deploy_grams\",\"type\":\"uint128\"},{\"name\":\"transfer_grams\",\"type\":\"uint128\"}],\"name\":\"transferToRecipient\",\"outputs\":[]},{\"inputs\":[{\"name\":\"recipient_public_key\",\"type\":\"uint256\"},{\"name\":\"recipient_address\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"deploy_grams\",\"type\":\"uint128\"},{\"name\":\"transfer_grams\",\"type\":\"uint128\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"transferToRecipientWithNotify\",\"outputs\":[]},{\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"}],\"name\":\"transfer\",\"outputs\":[]},{\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"transferWithNotify\",\"outputs\":[]},{\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"}],\"name\":\"transferFrom\",\"outputs\":[]},{\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"transferFromWithNotify\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"sender_public_key\",\"type\":\"uint256\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"send_gas_to\",\"type\":\"address\"},{\"name\":\"notify_receiver\",\"type\":\"bool\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"internalTransfer\",\"outputs\":[]},{\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"send_gas_to\",\"type\":\"address\"},{\"name\":\"notify_receiver\",\"type\":\"bool\"},{\"name\":\"payload\",\"type\":\"cell\"}],\"name\":\"internalTransferFrom\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"burnByOwner\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"burnByRoot\",\"outputs\":[]},{\"inputs\":[{\"name\":\"receive_callback_\",\"type\":\"address\"}],\"name\":\"setReceiveCallback\",\"outputs\":[]},{\"inputs\":[{\"name\":\"gas_dest\",\"type\":\"address\"}],\"name\":\"destroy\",\"outputs\":[]},{\"inputs\":[],\"name\":\"balance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"receive_callback\",\"outputs\":[{\"name\":\"receive_callback\",\"type\":\"address\"}]},{\"inputs\":[],\"name\":\"target_gas_balance\",\"outputs\":[{\"name\":\"target_gas_balance\",\"type\":\"uint128\"}]}],\"header\":[\"pubkey\",\"time\",\"expire\"]}"
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
func (contract *TONTokenWalletContract) TransferToRecipient(recipient_public_key string, recipient_address string, tokens string, deploy_grams string, transfer_grams string) *ContractMethod {
	input := fmt.Sprintf("{\"recipient_public_key\": \"%s\" ,\"recipient_address\": \"%s\" ,\"tokens\": \"%s\" ,\"deploy_grams\": \"%s\" ,\"transfer_grams\": \"%s\" }", recipient_public_key, recipient_address, tokens, deploy_grams, transfer_grams)
	return contract.CallContractMethod("transferToRecipient", input)
}
func (contract *TONTokenWalletContract) TransferToRecipientWithNotify(recipient_public_key string, recipient_address string, tokens string, deploy_grams string, transfer_grams string, payload string) *ContractMethod {
	input := fmt.Sprintf("{\"recipient_public_key\": \"%s\" ,\"recipient_address\": \"%s\" ,\"tokens\": \"%s\" ,\"deploy_grams\": \"%s\" ,\"transfer_grams\": \"%s\" ,\"payload\": \"%s\" }", recipient_public_key, recipient_address, tokens, deploy_grams, transfer_grams, payload)
	return contract.CallContractMethod("transferToRecipientWithNotify", input)
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
