// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ulimit

import (
	"testing"

	"github.com/lasthyphen/dijetsnodego/utils/logging"
)

// Test_SetDefault performs sanity checks for the os default.
func Test_SetDefault(t *testing.T) {
	err := Set(DefaultFDLimit, logging.NoLog{})
	if err != nil {
		t.Skipf("default fd-limit failed %v", err)
	}
}
