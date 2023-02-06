// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package proposer

import (
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/utils"
)

var _ utils.Sortable[validatorData] = validatorData{}

type validatorData struct {
	id     ids.NodeID
	weight uint64
}

func (d validatorData) Less(other validatorData) bool {
	return d.id.Less(other.id)
}
