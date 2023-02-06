// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package x

import (
	"github.com/lasthyphen/dijetsnodego/vms/avm/fxs"
	"github.com/lasthyphen/dijetsnodego/vms/avm/txs"
	"github.com/lasthyphen/dijetsnodego/vms/nftfx"
	"github.com/lasthyphen/dijetsnodego/vms/propertyfx"
	"github.com/lasthyphen/dijetsnodego/vms/secp256k1fx"
)

const (
	SECP256K1FxIndex = 0
	NFTFxIndex       = 1
	PropertyFxIndex  = 2
)

// Parser to support serialization and deserialization
var Parser txs.Parser

func init() {
	var err error
	Parser, err = txs.NewParser([]fxs.Fx{
		&secp256k1fx.Fx{},
		&nftfx.Fx{},
		&propertyfx.Fx{},
	})
	if err != nil {
		panic(err)
	}
}
