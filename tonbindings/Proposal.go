package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	ProposalAbi = "{\"ABI version\":2,\"data\":[{\"key\":1,\"name\":\"chainId\",\"type\":\"uint8\"},{\"key\":2,\"name\":\"nonce\",\"type\":\"uint64\"},{\"key\":3,\"name\":\"voteControllerAddress\",\"type\":\"address\"},{\"key\":4,\"name\":\"data\",\"type\":\"cell\"}],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_votersAmount\",\"type\":\"uint256\"},{\"name\":\"initializerChoice\",\"type\":\"uint8\"},{\"name\":\"initializerAddress\",\"type\":\"address\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"messageType\",\"type\":\"uint256\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getYesVotes\",\"outputs\":[{\"name\":\"yesVotes\",\"type\":\"uint256\"}]},{\"inputs\":[],\"name\":\"getNoVotes\",\"outputs\":[{\"name\":\"noVotes\",\"type\":\"uint256\"}]},{\"inputs\":[{\"name\":\"voter\",\"type\":\"address\"},{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"}],\"name\":\"voteByController\",\"outputs\":[]}],\"header\":[\"time\"]}"
	ProposalTvc = "te6ccgECFwEABK4AAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCxUFBBYC1o0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIds80wABn4ECANcYIPkBWPhC+RDyqN7TPwGOHfhDIbkgnzAg+COBA+iogggbd0Cgud6TIPhj4PI02DDTHwHbPPhHbvJ8CwYBPCLQ0wP6QDD4aak4ANwhxwDcIdMfId0B2zz4R27yfAYEUCCCEEDwDw664wIgghBNl5PpuuMCIIIQUENKdbrjAiCCEH2iaIa64wIRDwkHA3Aw+EFu4wDR2zwhwP+OIyPQ0wH6QDAxyM+HIM6AYM9Az4HPgc+T9omiGiHPC//JcPsA3jDjAH/4ZxQIEwAecHD4T3j0DpPXC/+RcOIxAtIw+EFu4wD4RvJzcfhm0//TB/pBldTR0PpA3/pBldTR0PpA39cN/5XU0dDT/9/R+Ez4SccF8uBn+AAk+G74TyQBUxB49A6T1wv/kXDipMjL/1l49EP4b/hQIwElyMsHWYEBC/RB+HAkwAELCgGkjkghcMjPhYDKAHPPQM6NBFAvrwgAAAAAAAAAAAAAAAAAAc8Wz4HPg8jPkDrDnRb4TM8W+ErPCwf4S88LPyLPC//4Tc8Uzclw+wDeXwXbPH/4ZxMBiO1E0CDXScIBjjfT/9M/0wDTB9M/+kDU0//0BPQE0wfTB9cLB/hz+HL4cfhw+G/4bvht+Gz4a/hqf/hh+Gb4Y/hijoDiDAL+9AVxIYBA9A6T1wsHkXDi+GpyIYBA9A6T1ws/kXDi+GtzIYBA9A6OJI0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABN/4bHQhgED0D46A3/htcPhubfhvbfhwcPhxcPhycPhzcAGAQPQO8r3XC//4YnD4Y3D4Zg4NAB5/+GGAZfhxgGb4coBn+HMBAogWA3Aw+EFu4wDR2zwhwP+OIyPQ0wH6QDAxyM+HIM6AYM9Az4HPgc+TNl5PpiHPC//JcPsA3jDjAH/4ZxQQEwAecHH4T3j0DpPXC/+RcOIxAv4w+EFu4wD6QNMH1w3/ldTR0NP/3/pBldTR0PpA39H4USDBApMwgGTe+En4TMcF8vT4UiDBApMwgGTeJPhQgQEL9AogkTHes/L0+E8jAVMQePQOk9cL/5Fw4qTIy/9ZePRD+G/4UCQBJMjLB1mBAQv0QfhwcfhPePQOk9cL/5FwFBIByOJw+E949A6T1wv/kXDioPhOvo5IIHDIz4WAygBzz0DOjQRQL68IAAAAAAAAAAAAAAAAAAHPFs+Bz4PIz5A6w50W+EzPFvhKzwsH+EvPCz8jzwv/+E3PFM3JcPsA3l8E2zx/+GcTAHT4QsjL//hDzws/+EbPCwD4SvhL+Ez4TfhO+E/4UPhR+FL4U16QywfLP87My//0APQAywfLB8sHye1UAHTtRNDT/9M/0wDTB9M/+kDU0//0BPQE0wfTB9cLB/hz+HL4cfhw+G/4bvht+Gz4a/hqf/hh+Gb4Y/hiAQr0pCD0oRYAAA=="
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
	params := json.RawMessage(fmt.Sprintf("{\"_votersAmount\": \"%s\" ,\"initializerChoice\": \"%s\" ,\"initializerAddress\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"messageType\": \"%s\" }", proposalDeployParams.VotersAmount, proposalDeployParams.InitializerChoice, proposalDeployParams.InitializerAddress, proposalDeployParams.HandlerAddress, proposalDeployParams.MessageType))
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
