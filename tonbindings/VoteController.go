package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	VoteControllerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_proposalCode\",\"type\":\"cell\"},{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"},{\"name\":\"_publicKey\",\"type\":\"uint256\"},{\"name\":\"_proposalPublicKey\",\"type\":\"uint256\"},{\"name\":\"_proposalVotersAmount\",\"type\":\"uint256\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"},{\"name\":\"initializerChoice\",\"type\":\"uint8\"},{\"name\":\"initializerAddress\",\"type\":\"address\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"messageType\",\"type\":\"uint256\"}],\"name\":\"createProposal\",\"outputs\":[{\"name\":\"proposalAddress\",\"type\":\"address\"}]},{\"inputs\":[],\"name\":\"getDeployInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"}],\"name\":\"setDeployInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"getProposalAddress\",\"outputs\":[{\"name\":\"proposal\",\"type\":\"address\"}]}],\"header\":[\"time\"]}"
	VoteControllerTvc = "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"
)

type VoteControllerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type VoteController struct {
	Ctx ContractContext
}
type VoteControllerDeployParams struct {
	ProposalCode         string
	DeployInitialValue   string
	PublicKey            string
	ProposalPublicKey    string
	ProposalVotersAmount string
}

func (c *VoteController) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{VoteControllerTvc})
}
func (c *VoteController) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(VoteControllerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *VoteController) Address() (string, error) {
	voteControllerDeployParams := VoteControllerDeployParams{
		DeployInitialValue:   "0x6",
		ProposalCode:         "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		ProposalPublicKey:    "0x7",
		ProposalVotersAmount: "0x7",
		PublicKey:            "0x7",
	}
	encodeMessage, err := c.DeployEncodeMessage(&voteControllerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *VoteController) New(address string) (*VoteControllerContract, error) {
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
	contract := VoteControllerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *VoteController) DeployEncodeMessage(voteControllerDeployParams *VoteControllerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         VoteControllerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_proposalCode\": \"%s\" ,\"_deployInitialValue\": \"%s\" ,\"_publicKey\": \"%s\" ,\"_proposalPublicKey\": \"%s\" ,\"_proposalVotersAmount\": \"%s\" }", voteControllerDeployParams.ProposalCode, voteControllerDeployParams.DeployInitialValue, voteControllerDeployParams.PublicKey, voteControllerDeployParams.ProposalPublicKey, voteControllerDeployParams.ProposalVotersAmount))
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
func (c *VoteController) Deploy(voteControllerDeployParams *VoteControllerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*VoteControllerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(voteControllerDeployParams)
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
	contract := VoteControllerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *VoteControllerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *VoteControllerContract) CreateProposal(chainId string, nonce string, data string, initializerChoice string, initializerAddress string, handlerAddress string, messageType string) *ContractMethod {
	input := fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" ,\"initializerChoice\": \"%s\" ,\"initializerAddress\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"messageType\": \"%s\" }", chainId, nonce, data, initializerChoice, initializerAddress, handlerAddress, messageType)
	return contract.CallContractMethod("createProposal", input)
}
func (contract *VoteControllerContract) GetDeployInitialValue() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getDeployInitialValue", input)
}
func (contract *VoteControllerContract) SetDeployInitialValue(_deployInitialValue string) *ContractMethod {
	input := fmt.Sprintf("{\"_deployInitialValue\": \"%s\" }", _deployInitialValue)
	return contract.CallContractMethod("setDeployInitialValue", input)
}
func (contract *VoteControllerContract) GetProposalAddress(chainId string, nonce string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" }", chainId, nonce, data)
	return contract.CallContractMethod("getProposalAddress", input)
}
func (contract *VoteControllerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *VoteControllerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *VoteControllerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
