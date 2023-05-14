package pkg

// TODO: Add tests

import (
	"crypto/hmac"
	"encoding/base32"
	"encoding/binary"
	"hash"
	"math"
	"strings"
)

type Algorithm func() hash.Hash

type Hotp struct {
	Secret    []byte
	Algorithm Algorithm
	Digits    int
}

func NewHotp(secret []byte, algorithm Algorithm, digits int) *Hotp {
	return &Hotp{
		Secret:    secret,
		Algorithm: algorithm,
		Digits:    digits,
	}
}

func NewHotpFromBase32(secret string, algorithm Algorithm, digits int) (*Hotp, error) {
	// Normally strings used for hotp do not contain padding
	hasPadding := strings.Contains(secret, "=")
	padding := base32.NoPadding

	if hasPadding {
		padding = base32.StdPadding
	}

	// Decode the secret
	decodedSecret, err := base32.StdEncoding.WithPadding(padding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return nil, err
	}

	return NewHotp(decodedSecret, algorithm, digits), nil
}

// Calculates the Hotp code, taking a counter as moving factor.
// TODO: Maybe change the output to a string and prepend the result with 0s
// TODO: Add tests
func (h *Hotp) Calculate(movingFactor uint64) uint32 {
	digest := calculateDigest(movingFactor, h.Algorithm, h.Secret)
	fullCode := encodeDigest(digest, calculateOffset(digest))
	return shortenCodeToDigits(fullCode, h.Digits)
}

// Calculates the Hotp code, taking a counter as moving factor.
// It uses a custom offset to extract 4 bytes from the HMAC-SHA Digest.
// Keep in mind that the max value of the offset is the last index of the resulting digest minus four bytes.
// Therefore, the offset has to be between (inclusive) 0 and 15 for SHA1, 27 for SHA256 and 59 for SHA512. // TODO: Revalidate that claim lol
// TODO: Add tests
// TODO: Decide if this should be exposed
func (h *Hotp) calculateCustomOffset(movingFactor uint64, offset uint8) uint32 {
	digest := calculateDigest(movingFactor, h.Algorithm, h.Secret)
	fullCode := encodeDigest(digest, offset)
	return shortenCodeToDigits(fullCode, h.Digits)
}

// Calculate hmac digest of the moving Factor
func calculateDigest(movingFactor uint64, algorithm Algorithm, secret []byte) []byte {
	hmacInstance := hmac.New(algorithm, secret)
	binary.Write(hmacInstance, binary.BigEndian, movingFactor)
	return hmacInstance.Sum(nil)
}

// Encode the digest as a 31 bit uint32 using the provided offset
func encodeDigest(digest []byte, offset uint8) uint32 {
	codeAsBytes := digest[offset : offset+4]
	codeAsUint := binary.BigEndian.Uint32(codeAsBytes)

	// Shorten the code to 31 bit
	return codeAsUint & 0x7fffffff
}

// Calculate the offset from last byte
func calculateOffset(digest []byte) uint8 {
	lastByte := digest[len(digest)-1]
	return lastByte & 0xF
}

// Shorten the code to the desired length.
// The length od the calculated code can, by design, not be higher than 10 characters.
// If the digit size is larger then 10, it just gets prefixed with zeros.
func shortenCodeToDigits(fullCode uint32, digits int) uint32 {
	if digits < 10 {
		modulusBase := uint32(math.Floor(math.Pow10(digits)))
		return fullCode % modulusBase
	}

	return fullCode
}
