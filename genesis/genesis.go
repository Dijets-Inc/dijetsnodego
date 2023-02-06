// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/utils/constants"
	"github.com/lasthyphen/dijetsnodego/utils/formatting"
	"github.com/lasthyphen/dijetsnodego/utils/formatting/address"
	"github.com/lasthyphen/dijetsnodego/utils/json"
	"github.com/lasthyphen/dijetsnodego/utils/set"
	"github.com/lasthyphen/dijetsnodego/vms/avm"
	"github.com/lasthyphen/dijetsnodego/vms/avm/fxs"
	"github.com/lasthyphen/dijetsnodego/vms/nftfx"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/api"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/genesis"
	"github.com/lasthyphen/dijetsnodego/vms/propertyfx"
	"github.com/lasthyphen/dijetsnodego/vms/secp256k1fx"

	xchaintxs "github.com/lasthyphen/dijetsnodego/vms/avm/txs"
	pchaintxs "github.com/lasthyphen/dijetsnodego/vms/platformvm/txs"
)

const (
	defaultEncoding    = formatting.Hex
	configChainIDAlias = "X"
)

var (
	errNoInitiallyStakedFunds = errors.New("initial staked funds cannot be empty")
	errNoSupply               = errors.New("initial supply must be > 0")
	errNoStakeDuration        = errors.New("initial stake duration must be > 0")
	errNoStakers              = errors.New("initial stakers must be > 0")
	errNoCChainGenesis        = errors.New("C-Chain genesis cannot be empty")
	errNoTxs                  = errors.New("genesis creates no transactions")
)

// validateInitialStakedFunds ensures all staked
// funds have allocations and that all staked
// funds are unique.
//
// This function assumes that NetworkID in *Config has already
// been checked for correctness.
func validateInitialStakedFunds(config *Config) error {
	if len(config.InitialStakedFunds) == 0 {
		return errNoInitiallyStakedFunds
	}

	allocationSet := set.Set[ids.ShortID]{}
	initialStakedFundsSet := set.Set[ids.ShortID]{}
	for _, allocation := range config.Allocations {
		// It is ok to have duplicates as different
		// ethAddrs could claim to the same djtxAddr.
		allocationSet.Add(allocation.DJTXAddr)
	}

	for _, staker := range config.InitialStakedFunds {
		if initialStakedFundsSet.Contains(staker) {
			djtxAddr, err := address.Format(
				configChainIDAlias,
				constants.GetHRP(config.NetworkID),
				staker.Bytes(),
			)
			if err != nil {
				return fmt.Errorf(
					"unable to format address from %s",
					staker.String(),
				)
			}

			return fmt.Errorf(
				"address %s is duplicated in initial staked funds",
				djtxAddr,
			)
		}
		initialStakedFundsSet.Add(staker)

		if !allocationSet.Contains(staker) {
			djtxAddr, err := address.Format(
				configChainIDAlias,
				constants.GetHRP(config.NetworkID),
				staker.Bytes(),
			)
			if err != nil {
				return fmt.Errorf(
					"unable to format address from %s",
					staker.String(),
				)
			}

			return fmt.Errorf(
				"address %s does not have an allocation to stake",
				djtxAddr,
			)
		}
	}

	return nil
}

