// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package djtx

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/database"
	"github.com/lasthyphen/dijetsnodego/database/memdb"
	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/snow/choices"
)

func TestStatusState(t *testing.T) {
	require := require.New(t)
	id0 := ids.GenerateTestID()

	db := memdb.New()
	s := NewStatusState(db)

	_, err := s.GetStatus(id0)
	require.Equal(database.ErrNotFound, err)

	_, err = s.GetStatus(id0)
	require.Equal(database.ErrNotFound, err)

	err = s.PutStatus(id0, choices.Accepted)
	require.NoError(err)

	status, err := s.GetStatus(id0)
	require.NoError(err)
	require.Equal(choices.Accepted, status)

	err = s.DeleteStatus(id0)
	require.NoError(err)

	_, err = s.GetStatus(id0)
	require.Equal(database.ErrNotFound, err)

	err = s.PutStatus(id0, choices.Accepted)
	require.NoError(err)

	s = NewStatusState(db)

	status, err = s.GetStatus(id0)
	require.NoError(err)
	require.Equal(choices.Accepted, status)
}
