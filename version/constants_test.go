// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrentRPCChainVMCompatible(t *testing.T) {
	compatibleVersions := RPCChainVMProtocolCompatibility[RPCChainVMProtocol]
	require.Contains(t, compatibleVersions, Current)
}
