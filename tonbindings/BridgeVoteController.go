package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	BridgeVoteControllerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_proposalCode\",\"type\":\"cell\"},{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"},{\"name\":\"_publicKey\",\"type\":\"uint256\"},{\"name\":\"_proposalPublicKey\",\"type\":\"uint256\"},{\"name\":\"_proposalVotersAmount\",\"type\":\"uint256\"},{\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"voter\",\"type\":\"address\"},{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"voteByBridge\",\"outputs\":[]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"},{\"name\":\"initializerChoice\",\"type\":\"uint8\"},{\"name\":\"initializerAddress\",\"type\":\"address\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"messageType\",\"type\":\"uint256\"}],\"name\":\"createProposal\",\"outputs\":[{\"name\":\"proposalAddress\",\"type\":\"address\"}]},{\"inputs\":[],\"name\":\"getDeployInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"}],\"name\":\"setDeployInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"getProposalAddress\",\"outputs\":[{\"name\":\"proposal\",\"type\":\"address\"}]}],\"header\":[\"time\"]}"
	BridgeVoteControllerTvc = "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"
)

type BridgeVoteControllerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type BridgeVoteController struct {
	Ctx ContractContext
}
type BridgeVoteControllerDeployParams struct {
	ProposalCode         string
	DeployInitialValue   string
	PublicKey            string
	ProposalPublicKey    string
	ProposalVotersAmount string
	BridgeAddress        string
}

func (c *BridgeVoteController) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{BridgeVoteControllerTvc})
}
func (c *BridgeVoteController) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(BridgeVoteControllerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *BridgeVoteController) Address() (string, error) {
	bridgeVoteControllerDeployParams := BridgeVoteControllerDeployParams{
		BridgeAddress:        "0:0000000000000000000000000000000000000000000000000000000000000000",
		DeployInitialValue:   "0x6",
		ProposalCode:         "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		ProposalPublicKey:    "0x7",
		ProposalVotersAmount: "0x7",
		PublicKey:            "0x7",
	}
	encodeMessage, err := c.DeployEncodeMessage(&bridgeVoteControllerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *BridgeVoteController) New(address string) (*BridgeVoteControllerContract, error) {
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
	contract := BridgeVoteControllerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *BridgeVoteController) DeployEncodeMessage(bridgeVoteControllerDeployParams *BridgeVoteControllerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         BridgeVoteControllerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_proposalCode\": \"%s\" ,\"_deployInitialValue\": \"%s\" ,\"_publicKey\": \"%s\" ,\"_proposalPublicKey\": \"%s\" ,\"_proposalVotersAmount\": \"%s\" ,\"_bridgeAddress\": \"%s\" }", bridgeVoteControllerDeployParams.ProposalCode, bridgeVoteControllerDeployParams.DeployInitialValue, bridgeVoteControllerDeployParams.PublicKey, bridgeVoteControllerDeployParams.ProposalPublicKey, bridgeVoteControllerDeployParams.ProposalVotersAmount, bridgeVoteControllerDeployParams.BridgeAddress))
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
func (c *BridgeVoteController) Deploy(bridgeVoteControllerDeployParams *BridgeVoteControllerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*BridgeVoteControllerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(bridgeVoteControllerDeployParams)
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
	contract := BridgeVoteControllerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *BridgeVoteControllerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *BridgeVoteControllerContract) VoteByBridge(voter string, choice string, chainId string, messageType string, handlerAddress string, nonce string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"voter\": \"%s\" ,\"choice\": \"%s\" ,\"chainId\": \"%s\" ,\"messageType\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" }", voter, choice, chainId, messageType, handlerAddress, nonce, data)
	return contract.CallContractMethod("voteByBridge", input)
}
func (contract *BridgeVoteControllerContract) CreateProposal(chainId string, nonce string, data string, initializerChoice string, initializerAddress string, handlerAddress string, messageType string) *ContractMethod {
	input := fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" ,\"initializerChoice\": \"%s\" ,\"initializerAddress\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"messageType\": \"%s\" }", chainId, nonce, data, initializerChoice, initializerAddress, handlerAddress, messageType)
	return contract.CallContractMethod("createProposal", input)
}
func (contract *BridgeVoteControllerContract) GetDeployInitialValue() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getDeployInitialValue", input)
}
func (contract *BridgeVoteControllerContract) SetDeployInitialValue(_deployInitialValue string) *ContractMethod {
	input := fmt.Sprintf("{\"_deployInitialValue\": \"%s\" }", _deployInitialValue)
	return contract.CallContractMethod("setDeployInitialValue", input)
}
func (contract *BridgeVoteControllerContract) GetProposalAddress(chainId string, nonce string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" }", chainId, nonce, data)
	return contract.CallContractMethod("getProposalAddress", input)
}
func (contract *BridgeVoteControllerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *BridgeVoteControllerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *BridgeVoteControllerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
