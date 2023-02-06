// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowstorm

import (
	"testing"
)

func TestDirectedConsensus(t *testing.T) {
	runConsensusTests(t, DirectedFactory{}, "DG")
}
