package otp_test

import (
	"crypto/sha1"
	"testing"

	"bode.fun/otp"
	"github.com/matryer/is"
)

// This test validates the implementation against the RFC 4226
// "HOTP: An HMAC-Based One-Time Password Algorithm"
// The values are available at Appendix D
// https://www.rfc-editor.org/rfc/rfc4226#page-32
func Test_Rfc4226(t *testing.T) {
	is := is.New(t)
	h := otp.NewHotp([]byte("12345678901234567890"), sha1.New, 6)

	{
		code := h.Calculate(0)
		is.Equal(uint32(755224), code)
	}

	{
		code := h.Calculate(1)
		is.Equal(uint32(287082), code)
	}

	{
		code := h.Calculate(2)
		is.Equal(uint32(359152), code)
	}

	{
		code := h.Calculate(3)
		is.Equal(uint32(969429), code)
	}

	{
		code := h.Calculate(4)
		is.Equal(uint32(338314), code)
	}

	{
		code := h.Calculate(5)
		is.Equal(uint32(254676), code)
	}

	{
		code := h.Calculate(6)
		is.Equal(uint32(287922), code)
	}

	{
		code := h.Calculate(7)
		is.Equal(uint32(162583), code)
	}

	{
		code := h.Calculate(8)
		is.Equal(uint32(399871), code)
	}

	{
		code := h.Calculate(9)
		is.Equal(uint32(520489), code)
	}
}
