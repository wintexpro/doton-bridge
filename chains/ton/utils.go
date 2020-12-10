// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"strconv"

	connection "github.com/ChainSafe/ChainBridge/connections/ton"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
)

type Message = map[string]interface{}

func LoadAbi(path, name string) client.Abi {
	content, err := ioutil.ReadFile(path + "/" + name + ".abi.json")
	if err != nil {
		panic(err)
	}
	abi := client.Abi{Type: client.ContractAbiType}
	if err = json.Unmarshal(content, &abi.Value); err != nil {
		panic(err)
	}

	return abi
}

func GetBlock(c *client.Client, blockNumber *big.Int) (*connection.BlockType, error) {
	params := client.ParamsOfQueryCollection{
		Collection: "blocks",
		Result:     "seq_no gen_utime",
		Limit:      null.Uint32From(1),
		Filter: json.RawMessage(`{
			"workchain_id":{"eq":-1},
			"status":{"eq": 2},
			"seq_no":{"eq": ` + blockNumber.String() + `}
		}`),
		Order: []client.OrderBy{{
			Path:      "seq_no",
			Direction: client.DescSortDirection,
		}},
	}

	res, err := c.NetQueryCollection(&params)

	if err != nil {
		return nil, err
	}

	if len(res.Result) <= 0 {
		return nil, errors.New("No blocks found")
	}

	currentBlock := &connection.BlockType{}

	s, err := json.Marshal(res.Result[0])

	if err != nil {
		return nil, err
	}

	json.Unmarshal(s, &currentBlock)

	return currentBlock, nil
}

func DecodeMessageBody(c *client.Client, message *Message, abi client.Abi) (*client.DecodedMessageBody, error) {
	params := client.ParamsOfDecodeMessageBody{
		Abi:  abi,
		Body: ((*message)["body"]).(string),
	}

	return c.AbiDecodeMessageBody(&params)
}

func GetMessage(c *client.Client, address string, prevBlock, currentBlock *connection.BlockType) (*[]Message, error) {
	messages := []Message{}

	params := client.ParamsOfQueryCollection{
		Collection: "messages",
		Result:     "id status created_at body",
		Filter: json.RawMessage(`{
			"status": { "eq": 5 },
			"src": { "eq": "` + address + `" },
			"created_at": {
				"ge": ` + strconv.FormatInt(prevBlock.CreatedAt, 10) + `,
				"lt": ` + strconv.FormatInt(currentBlock.CreatedAt, 10) + `
			}
		}`),
	}

	res, err := c.NetQueryCollection(&params)
	if err != nil {
		return &messages, err
	}

	if len(res.Result) <= 0 {
		return &messages, nil
	}

	for _, s := range res.Result {
		messages = append(messages, s.(Message))
	}

	return &messages, nil
}