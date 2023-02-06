// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package secp256k1fx

import (
	"github.com/lasthyphen/dijetsnodego/codec"
	"github.com/lasthyphen/dijetsnodego/utils/logging"
	"github.com/lasthyphen/dijetsnodego/utils/timer/mockable"
)

// VM that this Fx must be run by
type VM interface {
	CodecRegistry() codec.Registry
	Clock() *mockable.Clock
	Logger() logging.Logger
}

var _ VM = (*TestVM)(nil)

// TestVM is a minimal implementation of a VM
type TestVM struct {
	Clk   mockable.Clock
	Codec codec.Registry
	Log   logging.Logger
}

func (vm *TestVM) Clock() *mockable.Clock {
	return &vm.Clk
}

func (vm *TestVM) CodecRegistry() codec.Registry {
	return vm.Codec
}

func (vm *TestVM) Logger() logging.Logger {
	return vm.Log
}
