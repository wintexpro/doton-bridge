package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	VoteControllerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_proposalCode\",\"type\":\"cell\"},{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"},{\"name\":\"_publicKey\",\"type\":\"uint256\"},{\"name\":\"_proposalVotersAmount\",\"type\":\"uint256\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"},{\"name\":\"initializerChoice\",\"type\":\"uint8\"},{\"name\":\"initializerAddress\",\"type\":\"address\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"messageType\",\"type\":\"uint256\"}],\"name\":\"createProposal\",\"outputs\":[{\"name\":\"proposalAddress\",\"type\":\"address\"}]},{\"inputs\":[],\"name\":\"getDeployInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"}],\"name\":\"setDeployInitialValue\",\"outputs\":[{\"name\":\"value0\",\"type\":\"uint128\"}]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"getProposalAddress\",\"outputs\":[{\"name\":\"proposal\",\"type\":\"address\"}]}],\"header\":[\"time\"]}"
	VoteControllerTvc = "te6ccgECGQEABHkAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCxcFBBgCkiHbPNMAAY4SgQIA1xgg+QFY+EIg+GX5EPKo3tM/AY4d+EMhuSCfMCD4I4ED6KiCCBt3QKC53pMg+GPg8jTYMNMfAds8+Edu8nwVBgE0ItDXCwOpOADcIccA3CHTHyHdAds8+Edu8nwGA0AgghBKadv/u46A4CCCEFrNhWe7joDgIIIQYUsu0LrjAg8JBwNwMPhBbuMA0ds8IcD/jiMj0NMB+kAwMcjPhyDOgGDPQM+Bz4HPk4Usu0Ihzwt/yXD7AN4w4wB/+GcSCBQABPhMAiggghBPYtfguuMCIIIQWs2FZ7rjAgwKA3Qw+EFu4wDTf9HbPCHA/44jI9DTAfpAMDHIz4cgzoBgz0DPgc+Bz5NrNhWeIc8Lf8lw+wDeMNs8f/hnEgsUAEJw+E8gwQKTMIBk3vhFIG6SMHDe+Eu68vT4ACH4bPhMMTEDpjD4QW7jANMH0z/U0wf6QPpBldTR0PpA39cN/5XU0dDT/9/R2zwhwP+OIiPQ0wH6QDAxyM+HIM6AYM9Az4HPgc+TPYtfgiHPFslw+wDeMNs8f/hnEg0UAf6NCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAT4AG34QsjL/3BYgED0QyjIywdxWIBA9EMnyMs/cliAQPRD+ChzWIBA9BYmdFiAQPQXyPQAyfhKyM+EgPQA9ADPgckg+QDIz4oAQMv/ydD4TCHIz4WIzgH6AoBpDgBkz0DPg8+DIs8Uz4PIz5FBDSnW+E7PC/8ozwsHJ88WJs8WyCbPC//Nzclw+wAxMSAxbHECKCCCEBsCOW664wIgghBKadv/uuMCExADeDD4QW7jANMH0z/U0ds8IcD/jiIj0NMB+kAwMcjPhyDOgGDPQM+Bz4HPkymnb/4hzxbJcPsA3jDjAH/4ZxIRFADsjQhgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEbfhCyMv/cFiAQPRDJMjLB3FYgED0QyPIyz9yWIBA9EP4KHNYgED0FiJ0WIBA9BfI9ADJ+ErIz4SA9AD0AM+BySD5AMjPigBAy//J0GwSATBsMQBY7UTQ0//TP9MA1dcL//hu1NP/03/T/9cLB/hv+G34bPhr+Gp/+GH4Zvhj+GICZDD4QW7jAPhG8nNx+GbU03/T/9cN/5XU0dDT/9/R+AAj+Goi+Gwh+Gsg+G5fBNs8f/hnFRQAXvhCyMv/+EPPCz/4Rs8LAMj4TgHL//hK+Ev4TPhN+E9eUM8RzMv/y3/L/8sHye1UAWztRNAg10nCAY4p0//TP9MA1dcL//hu1NP/03/T/9cLB/hv+G34bPhr+Gp/+GH4Zvhj+GKOgOIWAVz0BYj4anD4a3D4bHD4bXD4bnD4b3ABgED0DvK91wv/+GJw+GNw+GZ/+GGAZfhvGAEK9KQg9KEYAAA="
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
	params := json.RawMessage(fmt.Sprintf("{\"_proposalCode\": \"%s\" ,\"_deployInitialValue\": \"%s\" ,\"_publicKey\": \"%s\" ,\"_proposalVotersAmount\": \"%s\" }", voteControllerDeployParams.ProposalCode, voteControllerDeployParams.DeployInitialValue, voteControllerDeployParams.PublicKey, voteControllerDeployParams.ProposalVotersAmount))
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
