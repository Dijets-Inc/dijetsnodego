// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avalanche

import (
	"testing"
)

func TestTopological(t *testing.T) {
	runConsensusTests(t, TopologicalFactory{})
}
