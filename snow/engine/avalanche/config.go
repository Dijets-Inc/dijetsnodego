// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avalanche

import (
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/snow/consensus/avalanche"
	"github.com/lasthyphen/dijetsnodego/snow/engine/avalanche/vertex"
	"github.com/lasthyphen/dijetsnodego/snow/engine/common"
	"github.com/lasthyphen/dijetsnodego/snow/validators"
)

// Config wraps all the parameters needed for an avalanche engine
type Config struct {
	Ctx *snow.ConsensusContext
	common.AllGetsServer
	VM         vertex.DAGVM
	Manager    vertex.Manager
	Sender     common.Sender
	Validators validators.Set

	Params    avalanche.Parameters
	Consensus avalanche.Consensus
}
