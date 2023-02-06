// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/vms/components/verify"
	"github.com/lasthyphen/dijetsnodego/vms/secp256k1fx"
)

type BurnOperation struct {
	secp256k1fx.Input `serialize:"true"`
}

func (*BurnOperation) InitCtx(*snow.Context) {}

func (*BurnOperation) Outs() []verify.State {
	return nil
}
