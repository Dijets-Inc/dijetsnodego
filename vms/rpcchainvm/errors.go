// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rpcchainvm

import (
	"github.com/lasthyphen/dijetsnodego/database"
	"github.com/lasthyphen/dijetsnodego/snow/engine/snowman/block"
)

var (
	errCodeToError = map[uint32]error{
		1: database.ErrClosed,
		2: database.ErrNotFound,
		3: block.ErrHeightIndexedVMNotImplemented,
		4: block.ErrIndexIncomplete,
		5: block.ErrStateSyncableVMNotImplemented,
	}
	errorToErrCode = map[error]uint32{
		database.ErrClosed:                     1,
		database.ErrNotFound:                   2,
		block.ErrHeightIndexedVMNotImplemented: 3,
		block.ErrIndexIncomplete:               4,
		block.ErrStateSyncableVMNotImplemented: 5,
	}
)

func errorToRPCError(err error) error {
	if _, ok := errorToErrCode[err]; ok {
		return nil
	}
	return err
}
