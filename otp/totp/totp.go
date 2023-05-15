package totp

import (
	"crypto/sha1"
	"math"
	"time"

	"bode.fun/otp"
	"bode.fun/otp/hotp"
)

type TotpOptions struct {
	Algorithm otp.Algorithm
	Digits    uint
	StepSize  uint
}

type TotpOption func(*TotpOptions)

func WithStepSize(stepSize uint) TotpOption {
	return func(to *TotpOptions) {
		to.StepSize = stepSize
	}
}

func WithDigits(digits uint) TotpOption {
	return func(to *TotpOptions) {
		to.Digits = digits
	}
}

func WithAlgorithm(algorithm otp.Algorithm) TotpOption {
	return func(to *TotpOptions) {
		to.Algorithm = algorithm
	}
}

const defaultDigits uint = 6
const defaultStepSize uint = 30

var defaultAlgorithm otp.Algorithm = sha1.New

type Totp struct {
	hotp     *hotp.Hotp
	stepSize uint
}

// Create a Totp instance from a base32 encoded secret.
// The algorithm, that is usually used, is sha1.
//
// Example:
//
//	totp := NewFromBase32("GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ",
//				WithAlgorithm(sha1.New),
//				WithDigits(6),
//			)
func NewFromBase32(secret string, options ...TotpOption) (*Totp, error) {
	opts := &TotpOptions{
		Algorithm: defaultAlgorithm,
		Digits:    defaultDigits,
		StepSize:  defaultStepSize,
	}

	for _, option := range options {
		option(opts)
	}

	hotp, err := hotp.NewFromBase32(
		secret,
		hotp.WithDigits(opts.Digits),
		hotp.WithAlgorithm(opts.Algorithm),
	)

	if err != nil {
		return nil, err
	}
	return &Totp{
		hotp:     hotp,
		stepSize: opts.StepSize,
	}, nil
}

// Create a Totp instance from a unencoded secret.
// The algorithm, that is usually used, is sha1.
//
// Example:
//
//	totp := New([]byte("12345678901234567890"),
//				WithAlgorithm(sha1.New),
//				WithDigits(6),
//			)
func New(secret []byte, options ...TotpOption) *Totp {
	opts := &TotpOptions{
		Algorithm: defaultAlgorithm,
		Digits:    defaultDigits,
		StepSize:  defaultStepSize,
	}

	for _, option := range options {
		option(opts)
	}

	hotp := hotp.New(secret,
		hotp.WithDigits(opts.Digits),
		hotp.WithAlgorithm(opts.Algorithm),
	)

	return &Totp{
		hotp:     hotp,
		stepSize: opts.StepSize,
	}
}

func (t *Totp) Digits() uint {
	return t.hotp.Digits()
}

func (t *Totp) Secret() []byte {
	return t.hotp.Secret()
}

func (t *Totp) Algorithm() otp.Algorithm {
	return t.hotp.Algorithm()
}

func (t *Totp) Calculate(movingFactor uint64) uint32 {
	flooredSeconds := float64(movingFactor)
	movingFactor = uint64(math.Floor(flooredSeconds / float64(t.stepSize)))
	return t.hotp.Calculate(movingFactor)
}

func (t *Totp) Now() uint32 {
	unixSeconds := time.Now().Unix()
	return t.Calculate(uint64(unixSeconds))
}

func (t *Totp) CalculateNow() uint32 {
	return t.Now()
}
