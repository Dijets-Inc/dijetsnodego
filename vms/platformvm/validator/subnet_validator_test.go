// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lasthyphen/dijetsnodego/ids"
	"github.com/lasthyphen/dijetsnodego/utils/constants"
)

func TestSubnetValidatorVerifySubnetID(t *testing.T) {
	require := require.New(t)

	// Error path
	{
		vdr := &SubnetValidator{
			Subnet: constants.PrimaryNetworkID,
		}

		require.Equal(errBadSubnetID, vdr.Verify())
	}

	// Happy path
	{
		vdr := &SubnetValidator{
			Subnet: ids.GenerateTestID(),
			Validator: Validator{
				Wght: 1,
			},
		}

		require.Equal(nil, vdr.Verify())
	}
}
