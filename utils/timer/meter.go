// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package timer

// Meter tracks the number of occurrences of a specified event
type Meter interface {
	// Notify this meter of a new event for it to rate
	Tick()
	// Return the number of events this meter is currently tracking
	Ticks() int
}
