// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package codec

// GeneralCodec marshals and unmarshals structs including interfaces
type GeneralCodec interface {
	Codec
	Registry
}
