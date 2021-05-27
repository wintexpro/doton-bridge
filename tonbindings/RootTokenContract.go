package tonbindings

import (
	"encoding/json"
	"fmt"
	client "github.com/radianceteam/ton-client-go/client"
	null "github.com/volatiletech/null"
)

const (
	RootTokenContractAbi = "{\"ABI version\":2,\"data\":[{\"key\":1,\"name\":\"_randomNonce\",\"type\":\"uint256\"},{\"key\":2,\"name\":\"name\",\"type\":\"bytes\"},{\"key\":3,\"name\":\"symbol\",\"type\":\"bytes\"},{\"key\":4,\"name\":\"decimals\",\"type\":\"uint8\"},{\"key\":5,\"name\":\"wallet_code\",\"type\":\"cell\"}],\"events\":[],\"functions\":[{\"inputs\":[{\"name\":\"root_public_key_\",\"type\":\"uint256\"},{\"name\":\"root_owner_address_\",\"type\":\"address\"}],\"name\":\"constructor\",\"outputs\":[]},{\"inputs\":[],\"name\":\"getDetails\",\"outputs\":[{\"components\":[{\"name\":\"name\",\"type\":\"bytes\"},{\"name\":\"symbol\",\"type\":\"bytes\"},{\"name\":\"decimals\",\"type\":\"uint8\"},{\"name\":\"wallet_code\",\"type\":\"cell\"},{\"name\":\"root_public_key\",\"type\":\"uint256\"},{\"name\":\"root_owner_address\",\"type\":\"address\"},{\"name\":\"total_supply\",\"type\":\"uint128\"},{\"name\":\"start_gas_balance\",\"type\":\"uint128\"},{\"name\":\"paused\",\"type\":\"bool\"}],\"name\":\"value0\",\"type\":\"tuple\"}]},{\"inputs\":[{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"}],\"name\":\"getWalletAddress\",\"outputs\":[{\"name\":\"value0\",\"type\":\"address\"}]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"},{\"name\":\"gas_back_address\",\"type\":\"address\"}],\"name\":\"deployWallet\",\"outputs\":[]},{\"inputs\":[{\"name\":\"grams\",\"type\":\"uint128\"},{\"name\":\"wallet_public_key_\",\"type\":\"uint256\"},{\"name\":\"owner_address_\",\"type\":\"address\"},{\"name\":\"gas_back_address\",\"type\":\"address\"}],\"name\":\"deployEmptyWallet\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"proxyBurn\",\"outputs\":[]},{\"inputs\":[{\"name\":\"tokens\",\"type\":\"uint128\"},{\"name\":\"sender_public_key\",\"type\":\"uint256\"},{\"name\":\"sender_address\",\"type\":\"address\"},{\"name\":\"callback_address\",\"type\":\"address\"},{\"name\":\"callback_payload\",\"type\":\"cell\"}],\"name\":\"tokensBurned\",\"outputs\":[]},{\"inputs\":[],\"name\":\"withdrawExtraGas\",\"outputs\":[]},{\"inputs\":[{\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"setPaused\",\"outputs\":[]},{\"inputs\":[{\"name\":\"callback_id\",\"type\":\"uint64\"},{\"name\":\"callback_addr\",\"type\":\"address\"}],\"name\":\"sendPausedCallbackTo\",\"outputs\":[]},{\"inputs\":[{\"name\":\"root_public_key_\",\"type\":\"uint256\"},{\"name\":\"root_owner_address_\",\"type\":\"address\"}],\"name\":\"transferOwner\",\"outputs\":[]},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"name\",\"type\":\"bytes\"}]},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"symbol\",\"type\":\"bytes\"}]},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8\"}]},{\"inputs\":[],\"name\":\"wallet_code\",\"outputs\":[{\"name\":\"wallet_code\",\"type\":\"cell\"}]},{\"inputs\":[],\"name\":\"total_supply\",\"outputs\":[{\"name\":\"total_supply\",\"type\":\"uint128\"}]},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"paused\",\"type\":\"bool\"}]}],\"header\":[\"pubkey\",\"time\",\"expire\"]}"
	RootTokenContractTvc = "te6ccgECPAEADigAAgE0AwEBAcACAEPQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgBCj/AIrtUyDjAyDA/+MCIMD+4wLyCzoHBDsBAAUC/I0IYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPhpIds80wABjh2BAgDXGCD5AQHTAAGU0/8DAZMC+ELiIPhl+RDyqJXTAAHyeuLTPwGOHfhDIbkgnzAg+COBA+iogggbd0Cgud6TIPhj4PI02DDTHwH4I7zyuTMGAhbTHwHbPPhHbo6A3goIA0Ii0NMD+kAw+GmpOACOgOAhxwDcIdMfId0B2zz4R26OgN43CggBBlvbPAkCDvhBbuMA2zw5OARYIIIQFZ7bPbuOgOAgghA4KCYau46A4CCCEGpjV1y7joDgIIIQeYWz9LuOgOAqHBILAzwgghByPcTOuuMCIIIQdkBjNrrjAiCCEHmFs/S64wIQDQwBVNs8+EvIi9wAAAAAAAAAAAAAAAAgzxbPgc+Bz5PmFs/SIc8UyXD7AH/4ZzkD/jD4QW7jANcNf5XU0dDTf9/XDf+V1NHQ0//f+kGV1NHQ+kDf+kGV1NHQ+kDf0fhWIMECkzCAZN4i+kJvE9cL/8MAIJQwI8AA3iCOEjAi+kJvE9cL/8AAIJQwI8MA3t/y9PhS+CdvENs8obV/tgly+wJtI8jL/3BYgED0Q/gocVg5MQ4B+oBA9Bb4TnJYgED0FyPIy/9zWIBA9EMidFiAQPQWyPQAyfhOyM+EgPQA9ADPgckg+QDIz4oAQMv/ydAlIcjPhYjOAfoCgGnPQM+Dz4MizxTPgc+RotV8/slw+wAxMCD6Qm8T1wv/wwCOFCDIz4WIzoBtz0DPgc+ByYEAgPsADwE+jhX4ScjPhYjOgG3PQM+Bz4HJgQCA+wDiXwTbPH/4ZzgD5jD4QW7jANcNf5XU0dDTf9/6QZXU0dD6QN/R+FMgwQKTMIBk3vhR+kJvE9cL/8MAIJcw+FH4SccF3iCOFDD4UMMAIJww+FD4RSBukjBw3rre3/L0+FH6Qm8T1wv/wACS+ACOgOIh+E8BoLV/+G8gyM+FiM45JhEBnI0EDmJaAAAAAAAAAAAAAAAAAAHPFs+Bz4HPkCz/PV4izwt/yXD7APhR+kJvE9cL/44V+FHIz4WIzoBtz0DPgc+ByYEAgPsA3lvbPH/4ZzgDQCCCEEWzvf27joDgIIIQaVUvMLuOgOAgghBqY1dcuuMCGRUTA/4w+EFu4wDXDX+V1NHQ03/f+kGV1NHQ+kDf+kGV1NHQ+kDf1NH4UyDBApMwgGTe+FH6Qm8T1wv/wwAglzD4UfhJxwXe8vT4J28Q2zyhtX9y+wJwI20iyMv/cFiAQPRD+ChxWIBA9Bb4TnJYgED0FyLIy/9zWIBA9EMhdFiAQPQWOTEUAZjI9ADJ+E7Iz4SA9AD0AM+BySD5AMjPigBAy//J0DFsISDIz4WIzoBtz0DPgc+Bz5BFl3vWJc8LfyPPFiLPFMmBAID7ADBfBNs8f/hnOAIoIIIQTKuUqLrjAiCCEGlVLzC64wIYFgL+MPhBbuMA1w3/ldTR0NP/3/pBldTR0PpA39H4ViDBApMwgGTeIfpCbxPXC//DACCUMCLAAN4gjhIwIfpCbxPXC//AACCUMCLDAN7f8vQhIW0iyMv/cFiAQPRD+ChxWIBA9Bb4TnJYgED0FyLIy/9zWIBA9EMhdFiAQPQWyPQAyTkXAZz4TsjPhID0APQAz4HJIPkAyM+KAEDL/8nQMWwhIDFsISHA/44iI9DTAfpAMDHIz4cgzoBgz0DPgc+Bz5OlVLzCIc8WyXD7AN4w2zx/+Gc4AVbbPPhPyIvcAAAAAAAAAAAAAAAAIM8Wz4HPgc+TMq5SoiHPC3/JcPsAf/hnOQIoIIIQPrN0nrrjAiCCEEWzvf264wIbGgFU2zz4TMiL3AAAAAAAAAAAAAAAACDPFs+Bz4HPkxbO9/YhzxTJcPsAf/hnOQLeMPhBbuMA0fhL+Ez4TfhO+FD4UfhP+FL4WG8JIcD/jkgj0NMB+kAwMcjPhyDOgGDPQM+Bz4PIz5L6zdJ6Im8pVQgpzxQozxQnzwsHJs8UJc8L/yTPFiPPC38izwt/Ic8KAGyRzclw+wDeMOMAf/hnOTgDQCCCEDCNZtG7joDgIIIQMlARAbuOgOAgghA4KCYauuMCJx8dAt4w+EFu4wDXDf+V1NHQ0//f+kGV1NHQ+kDf0fhTIMECkzCAZN74UfpCbxPXC//DACCXMPhR+EnHBd4gjhQw+FDDACCcMPhQ+EUgbpIwcN663t/y9PhWIMECkzCAZN4iwwAgmzAh+kJvE9cL/8AA3iA5HgFKjhIwIsAAIJswIfpCbxPXC//DAN7f8vT4ACH4cCD4cVvbPH/4ZzgCKCCCEDHDJJm64wIgghAyUBEBuuMCIiAC/jD4QW7jANcNf5XU0dDTf9/XDf+V1NHQ0//f+kGV1NHQ+kDf+kGV1NHQ+kDf1NH4VyDBApMwgGTe+Fiz8vQjI20iyMv/cFiAQPRD+ChxWIBA9Bb4TnJYgED0FyLIy/9zWIBA9EMhdFiAQPQWyPQAyfhOyM+EgPQA9ADPgckg+QA5IQLOyM+KAEDL/8nQMWwh+FUgwQKTMIBk3vhJIscF8vT4J28Q2zyhtX9y+wIl+E8BobV/+G8iyM+FiM6Abc9Az4HPg8jPkcd0ndInzwt/I88UJs8L/yXPFiLPFs3JgQCA+wAwXwXbPH/4ZzE4Av4w+EFu4wDXDX+V1NHQ03/f1w1/ldTR0NN/39cN/5XU0dDT/9/6QZXU0dD6QN/6QZXU0dD6QN/R+FMgwQKTMIBk3vhR+kJvE9cL/8MAIJcw+FH4SccF3iCOFDD4UMMAIJww+FD4RSBukjBw3rre3/L0JHC+8uBk+FYgwQKTMIBkOSMC+N4i+kJvE9cL/8MAIJQwI8AA3iCOEjAi+kJvE9cL/8AAIJQwI8MA3t/y9PhR+kJvE9cL/8AAkvgAjoDibSPIy/9wWIBA9EP4KHFYgED0FvhOcliAQPQXI8jL/3NYgED0QyJ0WIBA9BbI9ADJ+E7Iz4SA9AD0AM+BySD5AMgmJAHSz4oAQMv/ydAlIcjPhYjOAfoCgGnPQM+Dz4MizxTPgc+RotV8/slw+wAxIMjPhYjOjQQOYloAAAAAAAAAAAAAAAAAAc8Wz4HPgc+QLP89XibPC3/JcPsAJfhPAaC1f/hv+FH6Qm8T1wv/JQGGjjgh+kJvE9cL/8MAjhQhyM+FiM6Abc9Az4HPgcmBAID7AI4V+EnIz4WIzoBtz0DPgc+ByYEAgPsA4t4wXwXbPH/4ZzgBIPhS+CdvENs8obV/tgly+wIxAiggghAilxn3uuMCIIIQMI1m0brjAikoAVbbPPhYyIvcAAAAAAAAAAAAAAAAIM8Wz4HPgc+SwjWbRiHPCgDJcPsAf/hnOQFU2zz4TsiL3AAAAAAAAAAAAAAAACDPFs+Bz4HPkopcZ94hzxTJcPsAf/hnOQM+IIIJ9Rpmu46A4CCCEBUAWwe7joDgIIIQFZ7bPbrjAi8sKwKOMPhBbuMA0fhTIMECkzCAZN74UfpCbxPXC//DACCXMPhR+EnHBd7y9PhScvsC+FHIz4WIzoBtz0DPgc+ByYEAgPsA2zx/+Gc5OAIoIIIQDVr8crrjAiCCEBUAWwe64wIuLQFW2zz4TciL3AAAAAAAAAAAAAAAACDPFs+Bz4HPklQBbB4hzwsHyXD7AH/4ZzkCmjD4QW7jANIA0fhTIMECkzCAZN74UfpCbxPXC//DACCXMPhR+EnHBd4gjhQw+FDDACCcMPhQ+EUgbpIwcN663t/y9PgAIPh4MNs8f/hnOTgCJCCCCdU9HbrjAiCCCfUaZrrjAjIwA5Yw+EFu4wDTP/pBldTR0PpA39H4UvgnbxDbPKG1f7YJcvsCIMjPhYjOgG3PQM+Bz4HPkc4bw6Iizws/+FjPCgDJgQCA+wBb2zx/+Gc5MTgAGHBopvtglWim/mAx3wLcMPhBbuMA+Ebyc3H4ZtcN/5XU0dDT/9/6QZXU0dD6QN/R+FYgwQKTMIBk3iLDACCbMCH6Qm8T1wv/wADeII4SMCLAACCbMCH6Qm8T1wv/wwDe3/L0+AAh+HAg+HFw+G9w+Hj4J28Q+HJb2zx/+GczOAGw7UTQINdJwgGOS9P/0z/TANX6QNcLf/hy+HHT/9TU0wfU03/T/9MH0wfTB9MH0wfXCgD4ePh3+Hb4dfh0+HP4cPhv+G74bfhs+Gv4an/4Yfhm+GP4Yo6A4jQE/vQFcSGAQPQOk9cL/5Fw4vhqciGAQPQPjoDf+GtzIYBA9A+OgN/4bHQhgED0DpPXCweRcOL4bXUhgED0D46A3/hucPhvcPhwjQhgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE+HFw+HJw+HNw+HRw+HVw+HZw+Hc2NjY1AFpw+HhwAYBA9A7yvdcL//hicPhjcPhmf/hhgGT4c4Bl+HSAZ/h1gGr4doBr+HcBAog7AlYh1h8x+EFu4wD4ACDTHzIgghALP89Xup4h038zIPhPAaG1f/hvMN4wMNs8OTgAoPhCyMv/+EPPCz/4Rs8LAMj4UfhSAs7Lf/hK+Ev4TPhN+E74T/hQ+FP4VPhV+Fb4V/hYXtDPEcv/zMzLB8zLf8v/ywfLB8sHywfLB8oAye1UAJztRNDT/9M/0wDV+kDXC3/4cvhx0//U1NMH1NN/0//TB9MH0wfTB9MH1woA+Hj4d/h2+HX4dPhz+HD4b/hu+G34bPhr+Gp/+GH4Zvhj+GIBCvSkIPShOwAA"
)

