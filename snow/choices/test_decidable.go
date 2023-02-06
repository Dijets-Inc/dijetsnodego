// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package choices

import (
	"context"
	"fmt"

	"github.com/lasthyphen/dijetsnodego/ids"
)

var _ Decidable = (*TestDecidable)(nil)

// TestDecidable is a test Decidable
type TestDecidable struct {
	IDV              ids.ID
	AcceptV, RejectV error
	StatusV          Status
}

func (d *TestDecidable) ID() ids.ID {
	return d.IDV
}

func (d *TestDecidable) Accept(context.Context) error {
	switch d.StatusV {
	case Unknown, Rejected:
		return fmt.Errorf("invalid state transition from %s to %s",
			d.StatusV, Accepted)
	default:
		d.StatusV = Accepted
		return d.AcceptV
	}
}

func (d *TestDecidable) Reject(context.Context) error {
	switch d.StatusV {
	case Unknown, Accepted:
		return fmt.Errorf("invalid state transition from %s to %s",
			d.StatusV, Rejected)
	default:
		d.StatusV = Rejected
		return d.RejectV
	}
}

func (d *TestDecidable) Status() Status {
	return d.StatusV
}
