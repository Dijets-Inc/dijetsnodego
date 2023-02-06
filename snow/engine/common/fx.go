// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package common

import (
	"github.com/lasthyphen/dijetsnodego/ids"
)

// Fx wraps an instance of a feature extension
type Fx struct {
	ID ids.ID
	Fx interface{}
}
