package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	Tip3HandlerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_proposalCode\",\"type\":\"cell\"},{\"name\":\"_epochControllerPubKey\",\"type\":\"uint256\"},{\"name\":\"_tip3RootAddress\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"epochAddress\",\"type\":\"address\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"executeProposal\",\"outputs\":[]}],\"header\":[\"time\"]}"
	Tip3HandlerTvc = "te6ccgECEAEAApgAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCw4FBA8C1o0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIds80wABn4ECANcYIPkBWPhC+RDyqN7TPwGOHfhDIbkgnzAg+COBA+iogggbd0Cgud6TIPhj4PI02DDTHwHbPPhHbvJ8CAYBPCLQ0wP6QDD4aak4ANwhxwDcIdMfId0B2zz4R27yfAYCKCCCEA6w50W64wIgghBFgufyuuMCCgcCXjD4QW7jAPhG8nNx+GbU0//6QZXU0dD6QN/RIiL4ACH4aiD4a1sg+G1fA9s8f/hnCAwBUu1E0CDXScIBjhzT/9M/0wDU0//TB/ht+Gz4a/hqf/hh+Gb4Y/hijoDiCQGW9AWI+Gpw+Gtw+GyNCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAT4bXABgED0DvK91wv/+GJw+GNw+GZ/+GGAZfhsDwL2MPhBbuMA+kDTB9M/1w3/ldTR0NP/39TRJCQkI234S8jL/3BYgED0QyPIywdxWIBA9EMiyMs/cliAQPRDJHNYgED0FiF0WIBA9BfI9ADJ+ErIz4SA9AD0AM+ByfhMIMECkzCAZN74SfpCbxPXC/8i+QC68vT4TcjPhQjODQsBVo0EUBfXhAAAAAAAAAAAAAAAAAABzxbPgc+DJs8UyXD7ADBfBF8F4wB/+GcMAET4QsjL//hDzws/+EbPCwD4SvhL+Ez4TV4wzMv/ywfOye1UAD7tRNDT/9M/0wDU0//TB/ht+Gz4a/hqf/hh+Gb4Y/hiAQr0pCD0oQ8AAA=="
)

type Tip3HandlerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type Tip3Handler struct {
	Ctx ContractContext
}
type Tip3HandlerDeployParams struct {
	ProposalCode          string
	EpochControllerPubKey string
	Tip3RootAddress       string
}

func (c *Tip3Handler) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{Tip3HandlerTvc})
}
func (c *Tip3Handler) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(Tip3HandlerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *Tip3Handler) Address() (string, error) {
	tip3HandlerDeployParams := Tip3HandlerDeployParams{
		EpochControllerPubKey: "0x7",
		ProposalCode:          "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		Tip3RootAddress:       "0:0000000000000000000000000000000000000000000000000000000000000000",
	}
	encodeMessage, err := c.DeployEncodeMessage(&tip3HandlerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *Tip3Handler) New(address string) (*Tip3HandlerContract, error) {
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
	contract := Tip3HandlerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *Tip3Handler) DeployEncodeMessage(tip3HandlerDeployParams *Tip3HandlerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         Tip3HandlerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_proposalCode\": \"%s\" ,\"_epochControllerPubKey\": \"%s\" ,\"_tip3RootAddress\": \"%s\" }", tip3HandlerDeployParams.ProposalCode, tip3HandlerDeployParams.EpochControllerPubKey, tip3HandlerDeployParams.Tip3RootAddress))
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
func (c *Tip3Handler) Deploy(tip3HandlerDeployParams *Tip3HandlerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*Tip3HandlerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(tip3HandlerDeployParams)
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
	contract := Tip3HandlerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *Tip3HandlerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *Tip3HandlerContract) ExecuteProposal(epochAddress string, chainId string, nonce string, messageType string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"epochAddress\": \"%s\" ,\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"messageType\": \"%s\" ,\"data\": \"%s\" }", epochAddress, chainId, nonce, messageType, data)
	return contract.CallContractMethod("executeProposal", input)
}
func (contract *Tip3HandlerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *Tip3HandlerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *Tip3HandlerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
