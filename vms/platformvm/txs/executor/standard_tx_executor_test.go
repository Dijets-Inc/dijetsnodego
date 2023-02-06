// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/database"
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/utils"
	"github.com/lasthyphen/dijetsnodego/utils/constants"
	"github.com/lasthyphen/dijetsnodego/utils/crypto"
	"github.com/lasthyphen/dijetsnodego/utils/hashing"
	"github.com/lasthyphen/dijetsnodego/vms/components/djtx"
	"github.com/lasthyphen/dijetsnodego/vms/components/verify"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/config"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/fx"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/reward"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/state"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/status"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/txs"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/utxo"
	"github.com/lasthyphen/dijetsnodego/vms/secp256k1fx"
)

// This tests that the math performed during TransformSubnetTx execution can
// never overflow
const _ time.Duration = math.MaxUint32 * time.Second

func TestStandardTxExecutorAddValidatorTxEmptyID(t *testing.T) {
	env := newEnvironment( /*postBanff*/ false)
	env.ctx.Lock.Lock()
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()

	chainTime := env.state.GetTimestamp()
	startTime := defaultGenesisTime.Add(1 * time.Second)

	tests := []struct {
		banffTime     time.Time
		expectedError error
	}{
		{ // Case: Before banff
			banffTime:     chainTime.Add(1),
			expectedError: errEmptyNodeID,
		},
		{ // Case: At banff
			banffTime:     chainTime,
			expectedError: errEmptyNodeID,
		},
		{ // Case: After banff
			banffTime:     chainTime.Add(-1),
			expectedError: errEmptyNodeID,
		},
	}
	for _, test := range tests {
		// Case: Empty validator node ID after banff
		env.config.BanffTime = test.banffTime

		tx, err := env.txBuilder.NewAddValidatorTx( // create the tx
			env.config.MinValidatorStake,
			uint64(startTime.Unix()),
			uint64(defaultValidateEndTime.Unix()),
			ids.EmptyNodeID,
			ids.GenerateTestShortID(),
			reward.PercentDenominator,
			[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		stateDiff, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   stateDiff,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		require.ErrorIs(t, err, test.expectedError)
	}
}

func TestStandardTxExecutorAddDelegator(t *testing.T) {
	dummyHeight := uint64(1)
	rewardAddress := preFundedKeys[0].PublicKey().Address()
	nodeID := ids.NodeID(rewardAddress)

	newValidatorID := ids.GenerateTestNodeID()
	newValidatorStartTime := uint64(defaultValidateStartTime.Add(5 * time.Second).Unix())
	newValidatorEndTime := uint64(defaultValidateEndTime.Add(-5 * time.Second).Unix())

	// [addMinStakeValidator] adds a new validator to the primary network's
	// pending validator set with the minimum staking amount
	addMinStakeValidator := func(target *environment) {
		tx, err := target.txBuilder.NewAddValidatorTx(
			target.config.MinValidatorStake, // stake amount
			newValidatorStartTime,           // start time
			newValidatorEndTime,             // end time
			newValidatorID,                  // node ID
			rewardAddress,                   // Reward Address
			reward.PercentDenominator,       // Shares
			[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			ids.ShortEmpty,
		)
		if err != nil {
			t.Fatal(err)
		}

		staker, err := state.NewCurrentStaker(
			tx.ID(),
			tx.Unsigned.(*txs.AddValidatorTx),
			0,
		)
		if err != nil {
			t.Fatal(err)
		}

		target.state.PutCurrentValidator(staker)
		target.state.AddTx(tx, status.Committed)
		target.state.SetHeight(dummyHeight)
		if err := target.state.Commit(); err != nil {
			t.Fatal(err)
		}
	}

	// [addMaxStakeValidator] adds a new validator to the primary network's
	// pending validator set with the maximum staking amount
	addMaxStakeValidator := func(target *environment) {
		tx, err := target.txBuilder.NewAddValidatorTx(
			target.config.MaxValidatorStake, // stake amount
			newValidatorStartTime,           // start time
			newValidatorEndTime,             // end time
			newValidatorID,                  // node ID
			rewardAddress,                   // Reward Address
			reward.PercentDenominator,       // Shared
			[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			ids.ShortEmpty,
		)
		if err != nil {
			t.Fatal(err)
		}

		staker, err := state.NewCurrentStaker(
			tx.ID(),
			tx.Unsigned.(*txs.AddValidatorTx),
			0,
		)
		if err != nil {
			t.Fatal(err)
		}

		target.state.PutCurrentValidator(staker)
		target.state.AddTx(tx, status.Committed)
		target.state.SetHeight(dummyHeight)
		if err := target.state.Commit(); err != nil {
			t.Fatal(err)
		}
	}

	dummyH := newEnvironment( /*postBanff*/ false)
	currentTimestamp := dummyH.state.GetTimestamp()

	type test struct {
		stakeAmount   uint64
		startTime     uint64
		endTime       uint64
		nodeID        ids.NodeID
		rewardAddress ids.ShortID
		feeKeys       []*crypto.PrivateKeySECP256K1R
		setup         func(*environment)
		AP3Time       time.Time
		shouldErr     bool
		description   string
	}

	tests := []test{
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     uint64(defaultValidateStartTime.Unix()),
			endTime:       uint64(defaultValidateEndTime.Unix()) + 1,
			nodeID:        nodeID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         nil,
			AP3Time:       defaultGenesisTime,
			shouldErr:     true,
			description:   "validator stops validating primary network earlier than subnet",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     uint64(currentTimestamp.Add(MaxFutureStartTime + time.Second).Unix()),
			endTime:       uint64(currentTimestamp.Add(MaxFutureStartTime * 2).Unix()),
			nodeID:        nodeID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         nil,
			AP3Time:       defaultGenesisTime,
			shouldErr:     true,
			description:   fmt.Sprintf("validator should not be added more than (%s) in the future", MaxFutureStartTime),
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     uint64(defaultValidateStartTime.Unix()),
			endTime:       uint64(defaultValidateEndTime.Unix()) + 1,
			nodeID:        nodeID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         nil,
			AP3Time:       defaultGenesisTime,
			shouldErr:     true,
			description:   "end time is after the primary network end time",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     uint64(defaultValidateStartTime.Add(5 * time.Second).Unix()),
			endTime:       uint64(defaultValidateEndTime.Add(-5 * time.Second).Unix()),
			nodeID:        newValidatorID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         nil,
			AP3Time:       defaultGenesisTime,
			shouldErr:     true,
			description:   "validator not in the current or pending validator sets of the subnet",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     newValidatorStartTime - 1, // start validating subnet before primary network
			endTime:       newValidatorEndTime,
			nodeID:        newValidatorID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         addMinStakeValidator,
			AP3Time:       defaultGenesisTime,
			shouldErr:     true,
			description:   "validator starts validating subnet before primary network",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     newValidatorStartTime,
			endTime:       newValidatorEndTime + 1, // stop validating subnet after stopping validating primary network
			nodeID:        newValidatorID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         addMinStakeValidator,
			AP3Time:       defaultGenesisTime,
			shouldErr:     true,
			description:   "validator stops validating primary network before subnet",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     newValidatorStartTime, // same start time as for primary network
			endTime:       newValidatorEndTime,   // same end time as for primary network
			nodeID:        newValidatorID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         addMinStakeValidator,
			AP3Time:       defaultGenesisTime,
			shouldErr:     false,
			description:   "valid",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,                  // weight
			startTime:     uint64(currentTimestamp.Unix()),                  // start time
			endTime:       uint64(defaultValidateEndTime.Unix()),            // end time
			nodeID:        nodeID,                                           // node ID
			rewardAddress: rewardAddress,                                    // Reward Address
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]}, // tx fee payer
			setup:         nil,
			AP3Time:       defaultGenesisTime,
			shouldErr:     true,
			description:   "starts validating at current timestamp",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,                  // weight
			startTime:     uint64(defaultValidateStartTime.Unix()),          // start time
			endTime:       uint64(defaultValidateEndTime.Unix()),            // end time
			nodeID:        nodeID,                                           // node ID
			rewardAddress: rewardAddress,                                    // Reward Address
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[1]}, // tx fee payer
			setup: func(target *environment) { // Remove all UTXOs owned by keys[1]
				utxoIDs, err := target.state.UTXOIDs(
					preFundedKeys[1].PublicKey().Address().Bytes(),
					ids.Empty,
					math.MaxInt32)
				if err != nil {
					t.Fatal(err)
				}
				for _, utxoID := range utxoIDs {
					target.state.DeleteUTXO(utxoID)
				}
				target.state.SetHeight(dummyHeight)
				if err := target.state.Commit(); err != nil {
					t.Fatal(err)
				}
			},
			AP3Time:     defaultGenesisTime,
			shouldErr:   true,
			description: "tx fee paying key has no funds",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     newValidatorStartTime, // same start time as for primary network
			endTime:       newValidatorEndTime,   // same end time as for primary network
			nodeID:        newValidatorID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         addMaxStakeValidator,
			AP3Time:       defaultValidateEndTime,
			shouldErr:     false,
			description:   "over delegation before AP3",
		},
		{
			stakeAmount:   dummyH.config.MinDelegatorStake,
			startTime:     newValidatorStartTime, // same start time as for primary network
			endTime:       newValidatorEndTime,   // same end time as for primary network
			nodeID:        newValidatorID,
			rewardAddress: rewardAddress,
			feeKeys:       []*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			setup:         addMaxStakeValidator,
			AP3Time:       defaultGenesisTime,
			shouldErr:     true,
			description:   "over delegation after AP3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			freshTH := newEnvironment( /*postBanff*/ false)
			freshTH.config.ApricotPhase3Time = tt.AP3Time
			defer func() {
				if err := shutdownEnvironment(freshTH); err != nil {
					t.Fatal(err)
				}
			}()

			tx, err := freshTH.txBuilder.NewAddDelegatorTx(
				tt.stakeAmount,
				tt.startTime,
				tt.endTime,
				tt.nodeID,
				tt.rewardAddress,
				tt.feeKeys,
				ids.ShortEmpty,
			)
			if err != nil {
				t.Fatalf("couldn't build tx: %s", err)
			}
			if tt.setup != nil {
				tt.setup(freshTH)
			}

			onAcceptState, err := state.NewDiff(lastAcceptedID, freshTH)
			if err != nil {
				t.Fatal(err)
			}

			freshTH.config.BanffTime = onAcceptState.GetTimestamp()

			executor := StandardTxExecutor{
				Backend: &freshTH.backend,
				State:   onAcceptState,
				Tx:      tx,
			}
			err = tx.Unsigned.Visit(&executor)
			if err != nil && !tt.shouldErr {
				t.Fatalf("shouldn't have errored but got %s", err)
			} else if err == nil && tt.shouldErr {
				t.Fatalf("expected test to error but got none")
			}

			mempoolExecutor := MempoolTxVerifier{
				Backend:       &freshTH.backend,
				ParentID:      lastAcceptedID,
				StateVersions: freshTH,
				Tx:            tx,
			}
			err = tx.Unsigned.Visit(&mempoolExecutor)
			if err != nil && !tt.shouldErr {
				t.Fatalf("shouldn't have errored but got %s", err)
			} else if err == nil && tt.shouldErr {
				t.Fatalf("expected test to error but got none")
			}
		})
	}
}

