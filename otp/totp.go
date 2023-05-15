package otp

import (
	"math"
	"time"
)

const defaultStepSize uint = 30

type Totp struct {
	hotp     *Hotp
	StepSize uint
}

// Create a Totp instance from a base32 encoded secret.
// The algorithm, that is usually used, is sha1.
//
// Example:
//
//	totp := NewTotpFromBase32("GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ", sha1.New, 6)
func NewTotpFromBase32(secret string, algorithm Algorithm, digits uint) (*Totp, error) {
	hotp, err := NewHotpFromBase32(secret, algorithm, digits)
	if err != nil {
		return nil, err
	}
	return &Totp{
		hotp:     hotp,
		StepSize: defaultStepSize,
	}, nil
}

// Create a Totp instance from a unencoded secret.
// The algorithm, that is usually used, is sha1.
//
// Example:
//
//	totp := NewTotp([]byte("12345678901234567890"), sha1.New, 6)
func NewTotp(secret []byte, algorithm Algorithm, digits uint) *Totp {
	hotp := NewHotp(secret, algorithm, digits)
	return &Totp{
		hotp:     hotp,
		StepSize: defaultStepSize,
	}
}

func (t *Totp) Digits() uint {
	return t.hotp.Digits()
}

func (t *Totp) Secret() []byte {
	return t.hotp.Secret()
}

func (t *Totp) Algorithm() Algorithm {
	return t.hotp.Algorithm()
}

func (t *Totp) Calculate(movingFactor uint64) uint32 {
	flooredSeconds := float64(movingFactor)
	movingFactor = uint64(math.Floor(flooredSeconds / float64(t.StepSize)))
	return t.hotp.Calculate(movingFactor)
}

func (t *Totp) Now() uint32 {
	unixSeconds := time.Now().Unix()
	return t.Calculate(uint64(unixSeconds))
}

func (t *Totp) CalculateNow() uint32 {
	return t.Now()
}
