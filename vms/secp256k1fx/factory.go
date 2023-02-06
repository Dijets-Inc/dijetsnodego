// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package secp256k1fx

import (
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/vms"
)

var (
	_ vms.Factory = (*Factory)(nil)

	// ID that this Fx uses when labeled
	ID = ids.ID{'s', 'e', 'c', 'p', '2', '5', '6', 'k', '1', 'f', 'x'}
)

type Factory struct{}

func (*Factory) New(*snow.Context) (interface{}, error) {
	return &Fx{}, nil
}
