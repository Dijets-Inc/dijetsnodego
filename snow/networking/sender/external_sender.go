// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sender

import (
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/message"
	"github.com/lasthyphen/dijetsnodego/utils/set"
)

// ExternalSender sends consensus messages to other validators
// Right now this is implemented in the networking package
type ExternalSender interface {
	// Send a message to a specific set of nodes
	Send(
		msg message.OutboundMessage,
		nodeIDs set.Set[ids.NodeID],
		subnetID ids.ID,
		validatorOnly bool,
	) set.Set[ids.NodeID]

	// Send a message to a random group of nodes in a subnet.
	// Nodes are sampled based on their validator status.
	Gossip(
		msg message.OutboundMessage,
		subnetID ids.ID,
		validatorOnly bool,
		numValidatorsToSend int,
		numNonValidatorsToSend int,
		numPeersToSend int,
	) set.Set[ids.NodeID]
}
