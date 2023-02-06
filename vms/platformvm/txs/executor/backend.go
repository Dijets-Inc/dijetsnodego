// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/snow/uptime"
	"github.com/lasthyphen/dijetsnodego/utils"
	"github.com/lasthyphen/dijetsnodego/utils/timer/mockable"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/config"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/fx"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/reward"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/utxo"
)

type Backend struct {
	Config       *config.Config
	Ctx          *snow.Context
	Clk          *mockable.Clock
	Fx           fx.Fx
	FlowChecker  utxo.Verifier
	Uptimes      uptime.Manager
	Rewards      reward.Calculator
	Bootstrapped *utils.AtomicBool
}
