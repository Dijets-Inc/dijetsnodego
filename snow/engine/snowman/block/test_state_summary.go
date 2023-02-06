// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"context"
	"errors"
	"testing"

	"github.com/lasthyphen/dijetsnodego/ids"
)

var (
	_ StateSummary = (*TestStateSummary)(nil)

	errAccept = errors.New("unexpectedly called Accept")
)

type TestStateSummary struct {
	IDV     ids.ID
	HeightV uint64
	BytesV  []byte

	T          *testing.T
	CantAccept bool
	AcceptF    func(context.Context) (bool, error)
}

func (s *TestStateSummary) ID() ids.ID {
	return s.IDV
}

func (s *TestStateSummary) Height() uint64 {
	return s.HeightV
}

func (s *TestStateSummary) Bytes() []byte {
	return s.BytesV
}

func (s *TestStateSummary) Accept(ctx context.Context) (bool, error) {
	if s.AcceptF != nil {
		return s.AcceptF(ctx)
	}
	if s.CantAccept && s.T != nil {
		s.T.Fatal(errAccept)
	}
	return false, errAccept
}
