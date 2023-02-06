// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package bootstrap

import (
	"github.com/lasthyphen/dijetsnodego/snow/engine/avalanche/vertex"
	"github.com/lasthyphen/dijetsnodego/snow/engine/common"
	"github.com/lasthyphen/dijetsnodego/snow/engine/common/queue"
)

type Config struct {
	common.Config
	common.AllGetsServer

	// VtxBlocked tracks operations that are blocked on vertices
	VtxBlocked *queue.JobsWithMissing
	// TxBlocked tracks operations that are blocked on transactions
	TxBlocked *queue.Jobs

	Manager vertex.Manager
	VM      vertex.DAGVM
}
