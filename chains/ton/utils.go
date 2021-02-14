// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"strconv"

	connection "github.com/ChainSafe/ChainBridge/connections/ton"
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
)

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

func LoadTvc(path, name string) string {
	content, err := ioutil.ReadFile(path + "/" + name + ".tvc")
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(content)
}

func GetBlock(c *client.Client, blockNumber *big.Int, workchainID string) (*connection.BlockType, error) {
	filter := json.RawMessage(`{
		"workchain_id":{"eq":` + workchainID + `},
		"status":{"eq": 2},
		"seq_no":{"eq": ` + blockNumber.String() + `}
	}`)

	params := client.ParamsOfQueryCollection{
		Collection: "blocks",
		Result:     "seq_no gen_utime",
		Limit:      null.Uint32From(1),
		Filter:     filter,
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
		return nil, errors.New("No blocks found: " + string(filter))
	}

	currentBlock := &connection.BlockType{}

	s, err := json.Marshal(res.Result[0])

	if err != nil {
		return nil, err
	}

	json.Unmarshal(s, &currentBlock)

	return currentBlock, nil
}

func DecodeMessageBody(c *client.Client, message *json.RawMessage, abi client.Abi) (*client.DecodedMessageBody, error) {
	msg, err := message.MarshalJSON()
	if err != nil {
		return nil, err
	}

	params := client.ParamsOfDecodeMessageBody{
		Abi:  abi,
		Body: string(msg),
	}

	return c.AbiDecodeMessageBody(&params)
}

func GetMessage(c *client.Client, address string, prevBlock, currentBlock *connection.BlockType) (*[]json.RawMessage, error) {
	messages := []json.RawMessage{}
	// FIXME: receive dst_transaction and check aborted field
	params := client.ParamsOfQueryCollection{
		Collection: "messages",
		Result:     "id status created_at body src",
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
		messages = append(messages, s)
	}

	return &messages, nil
}
