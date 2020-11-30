// Copyright 2020 Wintex
// SPDX-License-Identifier: LGPL-3.0-only

package ton

import (
	"github.com/ByKeks/chainbridge-utils/msg"
	"github.com/ChainSafe/log15"
)

type eventName string
type eventHandler func(interface{}, log15.Logger) (msg.Message, error)
