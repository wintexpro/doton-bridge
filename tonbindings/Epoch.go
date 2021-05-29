package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	EpochAbi = "{\"ABI version\":2,\"data\":[{\"key\":1,\"name\":\"number\",\"type\":\"uint64\"},{\"key\":2,\"name\":\"voteControllerAddress\",\"type\":\"address\"}],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_proposalCode\",\"type\":\"cell\"},{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"},{\"name\":\"_publicKey\",\"type\":\"uint256\"},{\"name\":\"_proposalVotersAmount\",\"type\":\"uint256\"},{\"name\":\"_firstEraDuration\",\"type\":\"uint32\"},{\"name\":\"_secondEraDuration\",\"type\":\"uint32\"},{\"name\":\"_publicRandomness\",\"type\":\"uint256\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"registeringRelayer\",\"type\":\"address\"},{\"name\":\"signHighPart\",\"type\":\"uint256\"},{\"name\":\"signLowPart\",\"type\":\"uint256\"},{\"name\":\"pubkey\",\"type\":\"uint256\"}],\"name\":\"signup\",\"outputs\":[]},{\"inputs\":[{\"name\":\"signHighPart\",\"type\":\"uint256\"},{\"name\":\"signLowPart\",\"type\":\"uint256\"},{\"name\":\"pubkey\",\"type\":\"uint256\"}],\"name\":\"forceEra\",\"outputs\":[]},{\"inputs\":[{\"name\":\"voter\",\"type\":\"address\"},{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"voteByEpochController\",\"outputs\":[]},{\"inputs\":[{\"name\":\"relayer\",\"type\":\"address\"}],\"name\":\"isChoosen\",\"outputs\":[{\"name\":\"value0\",\"type\":\"bool\"}]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"},{\"name\":\"initializerChoice\",\"type\":\"uint8\"},{\"name\":\"initializerAddress\",\"type\":\"address\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"messageType\",\"type\":\"uint256\"}],\"name\":\"createProposal\",\"outputs\":[{\"name\":\"proposalAddress\",\"type\":\"address\"}]},{\"inputs\":[],\"name\":\"getDeployInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"}],\"name\":\"setDeployInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"getProposalAddress\",\"outputs\":[{\"name\":\"proposal\",\"type\":\"address\"}]},{\"inputs\":[],\"name\":\"firstEraEndsAt\",\"outputs\":[{\"name\":\"firstEraEndsAt\",\"type\":\"uint32\"}]},{\"inputs\":[],\"name\":\"secondEraEndsAt\",\"outputs\":[{\"name\":\"secondEraEndsAt\",\"type\":\"uint32\"}]}],\"header\":[\"time\"]}"
	EpochTvc = "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"
)

type EpochContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type Epoch struct {
	Ctx ContractContext
}
type EpochDeployParams struct {
	ProposalCode         string
	DeployInitialValue   string
	PublicKey            string
	ProposalVotersAmount string
	FirstEraDuration     string
	SecondEraDuration    string
	PublicRandomness     string
}
type EpochInitVars struct {
	Number                string
	VoteControllerAddress string
}