func TestStandardTxExecutorAddSubnetValidator(t *testing.T) {
	env := newEnvironment( /*postBanff*/ false)
	env.ctx.Lock.Lock()
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()

	nodeID := preFundedKeys[0].PublicKey().Address()
	env.config.BanffTime = env.state.GetTimestamp()

	{
		// Case: Proposed validator currently validating primary network
		// but stops validating subnet after stops validating primary network
		// (note that keys[0] is a genesis validator)
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,
			uint64(defaultValidateStartTime.Unix()),
			uint64(defaultValidateEndTime.Unix())+1,
			ids.NodeID(nodeID),
			testSubnet1.ID(),
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed because validator stops validating primary network earlier than subnet")
		}
	}

	{
		// Case: Proposed validator currently validating primary network
		// and proposed subnet validation period is subset of
		// primary network validation period
		// (note that keys[0] is a genesis validator)
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,
			uint64(defaultValidateStartTime.Unix()+1),
			uint64(defaultValidateEndTime.Unix()),
			ids.NodeID(nodeID),
			testSubnet1.ID(),
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Add a validator to pending validator set of primary network
	key, err := testKeyfactory.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	pendingDSValidatorID := ids.NodeID(key.PublicKey().Address())

	// starts validating primary network 10 seconds after genesis
	dsStartTime := defaultGenesisTime.Add(10 * time.Second)
	dsEndTime := dsStartTime.Add(5 * defaultMinStakingDuration)

	addDSTx, err := env.txBuilder.NewAddValidatorTx(
		env.config.MinValidatorStake, // stake amount
		uint64(dsStartTime.Unix()),   // start time
		uint64(dsEndTime.Unix()),     // end time
		pendingDSValidatorID,         // node ID
		nodeID,                       // reward address
		reward.PercentDenominator,    // shares
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
		ids.ShortEmpty,
	)
	if err != nil {
		t.Fatal(err)
	}

	{
		// Case: Proposed validator isn't in pending or current validator sets
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,
			uint64(dsStartTime.Unix()), // start validating subnet before primary network
			uint64(dsEndTime.Unix()),
			pendingDSValidatorID,
			testSubnet1.ID(),
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed because validator not in the current or pending validator sets of the primary network")
		}
	}

	staker, err := state.NewCurrentStaker(
		addDSTx.ID(),
		addDSTx.Unsigned.(*txs.AddValidatorTx),
		0,
	)
	if err != nil {
		t.Fatal(err)
	}

	env.state.PutCurrentValidator(staker)
	env.state.AddTx(addDSTx, status.Committed)
	dummyHeight := uint64(1)
	env.state.SetHeight(dummyHeight)
	if err := env.state.Commit(); err != nil {
		t.Fatal(err)
	}

	// Node with ID key.PublicKey().Address() now a pending validator for primary network

	{
		// Case: Proposed validator is pending validator of primary network
		// but starts validating subnet before primary network
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,
			uint64(dsStartTime.Unix())-1, // start validating subnet before primary network
			uint64(dsEndTime.Unix()),
			pendingDSValidatorID,
			testSubnet1.ID(),
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed because validator starts validating primary network before starting to validate primary network")
		}
	}

	{
		// Case: Proposed validator is pending validator of primary network
		// but stops validating subnet after primary network
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,
			uint64(dsStartTime.Unix()),
			uint64(dsEndTime.Unix())+1, // stop validating subnet after stopping validating primary network
			pendingDSValidatorID,
			testSubnet1.ID(),
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed because validator stops validating primary network after stops validating primary network")
		}
	}

	{
		// Case: Proposed validator is pending validator of primary network and
		// period validating subnet is subset of time validating primary network
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,
			uint64(dsStartTime.Unix()), // same start time as for primary network
			uint64(dsEndTime.Unix()),   // same end time as for primary network
			pendingDSValidatorID,
			testSubnet1.ID(),
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Case: Proposed validator start validating at/before current timestamp
	// First, advance the timestamp
	newTimestamp := defaultGenesisTime.Add(2 * time.Second)
	env.state.SetTimestamp(newTimestamp)

	{
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,               // weight
			uint64(newTimestamp.Unix()), // start time
			uint64(newTimestamp.Add(defaultMinStakingDuration).Unix()), // end time
			ids.NodeID(nodeID), // node ID
			testSubnet1.ID(),   // subnet ID
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed verification because starts validating at current timestamp")
		}
	}

	// reset the timestamp
	env.state.SetTimestamp(defaultGenesisTime)

	// Case: Proposed validator already validating the subnet
	// First, add validator as validator of subnet
	subnetTx, err := env.txBuilder.NewAddSubnetValidatorTx(
		defaultWeight,                           // weight
		uint64(defaultValidateStartTime.Unix()), // start time
		uint64(defaultValidateEndTime.Unix()),   // end time
		ids.NodeID(nodeID),                      // node ID
		testSubnet1.ID(),                        // subnet ID
		[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
		ids.ShortEmpty,
	)
	if err != nil {
		t.Fatal(err)
	}

	staker, err = state.NewCurrentStaker(
		subnetTx.ID(),
		subnetTx.Unsigned.(*txs.AddSubnetValidatorTx),
		0,
	)
	if err != nil {
		t.Fatal(err)
	}

	env.state.PutCurrentValidator(staker)
	env.state.AddTx(subnetTx, status.Committed)
	env.state.SetHeight(dummyHeight)
	if err := env.state.Commit(); err != nil {
		t.Fatal(err)
	}

	{
		// Node with ID nodeIDKey.PublicKey().Address() now validating subnet with ID testSubnet1.ID
		duplicateSubnetTx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,                           // weight
			uint64(defaultValidateStartTime.Unix()), // start time
			uint64(defaultValidateEndTime.Unix()),   // end time
			ids.NodeID(nodeID),                      // node ID
			testSubnet1.ID(),                        // subnet ID
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      duplicateSubnetTx,
		}
		err = duplicateSubnetTx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed verification because validator already validating the specified subnet")
		}
	}

	env.state.DeleteCurrentValidator(staker)
	env.state.SetHeight(dummyHeight)
	if err := env.state.Commit(); err != nil {
		t.Fatal(err)
	}

	{
		// Case: Too many signatures
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,                     // weight
			uint64(defaultGenesisTime.Unix()), // start time
			uint64(defaultGenesisTime.Add(defaultMinStakingDuration).Unix())+1, // end time
			ids.NodeID(nodeID), // node ID
			testSubnet1.ID(),   // subnet ID
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1], testSubnet1ControlKeys[2]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed verification because tx has 3 signatures but only 2 needed")
		}
	}

	{
		// Case: Too few signatures
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,                     // weight
			uint64(defaultGenesisTime.Unix()), // start time
			uint64(defaultGenesisTime.Add(defaultMinStakingDuration).Unix()), // end time
			ids.NodeID(nodeID), // node ID
			testSubnet1.ID(),   // subnet ID
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[2]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		// Remove a signature
		addSubnetValidatorTx := tx.Unsigned.(*txs.AddSubnetValidatorTx)
		input := addSubnetValidatorTx.SubnetAuth.(*secp256k1fx.Input)
		input.SigIndices = input.SigIndices[1:]
		// This tx was syntactically verified when it was created...pretend it wasn't so we don't use cache
		addSubnetValidatorTx.SyntacticallyVerified = false

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed verification because not enough control sigs")
		}
	}

	{
		// Case: Control Signature from invalid key (keys[3] is not a control key)
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,                     // weight
			uint64(defaultGenesisTime.Unix()), // start time
			uint64(defaultGenesisTime.Add(defaultMinStakingDuration).Unix()), // end time
			ids.NodeID(nodeID), // node ID
			testSubnet1.ID(),   // subnet ID
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], preFundedKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}
		// Replace a valid signature with one from keys[3]
		sig, err := preFundedKeys[3].SignHash(hashing.ComputeHash256(tx.Unsigned.Bytes()))
		if err != nil {
			t.Fatal(err)
		}
		copy(tx.Creds[0].(*secp256k1fx.Credential).Sigs[0][:], sig)

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed verification because a control sig is invalid")
		}
	}

	{
		// Case: Proposed validator in pending validator set for subnet
		// First, add validator to pending validator set of subnet
		tx, err := env.txBuilder.NewAddSubnetValidatorTx(
			defaultWeight,                       // weight
			uint64(defaultGenesisTime.Unix())+1, // start time
			uint64(defaultGenesisTime.Add(defaultMinStakingDuration).Unix())+1, // end time
			ids.NodeID(nodeID), // node ID
			testSubnet1.ID(),   // subnet ID
			[]*crypto.PrivateKeySECP256K1R{testSubnet1ControlKeys[0], testSubnet1ControlKeys[1]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		staker, err = state.NewCurrentStaker(
			subnetTx.ID(),
			subnetTx.Unsigned.(*txs.AddSubnetValidatorTx),
			0,
		)
		if err != nil {
			t.Fatal(err)
		}

		env.state.PutCurrentValidator(staker)
		env.state.AddTx(tx, status.Committed)
		env.state.SetHeight(dummyHeight)
		if err := env.state.Commit(); err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed verification because validator already in pending validator set of the specified subnet")
		}
	}
}

