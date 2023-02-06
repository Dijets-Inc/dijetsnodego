// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package djtx

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/database/memdb"
)

func TestSingletonState(t *testing.T) {
	require := require.New(t)

	db := memdb.New()
	s := NewSingletonState(db)

	isInitialized, err := s.IsInitialized()
	require.NoError(err)
	require.False(isInitialized)

	err = s.SetInitialized()
	require.NoError(err)

	isInitialized, err = s.IsInitialized()
	require.NoError(err)
	require.True(isInitialized)
}
