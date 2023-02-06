// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/vms"
)

var (
	_ vms.Factory = (*Factory)(nil)

	// ID that this Fx uses when labeled
	ID = ids.ID{'p', 'r', 'o', 'p', 'e', 'r', 't', 'y', 'f', 'x'}
)

type Factory struct{}

func (*Factory) New(*snow.Context) (interface{}, error) {
	return &Fx{}, nil
}
