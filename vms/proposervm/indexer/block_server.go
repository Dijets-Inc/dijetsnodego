// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package indexer

import (
	"context"

	"github.com/lasthyphen/dijetsnodego/database/versiondb"
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/snow/consensus/snowman"
)

// BlockServer represents all requests heightIndexer can issue
// against ProposerVM. All methods must be thread-safe.
type BlockServer interface {
	versiondb.Commitable

	// Note: this is a contention heavy call that should be avoided
	// for frequent/repeated indexer ops
	GetFullPostForkBlock(ctx context.Context, blkID ids.ID) (snowman.Block, error)
}
