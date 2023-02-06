// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

import (
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/snow/consensus/snowball"
	"github.com/lasthyphen/dijetsnodego/snow/consensus/snowman"
	"github.com/lasthyphen/dijetsnodego/snow/engine/common"
	"github.com/lasthyphen/dijetsnodego/snow/engine/snowman/block"
	"github.com/lasthyphen/dijetsnodego/snow/validators"
)

// Config wraps all the parameters needed for a snowman engine
type Config struct {
	common.AllGetsServer

	Ctx        *snow.ConsensusContext
	VM         block.ChainVM
	Sender     common.Sender
	Validators validators.Set
	Params     snowball.Parameters
	Consensus  snowman.Consensus
}
