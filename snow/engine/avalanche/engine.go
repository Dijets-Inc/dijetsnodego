// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avalanche

import (
	"context"

	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/snow/consensus/avalanche"
	"github.com/lasthyphen/dijetsnodego/snow/engine/common"
)

// Engine describes the events that can occur on a consensus instance
type Engine interface {
	common.Engine

	// GetVtx returns a vertex by its ID.
	// Returns an error if unknown.
	GetVtx(ctx context.Context, vtxID ids.ID) (avalanche.Vertex, error)
}
