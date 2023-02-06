// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"github.com/lasthyphen/dijetsnodego/utils/constants"
	"github.com/lasthyphen/dijetsnodego/utils/sampler"
)

// getIPs returns the beacon IPs for each network
func getIPs(networkID uint32) []string {
	switch networkID {
	case constants.MainnetID:
		return []string{
			"51.142.161.145:9651",
			"51.142.161.145:9653",
			"51.142.161.145:9655",
			"51.142.161.145:9657",
			"51.142.161.145:9659",
		}
	case constants.TahoeID:
		return []string{
			"51.142.161.145:9651",
			"51.142.161.145:9653",
			"51.142.161.145:9655",
			"51.142.161.145:9657",
			"51.142.161.145:9659",
		}
	default:
		return nil
	}
}

// getNodeIDs returns the beacon node IDs for each network
func getNodeIDs(networkID uint32) []string {
	switch networkID {
	case constants.MainnetID:
		return []string{
			"NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg",
			"NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ",
			"NodeID-NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN",
			"NodeID-GWPcbFJZFfZreETSoWjPimr846mXEKCtu",
			"NodeID-P7oB2McjBGgW2NXXWVYjV8JEDFoW9xDE5",
		}
	case constants.TahoeID:
		return []string{
			"NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg",
			"NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ",
			"NodeID-NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN",
			"NodeID-GWPcbFJZFfZreETSoWjPimr846mXEKCtu",
			"NodeID-P7oB2McjBGgW2NXXWVYjV8JEDFoW9xDE5",
		}
	default:
		return nil
	}
}

// SampleBeacons returns the some beacons this node should connect to
func SampleBeacons(networkID uint32, count int) ([]string, []string) {
	ips := getIPs(networkID)
	ids := getNodeIDs(networkID)

	if numIPs := len(ips); numIPs < count {
		count = numIPs
	}

	sampledIPs := make([]string, 0, count)
	sampledIDs := make([]string, 0, count)

	s := sampler.NewUniform()
	_ = s.Initialize(uint64(len(ips)))
	indices, _ := s.Sample(count)
	for _, index := range indices {
		sampledIPs = append(sampledIPs, ips[int(index)])
		sampledIDs = append(sampledIDs, ids[int(index)])
	}

	return sampledIPs, sampledIDs
}
