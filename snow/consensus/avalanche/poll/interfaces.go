// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package poll

import (
	"fmt"

	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/utils/formatting"
)

// Set is a collection of polls
type Set interface {
	fmt.Stringer

	Add(requestID uint32, vdrs ids.NodeIDBag) bool
	Vote(requestID uint32, vdr ids.NodeID, votes []ids.ID) []ids.UniqueBag
	Len() int
}

// Poll is an outstanding poll
type Poll interface {
	formatting.PrefixedStringer

	Vote(vdr ids.NodeID, votes []ids.ID)
	Finished() bool
	Result() ids.UniqueBag
}

// Factory creates a new Poll
type Factory interface {
	New(vdrs ids.NodeIDBag) Poll
}
