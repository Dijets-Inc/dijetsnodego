// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

// Allow vm to execute custom logic against the underlying transaction types.
type Visitor interface {
	AddValidatorTx(*AddValidatorTx) error
	AddSubnetValidatorTx(*AddSubnetValidatorTx) error
	AddDelegatorTx(*AddDelegatorTx) error
	CreateChainTx(*CreateChainTx) error
	CreateSubnetTx(*CreateSubnetTx) error
	ImportTx(*ImportTx) error
	ExportTx(*ExportTx) error
	AdvanceTimeTx(*AdvanceTimeTx) error
	RewardValidatorTx(*RewardValidatorTx) error
	RemoveSubnetValidatorTx(*RemoveSubnetValidatorTx) error
	TransformSubnetTx(*TransformSubnetTx) error
	AddPermissionlessValidatorTx(*AddPermissionlessValidatorTx) error
	AddPermissionlessDelegatorTx(*AddPermissionlessDelegatorTx) error
}
