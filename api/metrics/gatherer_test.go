// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package metrics

import (
	dto "github.com/prometheus/client_model/go"
)

var (
	hello      = "hello"
	world      = "world"
	helloWorld = "hello_world"
)

type testGatherer struct {
	mfs []*dto.MetricFamily
	err error
}

func (g *testGatherer) Gather() ([]*dto.MetricFamily, error) {
	return g.mfs, g.err
}
