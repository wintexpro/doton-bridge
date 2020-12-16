// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import "github.com/radianceteam/ton-client-go/client"

func NewClient() (*client.Client, error) {
	conn, err := client.NewClient(client.Config{
		Network: &client.NetworkConfig{ServerAddress: ""},
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}
