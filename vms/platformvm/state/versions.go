// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"github.com/lasthyphen/dijetsnodego/ids"
)

type Versions interface {
	GetState(blkID ids.ID) (Chain, bool)
}
