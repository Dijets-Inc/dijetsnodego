// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import "github.com/lasthyphen/dijetsnodego/ids"

// ValidatorID represents a validator that we gossip to other peers
type ValidatorID struct {
	// The validator's ID
	NodeID ids.NodeID
	// The Tx that added this into the validator set
	TxID ids.ID
}
