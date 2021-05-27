package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	EpochControllerAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_epochCode\",\"type\":\"cell\"},{\"name\":\"_proposalCode\",\"type\":\"cell\"},{\"name\":\"_deployInitialValue\",\"type\":\"uint128\"},{\"name\":\"_publicKey\",\"type\":\"uint256\"},{\"name\":\"_proposalVotersAmount\",\"type\":\"uint256\"},{\"name\":\"_bridgeAddress\",\"type\":\"address\"},{\"name\":\"_firstEraDuration\",\"type\":\"uint32\"},{\"name\":\"_secondEraDuration\",\"type\":\"uint32\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"number\",\"type\":\"uint64\"}],\"name\":\"getEpochAddress\",\"outputs\":[{\"name\":\"epoch\",\"type\":\"address\"}]},{\"inputs\":[{\"name\":\"epochNumber\",\"type\":\"uint64\"},{\"name\":\"voter\",\"type\":\"address\"},{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"}],\"name\":\"voteByBridge\",\"outputs\":[]},{\"inputs\":[{\"name\":\"newPublicrandomness\",\"type\":\"uint256\"}],\"name\":\"newEpoch\",\"outputs\":[]},{\"inputs\":[],\"name\":\"firstEraDuration\",\"outputs\":[{\"name\":\"firstEraDuration\",\"type\":\"uint32\"}]},{\"inputs\":[],\"name\":\"secondEraDuration\",\"outputs\":[{\"name\":\"secondEraDuration\",\"type\":\"uint32\"}]},{\"inputs\":[],\"name\":\"publicRandomness\",\"outputs\":[{\"name\":\"publicRandomness\",\"type\":\"uint256\"}]},{\"inputs\":[],\"name\":\"currentEpochNumber\",\"outputs\":[{\"name\":\"currentEpochNumber\",\"type\":\"uint64\"}]}],\"header\":[\"time\"]}"
	EpochControllerTvc = "te6ccgECHQEABjkAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCxsFBBwC1o0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIds80wABn4ECANcYIPkBWPhC+RDyqN7TPwGOHfhDIbkgnzAg+COBA+iogggbd0Cgud6TIPhj4PI02DDTHwHbPPhHbvJ8FQYBPCLQ0wP6QDD4aak4ANwhxwDcIdMfId0B2zz4R27yfAYEWCCCEBqTG3K7joDgIIIQKk2HdruOgOAgghBHIAL0u46A4CCCEGD6RBG7joDgFw4KBwIoIIIQUzzfiLrjAiCCEGD6RBG64wIJCANyMPhBbuMA0z/R2zwhwP+OIiPQ0wH6QDAxyM+HIM6AYM9Az4HPgc+Tg+kQRiHPFslw+wDeMOMAf/hnGhASAVbbPPhTyIvcAAAAAAAAAAAAAAAAIM8Wz4HPgc+TTPN+IiHPCx/JcPsAf/hnGgIoIIIQP0yu6rrjAiCCEEcgAvS64wINCwP+MPhBbuMA0z/6QNMH0wfXDf+V1NHQ0//f+kGV1NHQ+kDf1w0/ldTR0NM/39TR+EogwQKTMIBk3vhJ+E3HBfL0+EsgwQKTMIBk3ibAACCUMCbAAd/y9PgAJ9s8f8jPhYDKAHPPQM6NBFBTck4AAAAAAAAAAAAAAAAAAc8Wz4HPgxoQDAFYyM+RtNaoqijPFifPCwcmzwsHJc8L/yTPFiPPCz8izxTNyXH7AF8I2zx/+GcSAVbbPPhUyIvcAAAAAAAAAAAAAAAAIM8Wz4HPgc+S/TK7qiHPCx/JcPsAf/hnGgIoIIIQI3wAyLrjAiCCECpNh3a64wIRDwRyMPhBbuMA0//R+EwgwQKTMIBk3vhJ+FbbPMcF8vT4ACD4dfhWcaC1PyHbPDD4VqS1P/h2MNs8f/hnGhATEgDIjQhgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEbfhCyMv/cFiAQPRDIsjLP3FYgED0Q/gocliAQPQWyPQAyfhOyM+EgPQA9ADPgckg+QDIz4oAQMv/ydBsEgEwMQPOMPhBbuMA+Ebyc3H4ZtTU03/T/9cN/5XU0dDT/9/6QZXU0dD6QN/XDR+V1NHQ0x/f1w0fldTR0NMf39El+HAk+HEi+G0n+G4m+G8j+HIh+HMg+HT4EPh1cfh2cfhV2zwwXwjbPH/4ZxUTEgCU+ELIy//4Q88LP/hGzwsAyPhR+FL4VV4gy//L/8v/+Er4S/hM+E34TvhP+FD4U/hU+FZeoM8RywfLB8sHzszMy3/LH8sfyz/J7VQB9I0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPgAbfhCyMv/cFiAQPRDI8jLP3FYgED0Q/gocliAQPQWyPQAyfhOyM+EgPQA9ADPgckg+QDIz4oAQMv/ydD4UCHIz4WIzgH6AoBpz0DPg8+DIs8Uz4PIFABkz5BJWp6a+E/PFPhQzwt/+ELPC//4Us8L//hTzwsf+FTPCx8lzwv/zclw+wAxMSAxbCEBou1E0CDXScIBjkTT/9M/0wDV0//T/9cL//h1+HL4cdMH0wfTB/pA1NTTf9Mf0x/XCz/4dvh0+HP4cPhv+G74bfhs+Gv4an/4Yfhm+GP4Yo6A4hYC3PQFcPhqcPhrcPhsjQhgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE+G2I+G6I+G9w+HBw+HFw+HJw+HNw+HRw+HVw+HZwAYBA9A7yvdcL//hicPhjcPhmf/hhgG/4aoBw+GuAcfhsHBwCKCCCEBnJQrq64wIgghAakxtyuuMCGRgBVts8+FbIi9wAAAAAAAAAAAAAAAAgzxbPgc+Bz5JqTG3KIc8LP8lw+wB/+GcaAVbbPPhVyIvcAAAAAAAAAAAAAAAAIM8Wz4HPgc+SZyUK6iHPC//JcPsAf/hnGgCO7UTQ0//TP9MA1dP/0//XC//4dfhy+HHTB9MH0wf6QNTU03/TH9Mf1ws/+Hb4dPhz+HD4b/hu+G34bPhr+Gp/+GH4Zvhj+GIBCvSkIPShHAAA"
)

type EpochControllerContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type EpochController struct {
	Ctx ContractContext
}
type EpochControllerDeployParams struct {
	EpochCode            string
	ProposalCode         string
	DeployInitialValue   string
	PublicKey            string
	ProposalVotersAmount string
	BridgeAddress        string
	FirstEraDuration     string
	SecondEraDuration    string
}

func (c *EpochController) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{EpochControllerTvc})
}
func (c *EpochController) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(EpochControllerAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *EpochController) Address() (string, error) {
	epochControllerDeployParams := EpochControllerDeployParams{
		BridgeAddress:        "0:0000000000000000000000000000000000000000000000000000000000000000",
		DeployInitialValue:   "0x6",
		EpochCode:            "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		FirstEraDuration:     "0x4",
		ProposalCode:         "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		ProposalVotersAmount: "0x7",
		PublicKey:            "0x7",
		SecondEraDuration:    "0x4",
	}
	encodeMessage, err := c.DeployEncodeMessage(&epochControllerDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *EpochController) New(address string) (*EpochControllerContract, error) {
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
	contract := EpochControllerContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *EpochController) DeployEncodeMessage(epochControllerDeployParams *EpochControllerDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         EpochControllerTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_epochCode\": \"%s\" ,\"_proposalCode\": \"%s\" ,\"_deployInitialValue\": \"%s\" ,\"_publicKey\": \"%s\" ,\"_proposalVotersAmount\": \"%s\" ,\"_bridgeAddress\": \"%s\" ,\"_firstEraDuration\": \"%s\" ,\"_secondEraDuration\": \"%s\" }", epochControllerDeployParams.EpochCode, epochControllerDeployParams.ProposalCode, epochControllerDeployParams.DeployInitialValue, epochControllerDeployParams.PublicKey, epochControllerDeployParams.ProposalVotersAmount, epochControllerDeployParams.BridgeAddress, epochControllerDeployParams.FirstEraDuration, epochControllerDeployParams.SecondEraDuration))
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
func (c *EpochController) Deploy(epochControllerDeployParams *EpochControllerDeployParams, messageCallback func(event *client.ProcessingEvent)) (*EpochControllerContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(epochControllerDeployParams)
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
	contract := EpochControllerContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *EpochControllerContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *EpochControllerContract) GetEpochAddress(number string) *ContractMethod {
	input := fmt.Sprintf("{\"number\": \"%s\" }", number)
	return contract.CallContractMethod("getEpochAddress", input)
}
func (contract *EpochControllerContract) VoteByBridge(epochNumber string, voter string, choice string, chainId string, messageType string, handlerAddress string, nonce string, data string) *ContractMethod {
	input := fmt.Sprintf("{\"epochNumber\": \"%s\" ,\"voter\": \"%s\" ,\"choice\": \"%s\" ,\"chainId\": \"%s\" ,\"messageType\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" }", epochNumber, voter, choice, chainId, messageType, handlerAddress, nonce, data)
	return contract.CallContractMethod("voteByBridge", input)
}
func (contract *EpochControllerContract) NewEpoch(newPublicrandomness string) *ContractMethod {
	input := fmt.Sprintf("{\"newPublicrandomness\": \"%s\" }", newPublicrandomness)
	return contract.CallContractMethod("newEpoch", input)
}
func (contract *EpochControllerContract) FirstEraDuration() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("firstEraDuration", input)
}
func (contract *EpochControllerContract) SecondEraDuration() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("secondEraDuration", input)
}
func (contract *EpochControllerContract) PublicRandomness() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("publicRandomness", input)
}
func (contract *EpochControllerContract) CurrentEpochNumber() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("currentEpochNumber", input)
}
func (contract *EpochControllerContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *EpochControllerContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *EpochControllerContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
