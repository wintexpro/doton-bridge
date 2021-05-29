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
