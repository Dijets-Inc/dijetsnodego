// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"github.com/lasthyphen/dijetsnodego/codec"
	"github.com/lasthyphen/dijetsnodego/codec/linearcodec"
	"github.com/lasthyphen/dijetsnodego/utils/units"
	"github.com/lasthyphen/dijetsnodego/utils/wrappers"
)

const (
	codecVersion   uint16 = 0
	maxMessageSize        = 512 * units.KiB
	maxSliceLen           = maxMessageSize
)

// Codec does serialization and deserialization
var c codec.Manager

func init() {
	c = codec.NewManager(maxMessageSize)
	lc := linearcodec.NewCustomMaxLength(maxSliceLen)

	errs := wrappers.Errs{}
	errs.Add(
		lc.RegisterType(&Tx{}),
		c.RegisterCodec(codecVersion, lc),
	)
	if errs.Errored() {
		panic(errs.Err)
	}
}
