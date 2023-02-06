// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package secp256k1fx

import (
	"github.com/lasthyphen/dijetsnodego/vms/components/verify"
)

var _ verify.State = (*MintOutput)(nil)

type MintOutput struct {
	OutputOwners `serialize:"true"`
}

func (out *MintOutput) Verify() error {
	switch {
	case out == nil:
		return errNilOutput
	default:
		return out.OutputOwners.Verify()
	}
}

func (out *MintOutput) VerifyState() error {
	return out.Verify()
}