func TestStandardTxExecutorAddValidator(t *testing.T) {
	env := newEnvironment( /*postBanff*/ false)
	env.ctx.Lock.Lock()
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()

	nodeID := ids.GenerateTestNodeID()

	env.config.BanffTime = env.state.GetTimestamp()

	{
		// Case: Validator's start time too early
		tx, err := env.txBuilder.NewAddValidatorTx(
			env.config.MinValidatorStake,
			uint64(defaultValidateStartTime.Unix())-1,
			uint64(defaultValidateEndTime.Unix()),
			nodeID,
			ids.ShortEmpty,
			reward.PercentDenominator,
			[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should've errored because start time too early")
		}
	}

	{
		// Case: Validator's start time too far in the future
		tx, err := env.txBuilder.NewAddValidatorTx(
			env.config.MinValidatorStake,
			uint64(defaultValidateStartTime.Add(MaxFutureStartTime).Unix()+1),
			uint64(defaultValidateStartTime.Add(MaxFutureStartTime).Add(defaultMinStakingDuration).Unix()+1),
			nodeID,
			ids.ShortEmpty,
			reward.PercentDenominator,
			[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should've errored because start time too far in the future")
		}
	}

	{
		// Case: Validator already validating primary network
		tx, err := env.txBuilder.NewAddValidatorTx(
			env.config.MinValidatorStake,
			uint64(defaultValidateStartTime.Unix()),
			uint64(defaultValidateEndTime.Unix()),
			nodeID,
			ids.ShortEmpty,
			reward.PercentDenominator,
			[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should've errored because validator already validating")
		}
	}

	{
		// Case: Validator in pending validator set of primary network
		startTime := defaultGenesisTime.Add(1 * time.Second)
		tx, err := env.txBuilder.NewAddValidatorTx(
			env.config.MinValidatorStake,                            // stake amount
			uint64(startTime.Unix()),                                // start time
			uint64(startTime.Add(defaultMinStakingDuration).Unix()), // end time
			nodeID,
			ids.ShortEmpty,
			reward.PercentDenominator, // shares
			[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			ids.ShortEmpty, // change addr // key
		)
		if err != nil {
			t.Fatal(err)
		}

		staker, err := state.NewCurrentStaker(
			tx.ID(),
			tx.Unsigned.(*txs.AddValidatorTx),
			0,
		)
		if err != nil {
			t.Fatal(err)
		}

		env.state.PutCurrentValidator(staker)
		env.state.AddTx(tx, status.Committed)
		dummyHeight := uint64(1)
		env.state.SetHeight(dummyHeight)
		if err := env.state.Commit(); err != nil {
			t.Fatal(err)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed because validator in pending validator set")
		}
	}

	{
		// Case: Validator doesn't have enough tokens to cover stake amount
		tx, err := env.txBuilder.NewAddValidatorTx( // create the tx
			env.config.MinValidatorStake,
			uint64(defaultValidateStartTime.Unix()),
			uint64(defaultValidateEndTime.Unix()),
			nodeID,
			ids.ShortEmpty,
			reward.PercentDenominator,
			[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
			ids.ShortEmpty, // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		// Remove all UTXOs owned by preFundedKeys[0]
		utxoIDs, err := env.state.UTXOIDs(preFundedKeys[0].PublicKey().Address().Bytes(), ids.Empty, math.MaxInt32)
		if err != nil {
			t.Fatal(err)
		}
		for _, utxoID := range utxoIDs {
			env.state.DeleteUTXO(utxoID)
		}

		onAcceptState, err := state.NewDiff(lastAcceptedID, env)
		if err != nil {
			t.Fatal(err)
		}

		executor := StandardTxExecutor{
			Backend: &env.backend,
			State:   onAcceptState,
			Tx:      tx,
		}
		err = tx.Unsigned.Visit(&executor)
		if err == nil {
			t.Fatal("should have failed because tx fee paying key has no funds")
		}
	}
}

// Returns a RemoveSubnetValidatorTx that passes syntactic verification.
func newRemoveSubnetValidatorTx(t *testing.T) (*txs.RemoveSubnetValidatorTx, *txs.Tx) {
	t.Helper()

	creds := []verify.Verifiable{
		&secp256k1fx.Credential{
			Sigs: make([][65]byte, 1),
		},
		&secp256k1fx.Credential{
			Sigs: make([][65]byte, 1),
		},
	}
	unsignedTx := &txs.RemoveSubnetValidatorTx{
		BaseTx: txs.BaseTx{
			BaseTx: djtx.BaseTx{
				Ins: []*djtx.TransferableInput{{
					UTXOID: djtx.UTXOID{
						TxID: ids.GenerateTestID(),
					},
					Asset: djtx.Asset{
						ID: ids.GenerateTestID(),
					},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{0, 1},
						},
					},
				}},
				Outs: []*djtx.TransferableOutput{
					{
						Asset: djtx.Asset{
							ID: ids.GenerateTestID(),
						},
						Out: &secp256k1fx.TransferOutput{
							Amt: 1,
							OutputOwners: secp256k1fx.OutputOwners{
								Threshold: 1,
								Addrs:     []ids.ShortID{ids.GenerateTestShortID()},
							},
						},
					},
				},
				Memo: []byte("hi"),
			},
		},
		Subnet: ids.GenerateTestID(),
		NodeID: ids.GenerateTestNodeID(),
		SubnetAuth: &secp256k1fx.Credential{
			Sigs: make([][65]byte, 1),
		},
	}
	tx := &txs.Tx{
		Unsigned: unsignedTx,
		Creds:    creds,
	}
	if err := tx.Sign(txs.Codec, nil); err != nil {
		t.Fatal(err)
	}
	return unsignedTx, tx
}

// mock implementations that can be used in tests
// for verifying RemoveSubnetValidatorTx.
type removeSubnetValidatorTxVerifyEnv struct {
	banffTime   time.Time
	fx          *fx.MockFx
	flowChecker *utxo.MockVerifier
	unsignedTx  *txs.RemoveSubnetValidatorTx
	tx          *txs.Tx
	state       *state.MockDiff
	staker      *state.Staker
}

// Returns mock implementations that can be used in tests
// for verifying RemoveSubnetValidatorTx.
func newValidRemoveSubnetValidatorTxVerifyEnv(t *testing.T, ctrl *gomock.Controller) removeSubnetValidatorTxVerifyEnv {
	t.Helper()

	now := time.Now()
	mockFx := fx.NewMockFx(ctrl)
	mockFlowChecker := utxo.NewMockVerifier(ctrl)
	unsignedTx, tx := newRemoveSubnetValidatorTx(t)
	mockState := state.NewMockDiff(ctrl)
	return removeSubnetValidatorTxVerifyEnv{
		banffTime:   now,
		fx:          mockFx,
		flowChecker: mockFlowChecker,
		unsignedTx:  unsignedTx,
		tx:          tx,
		state:       mockState,
		staker: &state.Staker{
			TxID:     ids.GenerateTestID(),
			NodeID:   ids.GenerateTestNodeID(),
			Priority: txs.SubnetPermissionedValidatorCurrentPriority,
		},
	}
}

func TestStandardExecutorRemoveSubnetValidatorTx(t *testing.T) {
	type test struct {
		name        string
		newExecutor func(*gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor)
		shouldErr   bool
		expectedErr error
	}

	tests := []test{
		{
			name: "valid tx",
			newExecutor: func(ctrl *gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor) {
				env := newValidRemoveSubnetValidatorTxVerifyEnv(t, ctrl)

				// Set dependency expectations.
				env.state.EXPECT().GetCurrentValidator(env.unsignedTx.Subnet, env.unsignedTx.NodeID).Return(env.staker, nil).Times(1)
				subnetOwner := fx.NewMockOwner(ctrl)
				subnetTx := &txs.Tx{
					Unsigned: &txs.CreateSubnetTx{
						Owner: subnetOwner,
					},
				}
				env.state.EXPECT().GetTx(env.unsignedTx.Subnet).Return(subnetTx, status.Committed, nil).Times(1)
				env.fx.EXPECT().VerifyPermission(env.unsignedTx, env.unsignedTx.SubnetAuth, env.tx.Creds[len(env.tx.Creds)-1], subnetOwner).Return(nil).Times(1)
				env.flowChecker.EXPECT().VerifySpend(
					env.unsignedTx, env.state, env.unsignedTx.Ins, env.unsignedTx.Outs, env.tx.Creds[:len(env.tx.Creds)-1], gomock.Any(),
				).Return(nil).Times(1)
				env.state.EXPECT().DeleteCurrentValidator(env.staker)
				env.state.EXPECT().DeleteUTXO(gomock.Any()).Times(len(env.unsignedTx.Ins))
				env.state.EXPECT().AddUTXO(gomock.Any()).Times(len(env.unsignedTx.Outs))
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			shouldErr: false,
		},
		{
			name: "tx fails syntactic verification",
			newExecutor: func(ctrl *gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor) {
				env := newValidRemoveSubnetValidatorTxVerifyEnv(t, ctrl)
				// Setting the subnet ID to the Primary Network ID makes the tx fail syntactic verification
				env.tx.Unsigned.(*txs.RemoveSubnetValidatorTx).Subnet = constants.PrimaryNetworkID
				env.state = state.NewMockDiff(ctrl)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			shouldErr: true,
		},
		{
			name: "node isn't a validator of the subnet",
			newExecutor: func(ctrl *gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor) {
				env := newValidRemoveSubnetValidatorTxVerifyEnv(t, ctrl)
				env.state = state.NewMockDiff(ctrl)
				env.state.EXPECT().GetCurrentValidator(env.unsignedTx.Subnet, env.unsignedTx.NodeID).Return(nil, database.ErrNotFound)
				env.state.EXPECT().GetPendingValidator(env.unsignedTx.Subnet, env.unsignedTx.NodeID).Return(nil, database.ErrNotFound)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			shouldErr:   true,
			expectedErr: errNotValidator,
		},
		{
			name: "validator is permissionless",
			newExecutor: func(ctrl *gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor) {
				env := newValidRemoveSubnetValidatorTxVerifyEnv(t, ctrl)

				staker := *env.staker
				staker.Priority = txs.SubnetPermissionlessValidatorCurrentPriority

				// Set dependency expectations.
				env.state.EXPECT().GetCurrentValidator(env.unsignedTx.Subnet, env.unsignedTx.NodeID).Return(&staker, nil).Times(1)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			shouldErr:   true,
			expectedErr: errRemovePermissionlessValidator,
		},
		{
			name: "tx has no credentials",
			newExecutor: func(ctrl *gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor) {
				env := newValidRemoveSubnetValidatorTxVerifyEnv(t, ctrl)
				// Remove credentials
				env.tx.Creds = nil
				env.state = state.NewMockDiff(ctrl)
				env.state.EXPECT().GetCurrentValidator(env.unsignedTx.Subnet, env.unsignedTx.NodeID).Return(env.staker, nil)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			shouldErr:   true,
			expectedErr: errWrongNumberOfCredentials,
		},
		{
			name: "can't find subnet",
			newExecutor: func(ctrl *gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor) {
				env := newValidRemoveSubnetValidatorTxVerifyEnv(t, ctrl)
				env.state = state.NewMockDiff(ctrl)
				env.state.EXPECT().GetCurrentValidator(env.unsignedTx.Subnet, env.unsignedTx.NodeID).Return(env.staker, nil)
				env.state.EXPECT().GetTx(env.unsignedTx.Subnet).Return(nil, status.Unknown, database.ErrNotFound)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			shouldErr:   true,
			expectedErr: errCantFindSubnet,
		},
		{
			name: "no permission to remove validator",
			newExecutor: func(ctrl *gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor) {
				env := newValidRemoveSubnetValidatorTxVerifyEnv(t, ctrl)
				env.state = state.NewMockDiff(ctrl)
				env.state.EXPECT().GetCurrentValidator(env.unsignedTx.Subnet, env.unsignedTx.NodeID).Return(env.staker, nil)
				subnetOwner := fx.NewMockOwner(ctrl)
				subnetTx := &txs.Tx{
					Unsigned: &txs.CreateSubnetTx{
						Owner: subnetOwner,
					},
				}
				env.state.EXPECT().GetTx(env.unsignedTx.Subnet).Return(subnetTx, status.Committed, nil)
				env.fx.EXPECT().VerifyPermission(gomock.Any(), env.unsignedTx.SubnetAuth, env.tx.Creds[len(env.tx.Creds)-1], subnetOwner).Return(errors.New(""))
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			shouldErr:   true,
			expectedErr: errUnauthorizedSubnetModification,
		},
		{
			name: "flow checker failed",
			newExecutor: func(ctrl *gomock.Controller) (*txs.RemoveSubnetValidatorTx, *StandardTxExecutor) {
				env := newValidRemoveSubnetValidatorTxVerifyEnv(t, ctrl)
				env.state = state.NewMockDiff(ctrl)
				env.state.EXPECT().GetCurrentValidator(env.unsignedTx.Subnet, env.unsignedTx.NodeID).Return(env.staker, nil)
				subnetOwner := fx.NewMockOwner(ctrl)
				subnetTx := &txs.Tx{
					Unsigned: &txs.CreateSubnetTx{
						Owner: subnetOwner,
					},
				}
				env.state.EXPECT().GetTx(env.unsignedTx.Subnet).Return(subnetTx, status.Committed, nil)
				env.fx.EXPECT().VerifyPermission(gomock.Any(), env.unsignedTx.SubnetAuth, env.tx.Creds[len(env.tx.Creds)-1], subnetOwner).Return(nil)
				env.flowChecker.EXPECT().VerifySpend(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New(""))
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			shouldErr:   true,
			expectedErr: errFlowCheckFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			unsignedTx, executor := tt.newExecutor(ctrl)
			err := executor.RemoveSubnetValidatorTx(unsignedTx)
			if tt.shouldErr {
				require.Error(err)
				if tt.expectedErr != nil {
					require.ErrorIs(err, tt.expectedErr)
				}
				return
			}
			require.NoError(err)
		})
	}
}

// Returns a TransformSubnetTx that passes syntactic verification.
func newTransformSubnetTx(t *testing.T) (*txs.TransformSubnetTx, *txs.Tx) {
	t.Helper()

	creds := []verify.Verifiable{
		&secp256k1fx.Credential{
			Sigs: make([][65]byte, 1),
		},
		&secp256k1fx.Credential{
			Sigs: make([][65]byte, 1),
		},
	}
	unsignedTx := &txs.TransformSubnetTx{
		BaseTx: txs.BaseTx{
			BaseTx: djtx.BaseTx{
				Ins: []*djtx.TransferableInput{{
					UTXOID: djtx.UTXOID{
						TxID: ids.GenerateTestID(),
					},
					Asset: djtx.Asset{
						ID: ids.GenerateTestID(),
					},
					In: &secp256k1fx.TransferInput{
						Amt: 1,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{0, 1},
						},
					},
				}},
				Outs: []*djtx.TransferableOutput{
					{
						Asset: djtx.Asset{
							ID: ids.GenerateTestID(),
						},
						Out: &secp256k1fx.TransferOutput{
							Amt: 1,
							OutputOwners: secp256k1fx.OutputOwners{
								Threshold: 1,
								Addrs:     []ids.ShortID{ids.GenerateTestShortID()},
							},
						},
					},
				},
				Memo: []byte("hi"),
			},
		},
		Subnet:                   ids.GenerateTestID(),
		AssetID:                  ids.GenerateTestID(),
		InitialSupply:            10,
		MaximumSupply:            10,
		MinConsumptionRate:       0,
		MaxConsumptionRate:       reward.PercentDenominator,
		MinValidatorStake:        2,
		MaxValidatorStake:        10,
		MinStakeDuration:         1,
		MaxStakeDuration:         2,
		MinDelegationFee:         reward.PercentDenominator,
		MinDelegatorStake:        1,
		MaxValidatorWeightFactor: 1,
		UptimeRequirement:        reward.PercentDenominator,
		SubnetAuth: &secp256k1fx.Credential{
			Sigs: make([][65]byte, 1),
		},
	}
	tx := &txs.Tx{
		Unsigned: unsignedTx,
		Creds:    creds,
	}
	if err := tx.Sign(txs.Codec, nil); err != nil {
		t.Fatal(err)
	}
	return unsignedTx, tx
}

// mock implementations that can be used in tests
// for verifying TransformSubnetTx.
type transformSubnetTxVerifyEnv struct {
	banffTime   time.Time
	fx          *fx.MockFx
	flowChecker *utxo.MockVerifier
	unsignedTx  *txs.TransformSubnetTx
	tx          *txs.Tx
	state       *state.MockDiff
	staker      *state.Staker
}

// Returns mock implementations that can be used in tests
// for verifying TransformSubnetTx.
func newValidTransformSubnetTxVerifyEnv(t *testing.T, ctrl *gomock.Controller) transformSubnetTxVerifyEnv {
	t.Helper()

	now := time.Now()
	mockFx := fx.NewMockFx(ctrl)
	mockFlowChecker := utxo.NewMockVerifier(ctrl)
	unsignedTx, tx := newTransformSubnetTx(t)
	mockState := state.NewMockDiff(ctrl)
	return transformSubnetTxVerifyEnv{
		banffTime:   now,
		fx:          mockFx,
		flowChecker: mockFlowChecker,
		unsignedTx:  unsignedTx,
		tx:          tx,
		state:       mockState,
		staker: &state.Staker{
			TxID:   ids.GenerateTestID(),
			NodeID: ids.GenerateTestNodeID(),
		},
	}
}

func TestStandardExecutorTransformSubnetTx(t *testing.T) {
	type test struct {
		name        string
		newExecutor func(*gomock.Controller) (*txs.TransformSubnetTx, *StandardTxExecutor)
		err         error
	}

	tests := []test{
		{
			name: "tx fails syntactic verification",
			newExecutor: func(ctrl *gomock.Controller) (*txs.TransformSubnetTx, *StandardTxExecutor) {
				env := newValidTransformSubnetTxVerifyEnv(t, ctrl)
				// Setting the tx to nil makes the tx fail syntactic verification
				env.tx.Unsigned = (*txs.TransformSubnetTx)(nil)
				env.state = state.NewMockDiff(ctrl)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			err: txs.ErrNilTx,
		},
		{
			name: "max stake duration too large",
			newExecutor: func(ctrl *gomock.Controller) (*txs.TransformSubnetTx, *StandardTxExecutor) {
				env := newValidTransformSubnetTxVerifyEnv(t, ctrl)
				env.unsignedTx.MaxStakeDuration = math.MaxUint32
				env.state = state.NewMockDiff(ctrl)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime: env.banffTime,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			err: errMaxStakeDurationTooLarge,
		},
		{
			name: "fail subnet authorization",
			newExecutor: func(ctrl *gomock.Controller) (*txs.TransformSubnetTx, *StandardTxExecutor) {
				env := newValidTransformSubnetTxVerifyEnv(t, ctrl)
				// Remove credentials
				env.tx.Creds = nil
				env.state = state.NewMockDiff(ctrl)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime:        env.banffTime,
							MaxStakeDuration: math.MaxInt64,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			err: errWrongNumberOfCredentials,
		},
		{
			name: "flow checker failed",
			newExecutor: func(ctrl *gomock.Controller) (*txs.TransformSubnetTx, *StandardTxExecutor) {
				env := newValidTransformSubnetTxVerifyEnv(t, ctrl)
				env.state = state.NewMockDiff(ctrl)
				subnetOwner := fx.NewMockOwner(ctrl)
				subnetTx := &txs.Tx{
					Unsigned: &txs.CreateSubnetTx{
						Owner: subnetOwner,
					},
				}
				env.state.EXPECT().GetTx(env.unsignedTx.Subnet).Return(subnetTx, status.Committed, nil)
				env.state.EXPECT().GetSubnetTransformation(env.unsignedTx.Subnet).Return(nil, database.ErrNotFound).Times(1)
				env.fx.EXPECT().VerifyPermission(gomock.Any(), env.unsignedTx.SubnetAuth, env.tx.Creds[len(env.tx.Creds)-1], subnetOwner).Return(nil)
				env.flowChecker.EXPECT().VerifySpend(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errFlowCheckFailed)
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime:        env.banffTime,
							MaxStakeDuration: math.MaxInt64,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			err: errFlowCheckFailed,
		},
		{
			name: "valid tx",
			newExecutor: func(ctrl *gomock.Controller) (*txs.TransformSubnetTx, *StandardTxExecutor) {
				env := newValidTransformSubnetTxVerifyEnv(t, ctrl)

				// Set dependency expectations.
				subnetOwner := fx.NewMockOwner(ctrl)
				subnetTx := &txs.Tx{
					Unsigned: &txs.CreateSubnetTx{
						Owner: subnetOwner,
					},
				}
				env.state.EXPECT().GetTx(env.unsignedTx.Subnet).Return(subnetTx, status.Committed, nil).Times(1)
				env.state.EXPECT().GetSubnetTransformation(env.unsignedTx.Subnet).Return(nil, database.ErrNotFound).Times(1)
				env.fx.EXPECT().VerifyPermission(env.unsignedTx, env.unsignedTx.SubnetAuth, env.tx.Creds[len(env.tx.Creds)-1], subnetOwner).Return(nil).Times(1)
				env.flowChecker.EXPECT().VerifySpend(
					env.unsignedTx, env.state, env.unsignedTx.Ins, env.unsignedTx.Outs, env.tx.Creds[:len(env.tx.Creds)-1], gomock.Any(),
				).Return(nil).Times(1)
				env.state.EXPECT().AddSubnetTransformation(env.tx)
				env.state.EXPECT().SetCurrentSupply(env.unsignedTx.Subnet, env.unsignedTx.InitialSupply)
				env.state.EXPECT().DeleteUTXO(gomock.Any()).Times(len(env.unsignedTx.Ins))
				env.state.EXPECT().AddUTXO(gomock.Any()).Times(len(env.unsignedTx.Outs))
				e := &StandardTxExecutor{
					Backend: &Backend{
						Config: &config.Config{
							BanffTime:        env.banffTime,
							MaxStakeDuration: math.MaxInt64,
						},
						Bootstrapped: &utils.AtomicBool{},
						Fx:           env.fx,
						FlowChecker:  env.flowChecker,
						Ctx:          &snow.Context{},
					},
					Tx:    env.tx,
					State: env.state,
				}
				e.Bootstrapped.SetValue(true)
				return env.unsignedTx, e
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			unsignedTx, executor := tt.newExecutor(ctrl)
			err := executor.TransformSubnetTx(unsignedTx)
			require.ErrorIs(err, tt.err)
		})
	}
}