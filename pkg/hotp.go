package pkg

import (
	"crypto/hmac"
	"encoding/base32"
	"encoding/binary"
	"hash"
	"math"
	"strings"
)

type Algorithm func() hash.Hash

// Hotp is a counter based One Time Password algorithm.
//
// It uses a counter to verify the user. This counter has to be stored on
// the server and the client.
type Hotp struct {
	Secret    []byte
	Algorithm Algorithm
	Digits    uint
}

// Create a Hotp instance from a unencoded secret.
// The algorithm, that is usually used, is sha1.
//
// Example:
//
//	hotp := NewHotp([]byte("12345678901234567890"), sha1.New, 6)
func NewHotp(secret []byte, algorithm Algorithm, digits uint) *Hotp {
	return &Hotp{
		Secret:    secret,
		Algorithm: algorithm,
		Digits:    digits,
	}
}

// Create a Hotp instance from a base32 encoded secret.
// The algorithm, that is usually used, is sha1.
//
// Example:
//
//	hotp := NewHotpFromBase32("GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ", sha1.New, 6)
func NewHotpFromBase32(secret string, algorithm Algorithm, digits uint) (*Hotp, error) {
	// Usually strings, used for hotp, do not contain padding
	hasPadding := strings.Contains(secret, "=")
	padding := base32.NoPadding

	if hasPadding {
		padding = base32.StdPadding
	}

	// Decode the secret
	decodedSecret, err := base32.StdEncoding.
		WithPadding(padding).
		DecodeString(strings.ToUpper(secret))

	if err != nil {
		return nil, err
	}

	return NewHotp(decodedSecret, algorithm, digits), nil
}

// Calculates the Hotp code, taking a counter as moving factor.
// TODO: Maybe change the output to a string and prepend the result with 0s
func (h *Hotp) Calculate(movingFactor uint64) uint32 {
	digest := calculateDigest(movingFactor, h.Algorithm, h.Secret)
	offset := calculateOffset(digest)
	fullCode := encodeDigest(digest, offset)
	return shortenCodeToDigits(fullCode, h.Digits)
}

// Calculates the Hotp code, taking a counter as moving factor.
// It uses a custom offset to extract 4 bytes from the HMAC-SHA Digest.
// Keep in mind that the max value of the offset is the last index of the
// resulting digest minus four bytes.
// Therefore, the offset has to be between (inclusive) 0 and 16 for SHA1 (20 byte digest),
// 28 for SHA256 (32 byte digest) and 60 for SHA512 (64 byte digest).
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
func shortenCodeToDigits(fullCode uint32, digits uint) uint32 {
	if digits < 10 {
		modulusBase := uint32(math.Pow10(int(digits)))
		return fullCode % modulusBase
	}

	return fullCode
}
