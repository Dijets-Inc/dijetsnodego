// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/chains/atomic"
	"github.com/lasthyphen/dijetsnodego/database/prefixdb"
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/utils/crypto"
	"github.com/lasthyphen/dijetsnodego/vms/components/djtx"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/state"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/txs"
	"github.com/lasthyphen/dijetsnodego/vms/secp256k1fx"
)

func TestNewImportTx(t *testing.T) {
	env := newEnvironment( /*postBanff*/ false)
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()

	type test struct {
		description   string
		sourceChainID ids.ID
		sharedMemory  atomic.SharedMemory
		sourceKeys    []*crypto.PrivateKeySECP256K1R
		timestamp     time.Time
		shouldErr     bool
		shouldVerify  bool
	}

	factory := crypto.FactorySECP256K1R{}
	sourceKeyIntf, err := factory.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	sourceKey := sourceKeyIntf.(*crypto.PrivateKeySECP256K1R)

	cnt := new(byte)

	// Returns a shared memory where GetDatabase returns a database
	// where [recipientKey] has a balance of [amt]
	fundedSharedMemory := func(peerChain ids.ID, assets map[ids.ID]uint64) atomic.SharedMemory {
		*cnt++
		m := atomic.NewMemory(prefixdb.New([]byte{*cnt}, env.baseDB))

		sm := m.NewSharedMemory(env.ctx.ChainID)
		peerSharedMemory := m.NewSharedMemory(peerChain)

		for assetID, amt := range assets {
			// #nosec G404
			utxo := &djtx.UTXO{
				UTXOID: djtx.UTXOID{
					TxID:        ids.GenerateTestID(),
					OutputIndex: rand.Uint32(),
				},
				Asset: djtx.Asset{ID: assetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: amt,
					OutputOwners: secp256k1fx.OutputOwners{
						Locktime:  0,
						Addrs:     []ids.ShortID{sourceKey.PublicKey().Address()},
						Threshold: 1,
					},
				},
			}
			utxoBytes, err := txs.Codec.Marshal(txs.Version, utxo)
			if err != nil {
				t.Fatal(err)
			}
			inputID := utxo.InputID()
			if err := peerSharedMemory.Apply(map[ids.ID]*atomic.Requests{env.ctx.ChainID: {PutRequests: []*atomic.Element{{
				Key:   inputID[:],
				Value: utxoBytes,
				Traits: [][]byte{
					sourceKey.PublicKey().Address().Bytes(),
				},
			}}}}); err != nil {
				t.Fatal(err)
			}
		}

		return sm
	}

	customAssetID := ids.GenerateTestID()

	tests := []test{
		{
			description:   "can't pay fee",
			sourceChainID: env.ctx.XChainID,
			sharedMemory: fundedSharedMemory(
				env.ctx.XChainID,
				map[ids.ID]uint64{
					env.ctx.DJTXAssetID: env.config.TxFee - 1,
				},
			),
			sourceKeys: []*crypto.PrivateKeySECP256K1R{sourceKey},
			shouldErr:  true,
		},
		{
			description:   "can barely pay fee",
			sourceChainID: env.ctx.XChainID,
			sharedMemory: fundedSharedMemory(
				env.ctx.XChainID,
				map[ids.ID]uint64{
					env.ctx.DJTXAssetID: env.config.TxFee,
				},
			),
			sourceKeys:   []*crypto.PrivateKeySECP256K1R{sourceKey},
			shouldErr:    false,
			shouldVerify: true,
		},
		{
			description:   "attempting to import from C-chain",
			sourceChainID: cChainID,
			sharedMemory: fundedSharedMemory(
				cChainID,
				map[ids.ID]uint64{
					env.ctx.DJTXAssetID: env.config.TxFee,
				},
			),
			sourceKeys:   []*crypto.PrivateKeySECP256K1R{sourceKey},
			timestamp:    env.config.ApricotPhase5Time,
			shouldErr:    false,
			shouldVerify: true,
		},
		{
			description:   "attempting to import non-djtx from X-chain",
			sourceChainID: env.ctx.XChainID,
			sharedMemory: fundedSharedMemory(
				env.ctx.XChainID,
				map[ids.ID]uint64{
					env.ctx.DJTXAssetID: env.config.TxFee,
					customAssetID:       1,
				},
			),
			sourceKeys:   []*crypto.PrivateKeySECP256K1R{sourceKey},
			timestamp:    env.config.BanffTime,
			shouldErr:    false,
			shouldVerify: true,
		},
	}

	to := ids.GenerateTestShortID()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require := require.New(t)

			env.msm.SharedMemory = tt.sharedMemory
			tx, err := env.txBuilder.NewImportTx(
				tt.sourceChainID,
				to,
				tt.sourceKeys,
				ids.ShortEmpty,
			)
			if tt.shouldErr {
				require.Error(err)
				return
			}
			require.NoError(err)

			unsignedTx := tx.Unsigned.(*txs.ImportTx)
			require.NotEmpty(unsignedTx.ImportedInputs)
			require.Equal(len(tx.Creds), len(unsignedTx.Ins)+len(unsignedTx.ImportedInputs), "should have the same number of credentials as inputs")

			totalIn := uint64(0)
			for _, in := range unsignedTx.Ins {
				totalIn += in.Input().Amount()
			}
			for _, in := range unsignedTx.ImportedInputs {
				totalIn += in.Input().Amount()
			}
			totalOut := uint64(0)
			for _, out := range unsignedTx.Outs {
				totalOut += out.Out.Amount()
			}

			require.Equal(env.config.TxFee, totalIn-totalOut, "burned too much")

			fakedState, err := state.NewDiff(lastAcceptedID, env)
			require.NoError(err)

			fakedState.SetTimestamp(tt.timestamp)

			fakedParent := ids.GenerateTestID()
			env.SetState(fakedParent, fakedState)

			verifier := MempoolTxVerifier{
				Backend:       &env.backend,
				ParentID:      fakedParent,
				StateVersions: env,
				Tx:            tx,
			}
			err = tx.Unsigned.Visit(&verifier)
			if tt.shouldVerify {
				require.NoError(err)
			} else {
				require.Error(err)
			}
		})
	}
}
