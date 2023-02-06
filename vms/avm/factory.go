// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/vms"
)

var _ vms.Factory = (*Factory)(nil)

type Factory struct {
	TxFee            uint64
	CreateAssetTxFee uint64
}

func (f *Factory) New(*snow.Context) (interface{}, error) {
	return &VM{Factory: *f}, nil
}