func (c *Epoch) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{EpochTvc})
}
func (c *Epoch) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(EpochAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *Epoch) Address(epochInitVars *EpochInitVars) (string, error) {
	epochDeployParams := EpochDeployParams{
		DeployInitialValue:   "0x6",
		FirstEraDuration:     "0x4",
		ProposalCode:         "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		ProposalVotersAmount: "0x7",
		PublicKey:            "0x7",
		PublicRandomness:     "0x7",
		SecondEraDuration:    "0x4",
	}
	encodeMessage, err := c.DeployEncodeMessage(&epochDeployParams, epochInitVars)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *Epoch) New(address string, epochInitVars *EpochInitVars) (*EpochContract, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	if address == "" {
		address, err = c.Address(epochInitVars)
		if err != nil {
			return nil, err
		}
	}
	contract := EpochContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *Epoch) DeployEncodeMessage(epochDeployParams *EpochDeployParams, epochInitVars *EpochInitVars) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	initialVars := json.RawMessage(fmt.Sprintf("{\"number\": \"%s\" ,\"voteControllerAddress\": \"%s\" }", epochInitVars.Number, epochInitVars.VoteControllerAddress))
	deploySet := client.DeploySet{
		InitialData: initialVars,
		Tvc:         EpochTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_proposalCode\": \"%s\" ,\"_deployInitialValue\": \"%s\" ,\"_publicKey\": \"%s\" ,\"_proposalVotersAmount\": \"%s\" ,\"_firstEraDuration\": \"%s\" ,\"_secondEraDuration\": \"%s\" ,\"_publicRandomness\": \"%s\" }", epochDeployParams.ProposalCode, epochDeployParams.DeployInitialValue, epochDeployParams.PublicKey, epochDeployParams.ProposalVotersAmount, epochDeployParams.FirstEraDuration, epochDeployParams.SecondEraDuration, epochDeployParams.PublicRandomness))
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
func (c *Epoch) Deploy(epochDeployParams *EpochDeployParams, epochInitVars *EpochInitVars, messageCallback func(event *client.ProcessingEvent)) (*EpochContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(epochDeployParams, epochInitVars)
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
	contract := EpochContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *EpochContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *EpochContract) Signup(registeringRelayer string, signHighPart string, signLowPart string, pubkey string) *ContractMethod {
	input := fmt.Sprintf("{\"registeringRelayer\": \"%s\" ,\"signHighPart\": \"%s\" ,\"signLowPart\": \"%s\" ,\"pubkey\": \"%s\" }", registeringRelayer, signHighPart, signLowPart, pubkey)
	return contract.CallContractMethod("signup", input)
}
func (contract *EpochContract) ForceEra(signHighPart string, signLowPart string, pubkey string) *ContractMethod {
	input := fmt.Sprintf("{\"signHighPart\": \"%s\" ,\"signLowPart\": \"%s\" ,\"pubkey\": \"%s\" }", signHighPart, signLowPart, pubkey)
	return contract.CallContractMethod("forceEra", input)
}
func (contract *EpochContract) VoteByEpochController(voter string, choice string, chainId string, messageType string, handlerAddress string, nonce string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"voter\": \"%s\" ,\"choice\": \"%s\" ,\"chainId\": \"%s\" ,\"messageType\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" }", voter, choice, chainId, messageType, handlerAddress, nonce, data)
	return contract.CallContractMethod("voteByEpochController", input)
}
func (contract *EpochContract) IsChoosen(relayer string) *ContractMethod {
	input := fmt.Sprintf("{\"relayer\": \"%s\" }", relayer)
	return contract.CallContractMethod("isChoosen", input)
}
func (contract *EpochContract) CreateProposal(chainId string, nonce string, data string, initializerChoice string, initializerAddress string, handlerAddress string, messageType string) *ContractMethod {
	input := fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" ,\"initializerChoice\": \"%s\" ,\"initializerAddress\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"messageType\": \"%s\" }", chainId, nonce, data, initializerChoice, initializerAddress, handlerAddress, messageType)
	return contract.CallContractMethod("createProposal", input)
}
func (contract *EpochContract) GetDeployInitialValue() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getDeployInitialValue", input)
}
func (contract *EpochContract) SetDeployInitialValue(_deployInitialValue string) *ContractMethod {
	input := fmt.Sprintf("{\"_deployInitialValue\": \"%s\" }", _deployInitialValue)
	return contract.CallContractMethod("setDeployInitialValue", input)
}
func (contract *EpochContract) GetProposalAddress(chainId string, nonce string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" }", chainId, nonce, data)
	return contract.CallContractMethod("getProposalAddress", input)
}
func (contract *EpochContract) FirstEraEndsAt() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("firstEraEndsAt", input)
}
func (contract *EpochContract) SecondEraEndsAt() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("secondEraEndsAt", input)
}
func (contract *EpochContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *EpochContract) abiEncodeMessageForCall(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
	callSet := client.CallSet{
		FunctionName: functionName,
		Input:        json.RawMessage(input),
	}
	paramsAbiEncodeMessage := client.ParamsOfEncodeMessage{
		Abi:     contract.Abi,
		Address: null.StringFrom(contract.Address),
		CallSet: &callSet,
		Signer:  client.Signer{EnumTypeValue: client.NoneSigner{}},
	}
	return contract.Ctx.Conn.AbiEncodeMessage(&paramsAbiEncodeMessage)
}
func (contract *EpochContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *EpochContract) call(functionName string, input string) (*client.DecodedOutput, error) {
	message, err := contract.abiEncodeMessageForCall(functionName, input)
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
