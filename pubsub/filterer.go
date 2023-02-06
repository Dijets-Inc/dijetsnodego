// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package pubsub

type Filterer interface {
	Filter(connections []Filter) ([]bool, interface{})
}