type RootTokenContractContract struct {
	Abi     client.Abi
	Address string
	Ctx     ContractContext
}
type RootTokenContract struct {
	Ctx ContractContext
}
type RootTokenContractDeployParams struct {
	Rootpublickey    string
	Rootowneraddress string
}
type RootTokenContractInitVars struct {
	RandomNonce string
	Name        string
	Symbol      string
	Decimals    string
	Walletcode  string
}

func (c *RootTokenContract) Code() (*client.ResultOfGetCodeFromTvc, error) {
	return c.Ctx.Conn.BocGetCodeFromTvc(&client.ParamsOfGetCodeFromTvc{RootTokenContractTvc})
}
func (c *RootTokenContract) Abi() (*client.Abi, error) {
	abi := client.Abi{Type: client.ContractAbiType}
	if err := json.Unmarshal([]byte(RootTokenContractAbi), &abi.Value); err != nil {
		return nil, err
	}
	return &abi, nil
}
func (c *RootTokenContract) Address(rootTokenContractInitVars *RootTokenContractInitVars) (string, error) {
	rootTokenContractDeployParams := RootTokenContractDeployParams{
		Rootowneraddress: "0:0000000000000000000000000000000000000000000000000000000000000000",
		Rootpublickey:    "0x7",
	}
	encodeMessage, err := c.DeployEncodeMessage(&rootTokenContractDeployParams, rootTokenContractInitVars)
	if err != nil {
		return "", err
	}
	return encodeMessage.Address, nil
}
func (c *RootTokenContract) New(address string, rootTokenContractInitVars *RootTokenContractInitVars) (*RootTokenContractContract, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	if address == "" {
		address, err = c.Address(rootTokenContractInitVars)
		if err != nil {
			return nil, err
		}
	}
	contract := RootTokenContractContract{
		Abi:     *abi,
		Address: address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (c *RootTokenContract) DeployEncodeMessage(rootTokenContractDeployParams *RootTokenContractDeployParams, rootTokenContractInitVars *RootTokenContractInitVars) (*client.ResultOfEncodeMessage, error) {
	abi, err := c.Abi()
	if err != nil {
		return nil, err
	}
	initialVars := json.RawMessage(fmt.Sprintf("{\"_randomNonce\": \"%s\" ,\"name\": \"%s\" ,\"symbol\": \"%s\" ,\"decimals\": \"%s\" ,\"wallet_code\": \"%s\" }", rootTokenContractInitVars.RandomNonce, rootTokenContractInitVars.Name, rootTokenContractInitVars.Symbol, rootTokenContractInitVars.Decimals, rootTokenContractInitVars.Walletcode))
	deploySet := client.DeploySet{
		InitialData: initialVars,
		Tvc:         RootTokenContractTvc,
		WorkchainID: c.Ctx.WorkchainID,
	}
	params := json.RawMessage(fmt.Sprintf("{\"root_public_key_\": \"%s\" ,\"root_owner_address_\": \"%s\" }", rootTokenContractDeployParams.Rootpublickey, rootTokenContractDeployParams.Rootowneraddress))
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
func (c *RootTokenContract) Deploy(rootTokenContractDeployParams *RootTokenContractDeployParams, rootTokenContractInitVars *RootTokenContractInitVars, messageCallback func(event *client.ProcessingEvent)) (*RootTokenContractContract, error) {
	abi, err := c.Abi()
	encodeMessage, err := c.DeployEncodeMessage(rootTokenContractDeployParams, rootTokenContractInitVars)
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
	contract := RootTokenContractContract{
		Abi:     *abi,
		Address: encodeMessage.Address,
		Ctx:     c.Ctx,
	}
	return &contract, nil
}
func (contract *RootTokenContractContract) CallContractMethod(methodName string, input string) *ContractMethod {
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
func (contract *RootTokenContractContract) GetDetails() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("getDetails", input)
}
func (contract *RootTokenContractContract) GetWalletAddress(wallet_public_key_ string, owner_address_ string) *ContractMethod {
	input := fmt.Sprintf("{\"wallet_public_key_\": \"%s\" ,\"owner_address_\": \"%s\" }", wallet_public_key_, owner_address_)
	return contract.CallContractMethod("getWalletAddress", input)
}
func (contract *RootTokenContractContract) DeployWallet(tokens string, grams string, wallet_public_key_ string, owner_address_ string, gas_back_address string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"grams\": \"%s\" ,\"wallet_public_key_\": \"%s\" ,\"owner_address_\": \"%s\" ,\"gas_back_address\": \"%s\" }", tokens, grams, wallet_public_key_, owner_address_, gas_back_address)
	return contract.CallContractMethod("deployWallet", input)
}
func (contract *RootTokenContractContract) DeployEmptyWallet(grams string, wallet_public_key_ string, owner_address_ string, gas_back_address string) *ContractMethod {
	input := fmt.Sprintf("{\"grams\": \"%s\" ,\"wallet_public_key_\": \"%s\" ,\"owner_address_\": \"%s\" ,\"gas_back_address\": \"%s\" }", grams, wallet_public_key_, owner_address_, gas_back_address)
	return contract.CallContractMethod("deployEmptyWallet", input)
}
func (contract *RootTokenContractContract) Mint(tokens string, to string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"to\": \"%s\" }", tokens, to)
	return contract.CallContractMethod("mint", input)
}
func (contract *RootTokenContractContract) ProxyBurn(tokens string, sender_address string, callback_address string, callback_payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"sender_address\": \"%s\" ,\"callback_address\": \"%s\" ,\"callback_payload\": \"%s\" }", tokens, sender_address, callback_address, callback_payload)
	return contract.CallContractMethod("proxyBurn", input)
}
func (contract *RootTokenContractContract) TokensBurned(tokens string, sender_public_key string, sender_address string, callback_address string, callback_payload string) *ContractMethod {
	input := fmt.Sprintf("{\"tokens\": \"%s\" ,\"sender_public_key\": \"%s\" ,\"sender_address\": \"%s\" ,\"callback_address\": \"%s\" ,\"callback_payload\": \"%s\" }", tokens, sender_public_key, sender_address, callback_address, callback_payload)
	return contract.CallContractMethod("tokensBurned", input)
}
func (contract *RootTokenContractContract) WithdrawExtraGas() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("withdrawExtraGas", input)
}
func (contract *RootTokenContractContract) SetPaused(value string) *ContractMethod {
	input := fmt.Sprintf("{\"value\": \"%s\" }", value)
	return contract.CallContractMethod("setPaused", input)
}
func (contract *RootTokenContractContract) SendPausedCallbackTo(callback_id string, callback_addr string) *ContractMethod {
	input := fmt.Sprintf("{\"callback_id\": \"%s\" ,\"callback_addr\": \"%s\" }", callback_id, callback_addr)
	return contract.CallContractMethod("sendPausedCallbackTo", input)
}
func (contract *RootTokenContractContract) TransferOwner(root_public_key_ string, root_owner_address_ string) *ContractMethod {
	input := fmt.Sprintf("{\"root_public_key_\": \"%s\" ,\"root_owner_address_\": \"%s\" }", root_public_key_, root_owner_address_)
	return contract.CallContractMethod("transferOwner", input)
}
func (contract *RootTokenContractContract) Name() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("name", input)
}
func (contract *RootTokenContractContract) Symbol() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("symbol", input)
}
func (contract *RootTokenContractContract) Decimals() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("decimals", input)
}
func (contract *RootTokenContractContract) Wallet_code() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("wallet_code", input)
}
func (contract *RootTokenContractContract) Total_supply() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("total_supply", input)
}
func (contract *RootTokenContractContract) Paused() *ContractMethod {
	input := fmt.Sprintf("{}")
	return contract.CallContractMethod("paused", input)
}
func (contract *RootTokenContractContract) abiEncodeMessage(functionName string, input string) (*client.ResultOfEncodeMessage, error) {
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
func (contract *RootTokenContractContract) send(functionName string, input string, messageCallback func(event *client.ProcessingEvent)) (string, error) {
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
func (contract *RootTokenContractContract) call(functionName string, input string) (*client.DecodedOutput, error) {
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
