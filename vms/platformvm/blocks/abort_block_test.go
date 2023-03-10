// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/ids"
)

func TestNewBanffAbortBlock(t *testing.T) {
	require := require.New(t)

	timestamp := time.Now().Truncate(time.Second)
	parentID := ids.GenerateTestID()
	height := uint64(1337)
	blk, err := NewBanffAbortBlock(
		timestamp,
		parentID,
		height,
	)
	require.NoError(err)

	// Make sure the block is initialized
	require.NotNil(blk.Bytes())

	require.Equal(timestamp, blk.Timestamp())
	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}

func TestNewApricotAbortBlock(t *testing.T) {
	require := require.New(t)

	parentID := ids.GenerateTestID()
	height := uint64(1337)
	blk, err := NewApricotAbortBlock(
		parentID,
		height,
	)
	require.NoError(err)

	// Make sure the block is initialized
	require.NotNil(blk.Bytes())

	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}