// validateConfig returns an error if the provided
// *Config is not considered valid.
func validateConfig(networkID uint32, config *Config) error {
	if networkID != config.NetworkID {
		return fmt.Errorf(
			"networkID %d specified but genesis config contains networkID %d",
			networkID,
			config.NetworkID,
		)
	}

	initialSupply, err := config.InitialSupply()
	switch {
	case err != nil:
		return fmt.Errorf("unable to calculate initial supply: %w", err)
	case initialSupply == 0:
		return errNoSupply
	}

	startTime := time.Unix(int64(config.StartTime), 0)
	if time.Since(startTime) < 0 {
		return fmt.Errorf(
			"start time cannot be in the future: %s",
			startTime,
		)
	}

	// We don't impose any restrictions on the minimum
	// stake duration to enable complex testing configurations
	// but recommend setting a minimum duration of at least
	// 15 minutes.
	if config.InitialStakeDuration == 0 {
		return errNoStakeDuration
	}

	if len(config.InitialStakers) == 0 {
		return errNoStakers
	}

	offsetTimeRequired := config.InitialStakeDurationOffset * uint64(len(config.InitialStakers)-1)
	if offsetTimeRequired > config.InitialStakeDuration {
		return fmt.Errorf(
			"initial stake duration is %d but need at least %d with offset of %d",
			config.InitialStakeDuration,
			offsetTimeRequired,
			config.InitialStakeDurationOffset,
		)
	}

	if err := validateInitialStakedFunds(config); err != nil {
		return fmt.Errorf("initial staked funds validation failed: %w", err)
	}

	if len(config.CChainGenesis) == 0 {
		return errNoCChainGenesis
	}

	return nil
}

// FromFile returns the genesis data of the Platform Chain.
//
// Since an Avalanche network has exactly one Platform Chain, and the Platform
// Chain defines the genesis state of the network (who is staking, which chains
// exist, etc.), defining the genesis state of the Platform Chain is the same as
// defining the genesis state of the network.
//
// FromFile accepts:
// 1) The ID of the new network. [networkID]
// 2) The location of a custom genesis config to load. [filepath]
//
// If [filepath] is empty or the given network ID is Mainnet, Testnet, or Local, returns error.
// If [filepath] is non-empty and networkID isn't Mainnet, Testnet, or Local,
// loads the network genesis data from the config at [filepath].
//
// FromFile returns:
//  1. The byte representation of the genesis state of the platform chain
//     (ie the genesis state of the network)
//  2. The asset ID of DJTX
func FromFile(networkID uint32, filepath string) ([]byte, ids.ID, error) {
	switch networkID {
	case constants.MainnetID, constants.TahoeID, constants.LocalID:
		return nil, ids.ID{}, fmt.Errorf(
			"cannot override genesis config for standard network %s (%d)",
			constants.NetworkName(networkID),
			networkID,
		)
	}

	config, err := GetConfigFile(filepath)
	if err != nil {
		return nil, ids.ID{}, fmt.Errorf("unable to load provided genesis config at %s: %w", filepath, err)
	}

	if err := validateConfig(networkID, config); err != nil {
		return nil, ids.ID{}, fmt.Errorf("genesis config validation failed: %w", err)
	}

	return FromConfig(config)
}

// FromFlag returns the genesis data of the Platform Chain.
//
// Since an Avalanche network has exactly one Platform Chain, and the Platform
// Chain defines the genesis state of the network (who is staking, which chains
// exist, etc.), defining the genesis state of the Platform Chain is the same as
// defining the genesis state of the network.
//
// FromFlag accepts:
// 1) The ID of the new network. [networkID]
// 2) The content of a custom genesis config to load. [genesisContent]
//
// If [genesisContent] is empty or the given network ID is Mainnet, Testnet, or Local, returns error.
// If [genesisContent] is non-empty and networkID isn't Mainnet, Testnet, or Local,
// loads the network genesis data from [genesisContent].
//
// FromFlag returns:
//  1. The byte representation of the genesis state of the platform chain
//     (ie the genesis state of the network)
//  2. The asset ID of DJTX
func FromFlag(networkID uint32, genesisContent string) ([]byte, ids.ID, error) {
	switch networkID {
	case constants.MainnetID, constants.TahoeID, constants.LocalID:
		return nil, ids.ID{}, fmt.Errorf(
			"cannot override genesis config for standard network %s (%d)",
			constants.NetworkName(networkID),
			networkID,
		)
	}

	customConfig, err := GetConfigContent(genesisContent)
	if err != nil {
		return nil, ids.ID{}, fmt.Errorf("unable to load genesis content from flag: %w", err)
	}

	if err := validateConfig(networkID, customConfig); err != nil {
		return nil, ids.ID{}, fmt.Errorf("genesis config validation failed: %w", err)
	}

	return FromConfig(customConfig)
}

