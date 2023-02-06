// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

// Factory returns new instances of Consensus
type Factory interface {
	New() Consensus
}
