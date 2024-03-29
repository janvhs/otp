// TODO: Add Secret validation
// TODO: Add otp verification
// TODO: Add verification against a window
package hotp

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Algorithm string

const (
	Sha1   Algorithm = "sha1"
	Sha256 Algorithm = "sha256"
	Sha512 Algorithm = "sha512"
)

func (a Algorithm) ToHashFunction() func() hash.Hash {
	switch a {
	case Sha1:
		return sha1.New
	case Sha256:
		return sha256.New
	case Sha512:
		return sha512.New
	default:
		return nil
	}
}

type hotpOptions struct {
	algorithm Algorithm
	digits    uint
	account   string
	issuer    string
}

type HotpOption func(*hotpOptions)

func WithDigits(digits uint) HotpOption {
	return func(ho *hotpOptions) {
		ho.digits = digits
	}
}

func WithAlgorithm(algorithm Algorithm) HotpOption {
	return func(ho *hotpOptions) {
		ho.algorithm = algorithm
	}
}

func WithAccount(account string) HotpOption {
	return func(ho *hotpOptions) {
		ho.account = account
	}
}

func WithIssuer(issuer string) HotpOption {
	return func(ho *hotpOptions) {
		ho.issuer = issuer
	}
}

const defaultDigits uint = 6

var defaultAlgorithm Algorithm = Sha1

// Hotp is a stateful counter based One Time Password algorithm.
//
// It uses a counter to verify the user. This counter has to be stored on
// the server and the client.
type Hotp struct {
	secret    []byte
	algorithm Algorithm
	digits    uint
	account   string
	issuer    string
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
	opts := &hotpOptions{
		algorithm: defaultAlgorithm,
		digits:    defaultDigits,
	}

	for _, option := range options {
		option(opts)
	}

	return &Hotp{
		secret:    secret,
		algorithm: opts.algorithm,
		digits:    opts.digits,
		account:   opts.account,
		issuer:    opts.issuer,
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

func (h *Hotp) Algorithm() Algorithm {
	return h.algorithm
}

func (h *Hotp) Account() string {
	return h.account
}

func (h *Hotp) Issuer() string {
	return h.issuer
}

func (h *Hotp) Label() string {
	label := h.Account()

	if h.Issuer() != "" {
		label = label + ":" + h.Issuer()
	}

	return url.PathEscape(label)
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

// TODO: Add tests
// References: https://docs.yubico.com/yesdk/users-manual/application-oath/uri-string-format.html
func (h *Hotp) ToUrl(counter uint64) string {
	label := h.Account()

	if h.Issuer() != "" {
		label = label + ":" + h.Issuer()
	}

	otpUrl := &url.URL{
		Scheme: "otpauth",
		Host:   "hotp",
		Path:   label,
	}

	encodedSecret := base32.StdEncoding.EncodeToString(h.Secret())

	query := otpUrl.Query()

	query.Set("secret", encodedSecret)
	query.Set("counter", fmt.Sprint(counter))

	query.Set("algorithm", string(h.Algorithm()))

	query.Set("digits", fmt.Sprint(h.Digits()))

	if h.Issuer() != "" {
		query.Set("issuer", h.Issuer())
	}

	otpUrl.RawQuery = query.Encode()

	return otpUrl.String()
}

// TODO: Add some tests
func NewFromUrl(rawUrl string) (*Hotp, uint64, error) {
	var counter uint64

	otpUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, counter, err
	}

	encodedSecret := otpUrl.Query().Get("secret")

	hotpOptions := []HotpOption{}

	label := otpUrl.Path
	label = strings.TrimPrefix(label, "/")
	account, _, _ := strings.Cut(label, ":")
	if account != "" {
		hotpOptions = append(hotpOptions, WithAccount(account))
	}

	counterAsString := otpUrl.Query().Get("counter")
	if counterAsString != "" {
		counterAsInt, err := strconv.Atoi(counterAsString)
		if err != nil {
			return nil, counter, err
		}

		counter = uint64(counterAsInt)
	}

	digitsAsString := otpUrl.Query().Get("digits")
	if digitsAsString != "" {
		digitsAsInt, err := strconv.Atoi(digitsAsString)
		if err != nil {
			return nil, counter, err
		}

		digits := uint(digitsAsInt)
		hotpOptions = append(hotpOptions, WithDigits(digits))
	}

	issuer := otpUrl.Query().Get("issuer")
	if issuer != "" {
		hotpOptions = append(hotpOptions, WithIssuer(issuer))
	}

	algorithm := otpUrl.Query().Get("algorithm")
	if algorithm == "" {
		algorithm = string(Sha1)
	}
	hotpOptions = append(hotpOptions, WithAlgorithm(Algorithm(algorithm)))

	hotpInstance, err := NewFromBase32(encodedSecret, hotpOptions...)
	return hotpInstance, counter, err
}

// Calculate hmac digest of the moving Factor
func calculateDigest(movingFactor uint64, algorithm Algorithm, secret []byte) []byte {
	hmacInstance := hmac.New(algorithm.ToHashFunction(), secret)
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