// FromConfig returns:
//  1. The byte representation of the genesis state of the platform chain
//     (ie the genesis state of the network)
//  2. The asset ID of DJTX
func FromConfig(config *Config) ([]byte, ids.ID, error) {
	hrp := constants.GetHRP(config.NetworkID)

	amount := uint64(0)

	// Specify the genesis state of the AVM
	avmArgs := avm.BuildGenesisArgs{
		NetworkID: json.Uint32(config.NetworkID),
		Encoding:  defaultEncoding,
	}
	{
		djtx := avm.AssetDefinition{
			Name:         "Dijets",
			Symbol:       "DJTX",
			Denomination: 9,
			InitialState: map[string][]interface{}{},
		}
		memoBytes := []byte{}
		xAllocations := []Allocation(nil)
		for _, allocation := range config.Allocations {
			if allocation.InitialAmount > 0 {
				xAllocations = append(xAllocations, allocation)
			}
		}
		sortXAllocation(xAllocations)

		for _, allocation := range xAllocations {
			addr, err := address.FormatBech32(hrp, allocation.DJTXAddr.Bytes())
			if err != nil {
				return nil, ids.ID{}, err
			}

			djtx.InitialState["fixedCap"] = append(djtx.InitialState["fixedCap"], avm.Holder{
				Amount:  json.Uint64(allocation.InitialAmount),
				Address: addr,
			})
			memoBytes = append(memoBytes, allocation.ETHAddr.Bytes()...)
			amount += allocation.InitialAmount
		}

		var err error
		djtx.Memo, err = formatting.Encode(defaultEncoding, memoBytes)
		if err != nil {
			return nil, ids.Empty, fmt.Errorf("couldn't parse memo bytes to string: %w", err)
		}
		avmArgs.GenesisData = map[string]avm.AssetDefinition{
			"DJTX": djtx, // The AVM starts out with one asset: DJTX
		}
	}
	avmReply := avm.BuildGenesisReply{}

	avmSS := avm.CreateStaticService()
	err := avmSS.BuildGenesis(nil, &avmArgs, &avmReply)
	if err != nil {
		return nil, ids.ID{}, err
	}

	bytes, err := formatting.Decode(defaultEncoding, avmReply.Bytes)
	if err != nil {
		return nil, ids.ID{}, fmt.Errorf("couldn't parse avm genesis reply: %w", err)
	}
	djtxAssetID, err := DJTXAssetID(bytes)
	if err != nil {
		return nil, ids.ID{}, fmt.Errorf("couldn't generate DJTX asset ID: %w", err)
	}

	genesisTime := time.Unix(int64(config.StartTime), 0)
	initialSupply, err := config.InitialSupply()
	if err != nil {
		return nil, ids.ID{}, fmt.Errorf("couldn't calculate the initial supply: %w", err)
	}

	initiallyStaked := set.Set[ids.ShortID]{}
	initiallyStaked.Add(config.InitialStakedFunds...)
	skippedAllocations := []Allocation(nil)

	// Specify the initial state of the Platform Chain
	platformvmArgs := api.BuildGenesisArgs{
		DjtxAssetID:   djtxAssetID,
		NetworkID:     json.Uint32(config.NetworkID),
		Time:          json.Uint64(config.StartTime),
		InitialSupply: json.Uint64(initialSupply),
		Message:       config.Message,
		Encoding:      defaultEncoding,
	}
	for _, allocation := range config.Allocations {
		if initiallyStaked.Contains(allocation.DJTXAddr) {
			skippedAllocations = append(skippedAllocations, allocation)
			continue
		}
		addr, err := address.FormatBech32(hrp, allocation.DJTXAddr.Bytes())
		if err != nil {
			return nil, ids.ID{}, err
		}
		for _, unlock := range allocation.UnlockSchedule {
			if unlock.Amount > 0 {
				msgStr, err := formatting.Encode(defaultEncoding, allocation.ETHAddr.Bytes())
				if err != nil {
					return nil, ids.Empty, fmt.Errorf("couldn't encode message: %w", err)
				}
				platformvmArgs.UTXOs = append(platformvmArgs.UTXOs,
					api.UTXO{
						Locktime: json.Uint64(unlock.Locktime),
						Amount:   json.Uint64(unlock.Amount),
						Address:  addr,
						Message:  msgStr,
					},
				)
				amount += unlock.Amount
			}
		}
	}

	allNodeAllocations := splitAllocations(skippedAllocations, len(config.InitialStakers))
	endStakingTime := genesisTime.Add(time.Duration(config.InitialStakeDuration) * time.Second)
	stakingOffset := time.Duration(0)
	for i, staker := range config.InitialStakers {
		nodeAllocations := allNodeAllocations[i]
		endStakingTime := endStakingTime.Add(-stakingOffset)
		stakingOffset += time.Duration(config.InitialStakeDurationOffset) * time.Second

		destAddrStr, err := address.FormatBech32(hrp, staker.RewardAddress.Bytes())
		if err != nil {
			return nil, ids.ID{}, err
		}

		utxos := []api.UTXO(nil)
		for _, allocation := range nodeAllocations {
			addr, err := address.FormatBech32(hrp, allocation.DJTXAddr.Bytes())
			if err != nil {
				return nil, ids.ID{}, err
			}
			for _, unlock := range allocation.UnlockSchedule {
				msgStr, err := formatting.Encode(defaultEncoding, allocation.ETHAddr.Bytes())
				if err != nil {
					return nil, ids.Empty, fmt.Errorf("couldn't encode message: %w", err)
				}
				utxos = append(utxos, api.UTXO{
					Locktime: json.Uint64(unlock.Locktime),
					Amount:   json.Uint64(unlock.Amount),
					Address:  addr,
					Message:  msgStr,
				})
				amount += unlock.Amount
			}
		}

		delegationFee := json.Uint32(staker.DelegationFee)

		platformvmArgs.Validators = append(platformvmArgs.Validators,
			api.PermissionlessValidator{
				Staker: api.Staker{
					StartTime: json.Uint64(genesisTime.Unix()),
					EndTime:   json.Uint64(endStakingTime.Unix()),
					NodeID:    staker.NodeID,
				},
				RewardOwner: &api.Owner{
					Threshold: 1,
					Addresses: []string{destAddrStr},
				},
				Staked:             utxos,
				ExactDelegationFee: &delegationFee,
			},
		)
	}

	// Specify the chains that exist upon this network's creation
	genesisStr, err := formatting.Encode(defaultEncoding, []byte(config.CChainGenesis))
	if err != nil {
		return nil, ids.Empty, fmt.Errorf("couldn't encode message: %w", err)
	}
	platformvmArgs.Chains = []api.Chain{
		{
			GenesisData: avmReply.Bytes,
			SubnetID:    constants.PrimaryNetworkID,
			VMID:        constants.AVMID,
			FxIDs: []ids.ID{
				secp256k1fx.ID,
				nftfx.ID,
				propertyfx.ID,
			},
			Name: "X-Chain",
		},
		{
			GenesisData: genesisStr,
			SubnetID:    constants.PrimaryNetworkID,
			VMID:        constants.EVMID,
			Name:        "C-Chain",
		},
	}

	platformvmReply := api.BuildGenesisReply{}
	platformvmSS := api.StaticService{}
	if err := platformvmSS.BuildGenesis(nil, &platformvmArgs, &platformvmReply); err != nil {
		return nil, ids.ID{}, fmt.Errorf("problem while building platform chain's genesis state: %w", err)
	}

	genesisBytes, err := formatting.Decode(platformvmReply.Encoding, platformvmReply.Bytes)
	if err != nil {
		return nil, ids.ID{}, fmt.Errorf("problem parsing platformvm genesis bytes: %w", err)
	}

	return genesisBytes, djtxAssetID, nil
}

