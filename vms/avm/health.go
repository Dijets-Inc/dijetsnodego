// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import "context"

// TODO: add health checks
func (*VM) HealthCheck(context.Context) (interface{}, error) {
	return nil, nil
}
