// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package states

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/database"
	"github.com/lasthyphen/dijetsnodego/database/memdb"
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/utils/crypto"
	"github.com/lasthyphen/dijetsnodego/utils/units"
	"github.com/lasthyphen/dijetsnodego/vms/avm/fxs"
	"github.com/lasthyphen/dijetsnodego/vms/avm/txs"
	"github.com/lasthyphen/dijetsnodego/vms/components/djtx"
	"github.com/lasthyphen/dijetsnodego/vms/nftfx"
	"github.com/lasthyphen/dijetsnodego/vms/propertyfx"
	"github.com/lasthyphen/dijetsnodego/vms/secp256k1fx"
)

var (
	networkID uint32 = 10
	chainID          = ids.ID{5, 4, 3, 2, 1}
	assetID          = ids.ID{1, 2, 3}
	keys             = crypto.BuildTestKeys()
)

func TestTxState(t *testing.T) {
	require := require.New(t)

	db := memdb.New()
	parser, err := txs.NewParser([]fxs.Fx{
		&secp256k1fx.Fx{},
		&nftfx.Fx{},
		&propertyfx.Fx{},
	})
	require.NoError(err)

	stateIntf, err := NewTxState(db, parser, prometheus.NewRegistry())
	require.NoError(err)

	s := stateIntf.(*txState)

	_, err = s.GetTx(ids.Empty)
	require.Equal(database.ErrNotFound, err)

	tx := &txs.Tx{
		Unsigned: &txs.BaseTx{
			BaseTx: djtx.BaseTx{
				NetworkID:    networkID,
				BlockchainID: chainID,
				Ins: []*djtx.TransferableInput{{
					UTXOID: djtx.UTXOID{
						TxID:        ids.Empty,
						OutputIndex: 0,
					},
					Asset: djtx.Asset{ID: assetID},
					In: &secp256k1fx.TransferInput{
						Amt: 20 * units.KiloDjtx,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{
								0,
							},
						},
					},
				}},
			},
		},
	}

	err = tx.SignSECP256K1Fx(parser.Codec(), [][]*crypto.PrivateKeySECP256K1R{{keys[0]}})
	require.NoError(err)

	err = s.PutTx(ids.Empty, tx)
	require.NoError(err)

	loadedTx, err := s.GetTx(ids.Empty)
	require.NoError(err)
	require.Equal(tx.ID(), loadedTx.ID())

	s.txCache.Flush()

	loadedTx, err = s.GetTx(ids.Empty)
	require.NoError(err)
	require.Equal(tx.ID(), loadedTx.ID())

	err = s.DeleteTx(ids.Empty)
	require.NoError(err)

	_, err = s.GetTx(ids.Empty)
	require.Equal(database.ErrNotFound, err)

	s.txCache.Flush()

	_, err = s.GetTx(ids.Empty)
	require.Equal(database.ErrNotFound, err)
}
