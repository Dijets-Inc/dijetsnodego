// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vertex

import (
	"context"

	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/snow/consensus/snowstorm"
	"github.com/lasthyphen/dijetsnodego/snow/engine/common"
)

// DAGVM defines the minimum functionality that an avalanche VM must
// implement
type DAGVM interface {
	common.VM
	Getter

	// Return any transactions that have not been sent to consensus yet
	PendingTxs(ctx context.Context) []snowstorm.Tx

	// Convert a stream of bytes to a transaction or return an error
	ParseTx(ctx context.Context, txBytes []byte) (snowstorm.Tx, error)
}

// Getter defines the functionality for fetching a tx/block by its ID.
type Getter interface {
	// Retrieve a transaction that was submitted previously
	GetTx(ctx context.Context, txID ids.ID) (snowstorm.Tx, error)
}
