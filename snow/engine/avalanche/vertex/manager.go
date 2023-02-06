// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vertex

// Manager defines all the vertex related functionality that is required by the
// consensus engine.
type Manager interface {
	Builder
	Parser
	Storage
}
