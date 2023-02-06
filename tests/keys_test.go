// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/utils/crypto"
)

func TestLoadTestKeys(t *testing.T) {
	keys, err := LoadHexTestKeys("test.insecure.secp256k1.keys")
	require.NoError(t, err)
	for i, k := range keys {
		curAddr := encodeShortAddr(k)
		t.Logf("[%d] loaded %v", i, curAddr)
	}
}

func encodeShortAddr(pk *crypto.PrivateKeySECP256K1R) string {
	return pk.PublicKey().Address().String()
}
