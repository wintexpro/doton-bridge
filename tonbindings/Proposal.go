package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	ProposalAbi = "{\"ABI version\":2,\"data\":[{\"key\":1,\"name\":\"chainId\",\"type\":\"uint8\"},{\"key\":2,\"name\":\"nonce\",\"type\":\"uint64\"},{\"key\":3,\"name\":\"voteControllerAddress\",\"type\":\"address\"},{\"key\":4,\"name\":\"data\",\"type\":\"cell\"}],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_publicKey\",\"type\":\"uint256\"},{\"name\":\"_votersAmount\",\"type\":\"uint256\"},{\"name\":\"initializerChoice\",\"type\":\"uint8\"},{\"name\":\"initializerAddress\",\"type\":\"address\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"messageType\",\"type\":\"uint256\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getYesVotes\",\"outputs\":[{\"name\":\"yesVotes\",\"type\":\"uint256\"}]},{\"inputs\":[],\"name\":\"getNoVotes\",\"outputs\":[{\"name\":\"noVotes\",\"type\":\"uint256\"}]},{\"inputs\":[{\"name\":\"voter\",\"type\":\"address\"},{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"}],\"name\":\"voteByController\",\"outputs\":[]}],\"header\":[\"time\"]}"
	ProposalTvc = "te6ccgECHQEABckAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAib/APSkICLAAZL0oOGK7VNYMPShCAQBCvSkIPShBQIJnQAAAAwHBgBn12omhp/+mf6YBq64X//Dfpg+mf/SBqaf/6AnoC/Dj8OHw3fDb8Nnw1/DU//DD8M3wx/DFABt98IWRl//wh54Wf/CNnhYBkfCeA5f/8JXwl/CZ8JvwnfCh8KK84Z4jlg+Wf52Zl//oAegBk9qpAIBIAwJAcb/f40IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIe1E0CDXScIBjjDT/9M/0wDV1wv/+G/TB9M/+kDU0//0BPQF+HH4cPhu+G34bPhr+Gp/+GH4Zvhj+GIKAf6OffQFcSGAQPQOk9cLB5Fw4vhqciGAQPQOk9cLP5Fw4vhrcyGAQPQOjiSNCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATf+Gx0IYBA9A+SyMnf+G1w+G5w+G9t+HBt+HFwAYBA9A7yvdcL//hicPhjcPhmf/hhCwCs4tMAAZ+BAgDXGCD5AVj4QvkQ8qje0z8Bjh34QyG5IJ8wIPgjgQPoqIIIG3dAoLnekvhj4DDyNNjTHyHBAyKCEP////28sZNb8jzgAfAB+EdukzDyPN4CASAVDQIBIA8OAZ+72iaIb4QW6S8AzeIZnTH/hEWG91+GTf0XBw+FB49A6T1wv/kXDiMSHA/44jI9DTAfpAMDHIz4cgzoBgz0DPgc+Bz5P2iaIaIc8L/8lx+wCBECAUgSEAGft2Xk+n4QW6S8AzeIZnTH/hEWG91+GTf0XBx+FB49A6T1wv/kXDiMSHA/44jI9DTAfpAMDHIz4cgzoBgz0DPgc+Bz5M2Xk+mIc8L/8lx+wCARAISON/hEIG8TIW8S+ElVAm8RyHLPQMoAc89AzgH6AvQAgGjPQM+Bz4H4RG8VzwsfIc8L/8n4RG8U+wDiMJLwC95/+GcBCbY8A8OgEwH++EFukvAM3vpA0wfXDf+V1NHQ0//f+kGV1NHQ+kDf0fhJ+EzHBfLgZCP4UYEBC/QKIJEx3rPy4GT4UCMBUxB49A6T1wv/kXDipMjL/1l49EP4cPhRJAEkyMsHWYEBC/RB+HFx+FB49A6T1wv/kXDicPhQePQOk9cL/5Fw4qD4TxQAqL6OSSBwyM+FgMoAc89Azo0EUC+vCAAAAAAAAAAAAAAAAAABzxbPgc+DyM+QoQtcYvhCzwv/+ErPCwf4S88LPyPPC//4Tc8Uzclx+wDeXwTwC3/4ZwIBIBwWAQ+76f5+/4QW6BcC5I6A3vhG8nNx+GbT/9cN/5XU0dDT/9/XDQeV1NHQ0wff+kGV1NHQ+kDf+kGV1NHQ+kDf1w3/ldTR0NP/39H4ACX4biT4b/hQJAFTEHj0DpPXC/+RcOKkyMv/WXj0Q/hw+FEjASXIywdZgQEL9EH4cSTAARkYAKaOSSFwyM+FgMoAc89Azo0EUC+vCAAAAAAAAAAAAAAAAAABzxbPgc+DyM+QoQtcYvhCzwv/+ErPCwf4S88LPyLPC//4Tc8Uzclx+wDeXwbwC3/4ZwF07UTQINdJwgGOMNP/0z/TANXXC//4b9MH0z/6QNTT//QE9AX4cfhw+G74bfhs+Gv4an/4Yfhm+GP4YhoB/o599AVxIYBA9A6T1wsHkXDi+GpyIYBA9A6T1ws/kXDi+GtzIYBA9A6OJI0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABN/4bHQhgED0D5LIyd/4bXD4bnD4b234cG34cXABgED0DvK91wv/+GJw+GNw+GZ/+GEbAALiAJDdcCLQ0wP6QDD4aak4APhEf29xggiYloBvcm1vc3FvdPhk3CHHANwh0x8h3SHBAyKCEP////28sZNb8jzgAfAB+EdukzDyPN4="
)

type ProposalContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type Proposal struct {
	Ctx ContractContext
}
type ProposalDeployParams struct {
	PublicKey          string
	VotersAmount       string
	InitializerChoice  string
	InitializerAddress string
	HandlerAddress     string
	MessageType        string
}
type ProposalInitVars struct {
	ChainId               string
	Nonce                 string
	VoteControllerAddress string
	Data                  string
}

func (c *Proposal) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{ProposalTvc})
}
func (c *Proposal) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(ProposalAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *Proposal) Address(proposalInitVars *ProposalInitVars) (string, error) {
	proposalDeployParams := ProposalDeployParams{
		HandlerAddress:     "0:0000000000000000000000000000000000000000000000000000000000000000",
		InitializerAddress: "0:0000000000000000000000000000000000000000000000000000000000000000",
		InitializerChoice:  "0x2",
		MessageType:        "0x7",
		PublicKey:          "0x7",
		VotersAmount:       "0x7",
	}
	encodeMessage, err := c.DeployEncodeMessage(&proposalDeployParams, proposalInitVars)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *Proposal) New(address string, proposalInitVars *ProposalInitVars) (*ProposalContract, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	if address == "" {
		address, err = c.Address(proposalInitVars)
		if err != nil {
			return nil, err
		}
	}
	contract := ProposalContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *Proposal) DeployEncodeMessage(proposalDeployParams *ProposalDeployParams, proposalInitVars *ProposalInitVars) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	initialVars := json.RawMessage(fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"voteControllerAddress\": \"%s\" ,\"data\": \"%s\" }", proposalInitVars.ChainId, proposalInitVars.Nonce, proposalInitVars.VoteControllerAddress, proposalInitVars.Data))
	deploySet := client.DeploySet{
		InitialData: initialVars,
		Tvc:         ProposalTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_publicKey\": \"%s\" ,\"_votersAmount\": \"%s\" ,\"initializerChoice\": \"%s\" ,\"initializerAddress\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"messageType\": \"%s\" }", proposalDeployParams.PublicKey, proposalDeployParams.VotersAmount, proposalDeployParams.InitializerChoice, proposalDeployParams.InitializerAddress, proposalDeployParams.HandlerAddress, proposalDeployParams.MessageType))
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
func (c *Proposal) Deploy(proposalDeployParams *ProposalDeployParams, proposalInitVars *ProposalInitVars, messageCallback func(event *client.ProcessingEvent)) (*ProposalContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(proposalDeployParams, proposalInitVars)
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
	contract := ProposalContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *ProposalContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *ProposalContract) GetYesVotes() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getYesVotes", input)
}
func (contract *ProposalContract) GetNoVotes() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getNoVotes", input)
}
func (contract *ProposalContract) VoteByController(voter string, choice string, messageType string, handlerAddress string) *ContractMethod {
	input := fmt.Sprintf("{\"voter\": \"%s\" ,\"choice\": \"%s\" ,\"messageType\": \"%s\" ,\"handlerAddress\": \"%s\" }", voter, choice, messageType, handlerAddress)
	return contract.CallContractMethod("voteByController", input)
}
func (contract *ProposalContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *ProposalContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *ProposalContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
