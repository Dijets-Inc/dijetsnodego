// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"time"

	"github.com/lasthyphen/dijetsnodego/chains/atomic"
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/utils/set"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/blocks"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/state"
)

type standardBlockState struct {
	onAcceptFunc func()
	inputs       set.Set[ids.ID]
}

type proposalBlockState struct {
	initiallyPreferCommit bool
	onCommitState         state.Diff
	onAbortState          state.Diff
}

// The state of a block.
// Note that not all fields will be set for a given block.
type blockState struct {
	standardBlockState
	proposalBlockState
	statelessBlock blocks.Block
	onAcceptState  state.Diff

	timestamp      time.Time
	atomicRequests map[ids.ID]*atomic.Requests
}
