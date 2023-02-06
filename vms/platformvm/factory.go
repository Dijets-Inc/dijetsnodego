// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"github.com/lasthyphen/dijetsnodego/snow"
	"github.com/lasthyphen/dijetsnodego/vms"
	"github.com/lasthyphen/dijetsnodego/vms/platformvm/config"
)

var _ vms.Factory = (*Factory)(nil)

// Factory can create new instances of the Platform Chain
type Factory struct {
	config.Config
}

// New returns a new instance of the Platform Chain
func (f *Factory) New(*snow.Context) (interface{}, error) {
	return &VM{Factory: *f}, nil
}
