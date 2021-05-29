package tonbindings

import (
	"encoding/json"
	"fmt"

	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	MessageHandlerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"ProposalExecuted\",\"outputs\":[]}],\"functions\":[{\"inputs\":[{\"name\":\"_proposalCode\",\"type\":\"cell\"},{\"name\":\"_epochControllerPubKey\",\"type\":\"uint256\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"epochAddress\",\"type\":\"address\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"executeProposal\",\"outputs\":[]},{\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"receiveMessage\",\"outputs\":[]}],\"header\":[\"time\"]}"
	MessageHandlerTvc = "te6ccgECEQEAAs0AAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCw8FBBAC1o0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIds80wABn4ECANcYIPkBWPhC+RDyqN7TPwGOHfhDIbkgnzAg+COBA+iogggbd0Cgud6TIPhj4PI02DDTHwHbPPhHbvJ8CAYBPCLQ0wP6QDD4aak4ANwhxwDcIdMfId0B2zz4R27yfAYDPCCCEA6w50W64wIgghAPZGSHuuMCIIIQU6uwq7rjAgsKBwJEMPhBbuMA+Ebyc3H4ZtTT/9EhIfgAIfhqIPhrW1vbPH/4ZwgNAVDtRNAg10nCAY4b0//TP9MA1NP/1wsH+Gz4a/hqf/hh+Gb4Y/hijoDiCQFK9AWI+Gpw+Gtw+GxwAYBA9A7yvdcL//hicPhjcPhmf/hhgGX4bBABvjDTB9M/0//XDf+V1NHQ0//f0fhJ+kJvE9cL//go+kJvE9cL/7ry4G/Ii9wAAAAAAAAAAAAAAAAgzxbPgc+Bz5CJijO6JM8LByPPCz8izwv/Ic8L/8lw+wBfBOMAf/hnDQL2MPhBbuMA+kDTB9M/1w3/ldTR0NP/39TRJCQkI234S8jL/3BYgED0QyPIywdxWIBA9EMiyMs/cliAQPRDJHNYgED0FiF0WIBA9BfI9ADJ+ErIz4SA9AD0AM+ByfhMIMECkzCAZN74SfpCbxPXC/8i+QC68vT4KMjPhQjODgwBVo0EUBfXhAAAAAAAAAAAAAAAAAABzxbPgc+DJs8UyXD7ADBfBF8F4wB/+GcNAD74QsjL//hDzws/+EbPCwD4SvhL+ExeIMzL/8sHye1UADztRNDT/9M/0wDU0//XCwf4bPhr+Gp/+GH4Zvhj+GIBCvSkIPShEAAA"
)

type MessageHandlerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type MessageHandler struct {
	Ctx ContractContext
}
type MessageHandlerDeployParams struct {
	ProposalCode          string
	EpochControllerPubKey string
}

func (c *MessageHandler) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{MessageHandlerTvc})
}
func (c *MessageHandler) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(MessageHandlerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *MessageHandler) Address() (string, error) {
	messageHandlerDeployParams := MessageHandlerDeployParams{
		EpochControllerPubKey: "0x7",
		ProposalCode:          "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
	}
	encodeMessage, err := c.DeployEncodeMessage(&messageHandlerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *MessageHandler) New(address string) (*MessageHandlerContract, error) {
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
	contract := MessageHandlerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *MessageHandler) DeployEncodeMessage(messageHandlerDeployParams *MessageHandlerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         MessageHandlerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_proposalCode\": \"%s\" ,\"_epochControllerPubKey\": \"%s\" }", messageHandlerDeployParams.ProposalCode, messageHandlerDeployParams.EpochControllerPubKey))
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
func (c *MessageHandler) Deploy(messageHandlerDeployParams *MessageHandlerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*MessageHandlerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(messageHandlerDeployParams)
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
	contract := MessageHandlerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *MessageHandlerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *MessageHandlerContract) ExecuteProposal(epochAddress string, chainId string, nonce string, messageType string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"epochAddress\": \"%s\" ,\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"messageType\": \"%s\" ,\"data\": \"%s\" }", epochAddress, chainId, nonce, messageType, data)
	return contract.CallContractMethod("executeProposal", input)
}
func (contract *MessageHandlerContract) ReceiveMessage(chainId string, nonce string, messageType string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"chainId\": \"%s\" ,\"nonce\": \"%s\" ,\"messageType\": \"%s\" ,\"data\": \"%s\" }", chainId, nonce, messageType, data)
	return contract.CallContractMethod("receiveMessage", input)
}
func (contract *MessageHandlerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *MessageHandlerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *MessageHandlerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
