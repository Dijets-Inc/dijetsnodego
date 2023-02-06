// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"github.com/lasthyphen/dijetsnodego/vms/secp256k1fx"
)

type Credential struct {
	secp256k1fx.Credential `serialize:"true"`
}
