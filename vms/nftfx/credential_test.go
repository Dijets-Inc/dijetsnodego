// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nftfx

import (
	"testing"

	"github.com/lasthyphen/dijetsnodego/vms/components/verify"
)

func TestCredentialState(t *testing.T) {
	intf := interface{}(&Credential{})
	if _, ok := intf.(verify.State); ok {
		t.Fatalf("shouldn't be marked as state")
	}
}
