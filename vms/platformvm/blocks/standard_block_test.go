// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/vms/components/djtx"
	"github.com/lasthyphen/dijetsnodego/vms/components/verify"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/txs"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/validator"
	"github.com/lasthyphen/dijetsnodego/vms/secp256k1fx"
)

func TestNewBanffStandardBlock(t *testing.T) {
	require := require.New(t)

	timestamp := time.Now().Truncate(time.Second)
	parentID := ids.GenerateTestID()
	height := uint64(1337)

	tx := &txs.Tx{
		Unsigned: &txs.AddValidatorTx{
			BaseTx: txs.BaseTx{
				BaseTx: djtx.BaseTx{
					Ins:  []*djtx.TransferableInput{},
					Outs: []*djtx.TransferableOutput{},
				},
			},
			StakeOuts: []*djtx.TransferableOutput{},
			Validator: validator.Validator{},
			RewardsOwner: &secp256k1fx.OutputOwners{
				Addrs: []ids.ShortID{},
			},
		},
		Creds: []verify.Verifiable{},
	}
	require.NoError(tx.Sign(txs.Codec, nil))

	blk, err := NewBanffStandardBlock(
		timestamp,
		parentID,
		height,
		[]*txs.Tx{tx},
	)
	require.NoError(err)

	// Make sure the block and tx are initialized
	require.NotNil(blk.Bytes())
	require.NotNil(blk.Transactions[0].Bytes())
	require.NotEqual(ids.Empty, blk.Transactions[0].ID())
	require.Equal(tx.Bytes(), blk.Transactions[0].Bytes())
	require.Equal(timestamp, blk.Timestamp())
	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}

func TestNewApricotStandardBlock(t *testing.T) {
	require := require.New(t)

	parentID := ids.GenerateTestID()
	height := uint64(1337)

	tx := &txs.Tx{
		Unsigned: &txs.AddValidatorTx{
			BaseTx: txs.BaseTx{
				BaseTx: djtx.BaseTx{
					Ins:  []*djtx.TransferableInput{},
					Outs: []*djtx.TransferableOutput{},
				},
			},
			StakeOuts: []*djtx.TransferableOutput{},
			Validator: validator.Validator{},
			RewardsOwner: &secp256k1fx.OutputOwners{
				Addrs: []ids.ShortID{},
			},
		},
		Creds: []verify.Verifiable{},
	}
	require.NoError(tx.Sign(txs.Codec, nil))

	blk, err := NewApricotStandardBlock(
		parentID,
		height,
		[]*txs.Tx{tx},
	)
	require.NoError(err)

	// Make sure the block and tx are initialized
	require.NotNil(blk.Bytes())
	require.NotNil(blk.Transactions[0].Bytes())
	require.NotEqual(ids.Empty, blk.Transactions[0].ID())
	require.Equal(tx.Bytes(), blk.Transactions[0].Bytes())
	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}
