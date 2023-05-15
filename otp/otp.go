package otp

import "hash"

type Algorithm func() hash.Hash
