// TODO: Add Secret validation and otp verification
// TODO: Add otp verification
// TODO: Add remaining time calculation
// TODO: Add url serialization and deserialization
// TODO: Add account and issuer
package hotp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"math"
	"strings"

	"bode.fun/otp"
)

type HotpOptions struct {
	// TODO: Change this to a serialisable format
	Algorithm otp.Algorithm
	Digits    uint
}

type HotpOption func(*HotpOptions)

func WithDigits(digits uint) HotpOption {
	return func(ho *HotpOptions) {
		ho.Digits = digits
	}
}

func WithAlgorithm(algorithm otp.Algorithm) HotpOption {
	return func(ho *HotpOptions) {
		ho.Algorithm = algorithm
	}
}

const defaultDigits uint = 6

var defaultAlgorithm otp.Algorithm = sha1.New

// Hotp is a stateful counter based One Time Password algorithm.
//
// It uses a counter to verify the user. This counter has to be stored on
// the server and the client.
type Hotp struct {
	secret    []byte
	algorithm otp.Algorithm
	digits    uint
}

// Create a Hotp instance from a unencoded secret.
// The algorithm, that is usually used, is sha1.
//
// Example:
//
//	hotp := New([]byte("12345678901234567890"),
//				WithAlgorithm(sha1.New),
//				WithDigits(6),
//			)
func New(secret []byte, options ...HotpOption) *Hotp {
	opts := &HotpOptions{
		Algorithm: defaultAlgorithm,
		Digits:    defaultDigits,
	}

	for _, option := range options {
		option(opts)
	}

	return &Hotp{
		secret:    secret,
		algorithm: opts.Algorithm,
		digits:    opts.Digits,
	}
}

// Create a Hotp instance from a base32 encoded secret.
// The algorithm, that is usually used, is sha1.
//
// Example:
//
//	hotp := NewFromBase32("GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ",
//				WithAlgorithm(sha1.New),
//				WithDigits(6),
//			)
func NewFromBase32(secret string, options ...HotpOption) (*Hotp, error) {
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

	return New(decodedSecret, options...), nil
}

func (h *Hotp) Digits() uint {
	return h.digits
}

func (h *Hotp) Secret() []byte {
	return h.secret
}

func (h *Hotp) Algorithm() otp.Algorithm {
	return h.algorithm
}

// Calculates the Hotp code, taking a counter as moving factor.
//
// TODO: Maybe change the output to a string and prepend the result with 0s
func (h *Hotp) Calculate(movingFactor uint64) uint32 {
	digest := calculateDigest(movingFactor, h.algorithm, h.secret)
	offset := calculateOffset(digest)
	fullCode := encodeDigest(digest, offset)
	return shortenCodeToDigits(fullCode, h.digits)
}

// Calculates the Hotp code, taking a counter as moving factor.
// It uses a custom offset to extract 4 bytes from the HMAC-SHA Digest.
// Keep in mind that the max value of the offset is the last index of the
// resulting digest minus four bytes.
// Therefore, the offset has to be between (inclusive) 0 and 16 for SHA1 (20 byte digest),
// 28 for SHA256 (32 byte digest) and 60 for SHA512 (64 byte digest).
//
// TODO: Decide if this should be exposed
func (h *Hotp) calculateCustomOffset(movingFactor uint64, offset uint8) uint32 {
	digest := calculateDigest(movingFactor, h.algorithm, h.secret)
	fullCode := encodeDigest(digest, offset)
	return shortenCodeToDigits(fullCode, h.digits)
}

// Calculate hmac digest of the moving Factor
func calculateDigest(movingFactor uint64, algorithm otp.Algorithm, secret []byte) []byte {
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
