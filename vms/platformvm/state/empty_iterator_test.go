// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyIterator(t *testing.T) {
	require := require.New(t)
	require.False(EmptyIterator.Next())

	EmptyIterator.Release()

	require.False(EmptyIterator.Next())
	require.Nil(EmptyIterator.Value())
}