func splitAllocations(allocations []Allocation, numSplits int) [][]Allocation {
	totalAmount := uint64(0)
	for _, allocation := range allocations {
		for _, unlock := range allocation.UnlockSchedule {
			totalAmount += unlock.Amount
		}
	}

	nodeWeight := totalAmount / uint64(numSplits)
	allNodeAllocations := make([][]Allocation, 0, numSplits)

	currentNodeAllocation := []Allocation(nil)
	currentNodeAmount := uint64(0)
	for _, allocation := range allocations {
		currentAllocation := allocation
		// Already added to the X-chain
		currentAllocation.InitialAmount = 0
		// Going to be added until the correct amount is reached
		currentAllocation.UnlockSchedule = nil

		for _, unlock := range allocation.UnlockSchedule {
			unlock := unlock
			for currentNodeAmount+unlock.Amount > nodeWeight && len(allNodeAllocations) < numSplits-1 {
				amountToAdd := nodeWeight - currentNodeAmount
				currentAllocation.UnlockSchedule = append(currentAllocation.UnlockSchedule, LockedAmount{
					Amount:   amountToAdd,
					Locktime: unlock.Locktime,
				})
				unlock.Amount -= amountToAdd

				currentNodeAllocation = append(currentNodeAllocation, currentAllocation)

				allNodeAllocations = append(allNodeAllocations, currentNodeAllocation)

				currentNodeAllocation = nil
				currentNodeAmount = 0

				currentAllocation = allocation
				// Already added to the X-chain
				currentAllocation.InitialAmount = 0
				// Going to be added until the correct amount is reached
				currentAllocation.UnlockSchedule = nil
			}

			if unlock.Amount == 0 {
				continue
			}

			currentAllocation.UnlockSchedule = append(currentAllocation.UnlockSchedule, LockedAmount{
				Amount:   unlock.Amount,
				Locktime: unlock.Locktime,
			})
			currentNodeAmount += unlock.Amount
		}

		if len(currentAllocation.UnlockSchedule) > 0 {
			currentNodeAllocation = append(currentNodeAllocation, currentAllocation)
		}
	}

	return append(allNodeAllocations, currentNodeAllocation)
}

