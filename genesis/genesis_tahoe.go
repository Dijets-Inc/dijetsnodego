// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	_ "embed"

	"github.com/lasthyphen/dijetsnodego/utils/units"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/reward"
)

var (
	//go:embed genesis_tahoe.json
	tahoeGenesisConfigJSON []byte

	// TahoeParams are the params used for the tahoe testnet
	TahoeParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                         units.MilliDjtx,
			CreateAssetTxFee:              10 * units.MilliDjtx,
			CreateSubnetTxFee:             100 * units.MilliDjtx,
			TransformSubnetTxFee:          100 * units.MilliDjtx,
			CreateBlockchainTxFee:         100 * units.MilliDjtx,
			AddPrimaryNetworkValidatorFee: 0,
			AddPrimaryNetworkDelegatorFee: 0,
			AddSubnetValidatorFee:         units.MilliDjtx,
			AddSubnetDelegatorFee:         units.MilliDjtx,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 1 * units.Djtx,
			MaxValidatorStake: 3 * units.MegaDjtx,
			MinDelegatorStake: 1 * units.Djtx,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  24 * time.Hour,
			MaxStakeDuration:  365 * 24 * time.Hour,
			RewardConfig: reward.Config{
				MaxConsumptionRate: .12 * reward.PercentDenominator,
				MinConsumptionRate: .10 * reward.PercentDenominator,
				MintingPeriod:      365 * 24 * time.Hour,
				SupplyCap:          666666666 * units.Djtx,
			},
		},
	}
)
