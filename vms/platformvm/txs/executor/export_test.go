// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/utils/crypto"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/state"
)

func TestNewExportTx(t *testing.T) {
	env := newEnvironment( /*postBanff*/ true)
	env.ctx.Lock.Lock()
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()

	type test struct {
		description        string
		destinationChainID ids.ID
		sourceKeys         []*crypto.PrivateKeySECP256K1R
		timestamp          time.Time
		shouldErr          bool
		shouldVerify       bool
	}

	sourceKey := preFundedKeys[0]

	tests := []test{
		{
			description:        "P->X export",
			destinationChainID: xChainID,
			sourceKeys:         []*crypto.PrivateKeySECP256K1R{sourceKey},
			timestamp:          defaultValidateStartTime,
			shouldErr:          false,
			shouldVerify:       true,
		},
		{
			description:        "P->C export",
			destinationChainID: cChainID,
			sourceKeys:         []*crypto.PrivateKeySECP256K1R{sourceKey},
			timestamp:          env.config.ApricotPhase5Time,
			shouldErr:          false,
			shouldVerify:       true,
		},
	}

	to := ids.GenerateTestShortID()
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require := require.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tx, err := env.txBuilder.NewExportTx(
				defaultBalance-defaultTxFee, // Amount of tokens to export
				tt.destinationChainID,
				to,
				tt.sourceKeys,
				ids.ShortEmpty, // Change address
			)
			if tt.shouldErr {
				require.Error(err)
				return
			}
			require.NoError(err)

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
