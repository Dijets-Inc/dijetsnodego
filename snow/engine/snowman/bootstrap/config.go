// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package bootstrap

import (
	"github.com/lasthyphen/dijetsnodego/snow/engine/common"
	"github.com/lasthyphen/dijetsnodego/snow/engine/common/queue"
	"github.com/lasthyphen/dijetsnodego/snow/engine/snowman/block"
)

type Config struct {
	common.Config
	common.AllGetsServer

	// Blocked tracks operations that are blocked on blocks
	//
	// It should be guaranteed that `MissingIDs` should contain all IDs
	// referenced by the `MissingDependencies` that have not already been added
	// to the queue.
	Blocked *queue.JobsWithMissing

	VM block.ChainVM

	Bootstrapped func()
}
