// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package signer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	require := require.New(t)

	noSigner := &Empty{}
	require.NoError(noSigner.Verify())
	require.Nil(noSigner.Key())
}
