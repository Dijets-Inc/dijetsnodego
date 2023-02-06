// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package summary

import (
	"errors"
	"math"

	"github.com/lasthyphen/dijetsnodego/codec"
	"github.com/lasthyphen/dijetsnodego/codec/linearcodec"
)

const codecVersion = 0

var (
	c codec.Manager

	errWrongCodecVersion = errors.New("wrong codec version")
)

func init() {
	lc := linearcodec.NewCustomMaxLength(math.MaxUint32)
	c = codec.NewManager(math.MaxInt32)
	if err := c.RegisterCodec(codecVersion, lc); err != nil {
		panic(err)
	}
}
