// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"time"

	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/message"
	"github.com/lasthyphen/dijetsnodego/network/throttling"
	"github.com/lasthyphen/dijetsnodego/snow/networking/router"
	"github.com/lasthyphen/dijetsnodego/snow/networking/tracker"
	"github.com/lasthyphen/dijetsnodego/snow/uptime"
	"github.com/lasthyphen/dijetsnodego/snow/validators"
	"github.com/lasthyphen/dijetsnodego/utils/logging"
	"github.com/lasthyphen/dijetsnodego/utils/set"
	"github.com/lasthyphen/dijetsnodego/utils/timer/mockable"
	"github.com/lasthyphen/dijetsnodego/version"
)

type Config struct {
	// Size, in bytes, of the buffer this peer reads messages into
	ReadBufferSize int
	// Size, in bytes, of the buffer this peer writes messages into
	WriteBufferSize int
	Clock           mockable.Clock
	Metrics         *Metrics
	MessageCreator  message.Creator

	Log                  logging.Logger
	InboundMsgThrottler  throttling.InboundMsgThrottler
	Network              Network
	Router               router.InboundHandler
	VersionCompatibility version.Compatibility
	MySubnets            set.Set[ids.ID]
	Beacons              validators.Set
	NetworkID            uint32
	PingFrequency        time.Duration
	PongTimeout          time.Duration
	MaxClockDifference   time.Duration

	// Unix time of the last message sent and received respectively
	// Must only be accessed atomically
	LastSent, LastReceived int64

	// Tracks CPU/disk usage caused by each peer.
	ResourceTracker tracker.ResourceTracker

	// Tracks which peer knows about which peers
	GossipTracker GossipTracker

	// Calculates uptime of peers
	UptimeCalculator uptime.Calculator

	// Signs my IP so I can send my signed IP address in the Version message
	IPSigner *IPSigner
}
