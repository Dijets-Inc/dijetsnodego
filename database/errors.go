// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package database

import "errors"

// common errors
var (
	ErrClosed   = errors.New("closed")
	ErrNotFound = errors.New("not found")
)
