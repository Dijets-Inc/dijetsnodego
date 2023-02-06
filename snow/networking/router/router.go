// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package router

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lasthyphen/dijetsnodego/api/health"
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/message"
	"github.com/lasthyphen/dijetsnodego/snow/networking/benchlist"
	"github.com/lasthyphen/dijetsnodego/snow/networking/handler"
	"github.com/lasthyphen/dijetsnodego/snow/networking/timeout"
	"github.com/lasthyphen/dijetsnodego/utils/logging"
	"github.com/lasthyphen/dijetsnodego/utils/set"
)

// Router routes consensus messages to the Handler of the consensus
// engine that the messages are intended for
type Router interface {
	ExternalHandler
	InternalHandler

	Initialize(
		nodeID ids.NodeID,
		log logging.Logger,
		timeouts timeout.Manager,
		shutdownTimeout time.Duration,
		criticalChains set.Set[ids.ID],
		whiteListedSubnets set.Set[ids.ID],
		onFatal func(exitCode int),
		healthConfig HealthConfig,
		metricsNamespace string,
		metricsRegisterer prometheus.Registerer,
	) error
	Shutdown(context.Context)
	AddChain(ctx context.Context, chain handler.Handler)
	health.Checker
}

// InternalHandler deals with messages internal to this node
type InternalHandler interface {
	benchlist.Benchable

	RegisterRequest(
		ctx context.Context,
		nodeID ids.NodeID,
		sourceChainID ids.ID,
		destinationChainID ids.ID,
		requestID uint32,
		op message.Op,
		failedMsg message.InboundMessage,
	)
}
