// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"github.com/radianceteam/ton-client-go/client"
	"github.com/volatiletech/null"
)

func NewClient() (*client.Client, error) {
	conn, err := client.NewClient(client.Config{
		Network: &client.NetworkConfig{ServerAddress: null.StringFrom("")},
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