func VMGenesis(genesisBytes []byte, vmID ids.ID) (*pchaintxs.Tx, error) {
	genesis, err := genesis.Parse(genesisBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse genesis: %w", err)
	}
	for _, chain := range genesis.Chains {
		uChain := chain.Unsigned.(*pchaintxs.CreateChainTx)
		if uChain.VMID == vmID {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("couldn't find blockchain with VM ID %s", vmID)
}

func DJTXAssetID(avmGenesisBytes []byte) (ids.ID, error) {
	parser, err := xchaintxs.NewParser([]fxs.Fx{
		&secp256k1fx.Fx{},
	})
	if err != nil {
		return ids.Empty, err
	}

	genesisCodec := parser.GenesisCodec()
	genesis := avm.Genesis{}
	if _, err := genesisCodec.Unmarshal(avmGenesisBytes, &genesis); err != nil {
		return ids.Empty, err
	}

	if len(genesis.Txs) == 0 {
		return ids.Empty, errNoTxs
	}
	genesisTx := genesis.Txs[0]

	tx := xchaintxs.Tx{Unsigned: &genesisTx.CreateAssetTx}
	if err := parser.InitializeGenesisTx(&tx); err != nil {
		return ids.Empty, err
	}
	return tx.ID(), nil
}

type innerSortXAllocation []Allocation

func (xa innerSortXAllocation) Less(i, j int) bool {
	return xa[i].InitialAmount < xa[j].InitialAmount ||
		(xa[i].InitialAmount == xa[j].InitialAmount &&
			bytes.Compare(xa[i].DJTXAddr.Bytes(), xa[j].DJTXAddr.Bytes()) == -1)
}

func (xa innerSortXAllocation) Len() int      { return len(xa) }
func (xa innerSortXAllocation) Swap(i, j int) { xa[j], xa[i] = xa[i], xa[j] }

func sortXAllocation(a []Allocation) { sort.Sort(innerSortXAllocation(a)) }