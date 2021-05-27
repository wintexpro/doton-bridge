package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	BridgeAbi = "{\"ABI version\":2,\"data\":[],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"_relayerInitState\",\"type\":\"cell\"},{\"name\":\"_accessControllerAddress\",\"type\":\"address\"},{\"name\":\"_voteControllerAddress\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"handlerAddress\",\"type\":\"address\"},{\"name\":\"relayerPubKey\",\"type\":\"uint256\"}],\"name\":\"adminSetHandler\",\"outputs\":[]},{\"inputs\":[{\"name\":\"epochNumber\",\"type\":\"uint64\"},{\"name\":\"choice\",\"type\":\"uint8\"},{\"name\":\"chainId\",\"type\":\"uint8\"},{\"name\":\"messageType\",\"type\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint64\"},{\"name\":\"data\",\"type\":\"cell\"},{\"name\":\"relayerPubKey\",\"type\":\"uint256\"}],\"name\":\"relayerVoteForProposal\",\"outputs\":[]},{\"inputs\":[{\"name\":\"messageType\",\"type\":\"uint256\"}],\"name\":\"getHandlerAddressByMessageType\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]}],\"header\":[\"time\"]}"
	BridgeTvc = "te6ccgECFgEABPwAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCxQFBBUC1o0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIds80wABn4ECANcYIPkBWPhC+RDyqN7TPwGOHfhDIbkgnzAg+COBA+iogggbd0Cgud6TIPhj4PI02DDTHwHbPPhHbvJ8CAYBPCLQ0wP6QDD4aak4ANwhxwDcIdMfId0B2zz4R27yfAYEUCCCEBZaU6O64wIgghBLaU3QuuMCIIIQWIIfebrjAiCCEGdtZCC64wIMCwoHAlgw+EFu4wD4RvJzcfhm1PpA+kGV1NHQ+kDf0fgAIvhsIfhqIPhrXwPbPH/4ZwgPAXDtRNAg10nCAY4r0//TP9MA+kD6QNT0BNMH0wfXCwf4cPhv+G74bfhs+Gv4an/4Yfhm+GP4Yo6A4gkB/vQFjQhgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE+GqNCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAT4a4j4bG34bXD4bnD4b3D4cHABgED0DvK91wv/+GJw+GNw+GZ/+GGAZfhugGb4b4Bn+HAVA6Iw+EFu4wDT//pBldTR0PpA39cN/5XU0dDT/9/RIPhMIds8+E4gwQKTMIBk3vhJ+kJvE9cL/yL5ALry9PhNJQElWYEBAPQW+G0wMF8D2zx/+GcTEQ8CzjD4QW7jANP/0SD4TYEBAPQOjiSNCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATfMSHA/44iI9DTAfpAMDHIz4cgzoBgz0DPgc+Bz5MtpTdCIc8WyXD7AN4w4wB/+GcTDwTEMPhBbuMA0z/TB9MH0//TP9TXDf+V1NHQ0//f0SD4TCHbPPhOIMECkzCAZN74SfpCbxPXC/8i+QC68vT4TyDBApMwgGTe2zyCEBfXhAC+8vT4UCDBApMwgGTeJvhNgQEA9A4TERANAf6OJI0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABN+NCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATHBbPy9PhLf8jPhYDKAHPPQM6NBFBTck4AAAAAAAAAAAAAAAAAAc8Wz4HPg8jPkRyAC9IqDgG2zws/+EnPFinPCwcozwsHJ88L/yf4TYEBAPQOjiSNCGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATfzxbIJ88LPybPFM3NyXH7ADAwXwfjAH/4Zw8AWvhCyMv/+EPPCz/4Rs8LAPhK+Ev4TPhN+E74T/hQXmDOzsz0AMsHywfLB8ntVAAYcGim+2CVaKb+YDHfAe4h0Mgh0wAzwACTcc9AmnHPQSHTHzPPCx/iIdMAM8AAk3HPQJpxz0Eh0wEzzwsB4iHTADPAAJNxz0CYcc9BIdQzzxTiIdMAM8MB8ndxz0HII88L/yLUNND0BAEicCKAQPRDMSDI9AAgySXMNV8EIdMAM8AAk3HPQBIAHJhxz0Eh1DPPFOIgyWxBAFztRNDT/9M/0wD6QPpA1PQE0wfTB9cLB/hw+G/4bvht+Gz4a/hqf/hh+Gb4Y/hiAQr0pCD0oRUAAA=="
)

type BridgeContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type Bridge struct {
	Ctx ContractContext
}
type BridgeDeployParams struct {
	RelayerInitState        string
	AccessControllerAddress string
	VoteControllerAddress   string
}

func (c *Bridge) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{BridgeTvc})
}
func (c *Bridge) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(BridgeAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *Bridge) Address() (string, error) {
	bridgeDeployParams := BridgeDeployParams{
		AccessControllerAddress: "0:0000000000000000000000000000000000000000000000000000000000000000",
		RelayerInitState:        "te6ccgEBAQEAOgAAb8AP8AxBKu/bcRJNoYhaIb6rV2wlmZXVB48QbP/KW5bki0ICXcMBftcYAAAAAAAAAB2LXmIPSAAE",
		VoteControllerAddress:   "0:0000000000000000000000000000000000000000000000000000000000000000",
	}
	encodeMessage, err := c.DeployEncodeMessage(&bridgeDeployParams)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *Bridge) New(address string) (*BridgeContract, error) {
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
	contract := BridgeContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *Bridge) DeployEncodeMessage(bridgeDeployParams *BridgeDeployParams) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	deploySet := client.DeploySet{
		Tvc:         BridgeTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"_relayerInitState\": \"%s\" ,\"_accessControllerAddress\": \"%s\" ,\"_voteControllerAddress\": \"%s\" }", bridgeDeployParams.RelayerInitState, bridgeDeployParams.AccessControllerAddress, bridgeDeployParams.VoteControllerAddress))
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
func (c *Bridge) Deploy(bridgeDeployParams *BridgeDeployParams, messageCallback func(event *client.ProcessingEvent)) (*BridgeContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(bridgeDeployParams)
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
	contract := BridgeContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *BridgeContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *BridgeContract) AdminSetHandler(messageType string, handlerAddress string, relayerPubKey string) *ContractMethod {
	input := fmt.Sprintf("{\"messageType\": \"%s\" ,\"handlerAddress\": \"%s\" ,\"relayerPubKey\": \"%s\" }", messageType, handlerAddress, relayerPubKey)
	return contract.CallContractMethod("adminSetHandler", input)
}
func (contract *BridgeContract) RelayerVoteForProposal(epochNumber string, choice string, chainId string, messageType string, nonce string, data string, relayerPubKey string) *ContractMethod {
	input := fmt.Sprintf("{\"epochNumber\": \"%s\" ,\"choice\": \"%s\" ,\"chainId\": \"%s\" ,\"messageType\": \"%s\" ,\"nonce\": \"%s\" ,\"data\": \"%s\" ,\"relayerPubKey\": \"%s\" }", epochNumber, choice, chainId, messageType, nonce, data, relayerPubKey)
	return contract.CallContractMethod("relayerVoteForProposal", input)
}
func (contract *BridgeContract) GetHandlerAddressByMessageType(messageType string) *ContractMethod {
	input := fmt.Sprintf("{\"messageType\": \"%s\" }", messageType)
	return contract.CallContractMethod("getHandlerAddressByMessageType", input)
}
func (contract *BridgeContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *BridgeContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *BridgeContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
