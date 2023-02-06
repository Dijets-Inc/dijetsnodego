// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package codec

import "github.com/lasthyphen/dijetsnodego/utils/wrappers"

// Codec marshals and unmarshals
type Codec interface {
	MarshalInto(interface{}, *wrappers.Packer) error
	Unmarshal([]byte, interface{}) error
}
