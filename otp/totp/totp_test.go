package totp_test

import (
	"testing"

	"bode.fun/otp/totp"
	"github.com/matryer/is"
)

// This test validates the implementation against the RFC 6238
// "TOTP: Time-Based One-Time Password Algorithm"
// The values are available at Appendix B
// https://www.rfc-editor.org/rfc/rfc6238#appendix-B
// It just checks sha1
func Test_Rfc6238_Sha1(t *testing.T) {
	is := is.New(t)
	totp := totp.New(
		[]byte("12345678901234567890"),
		totp.WithDigits(8),
	)

	{
		code := totp.Calculate(59)
		is.Equal(uint32(94287082), code)
	}

	{
		code := totp.Calculate(1111111109)
		is.Equal(uint32(7081804), code) // 07081804
	}

	{
		code := totp.Calculate(1111111111)
		is.Equal(uint32(14050471), code)
	}

	{
		code := totp.Calculate(1234567890)
		is.Equal(uint32(89005924), code)
	}

	{
		code := totp.Calculate(2000000000)
		is.Equal(uint32(69279037), code)
	}

	{
		code := totp.Calculate(20000000000)
		is.Equal(uint32(65353130), code)
	}
}

// This test validates the implementation against the RFC 6238
// "TOTP: Time-Based One-Time Password Algorithm"
// The values are available at Appendix B
// https://www.rfc-editor.org/rfc/rfc6238#appendix-B
// It just checks sha256
func Test_Rfc6238_Sha256(t *testing.T) {
	is := is.New(t)
	totp := totp.New(
		[]byte("12345678901234567890123456789012"),
		totp.WithAlgorithm(totp.Sha256),
		totp.WithDigits(8),
	)

	{
		code := totp.Calculate(59)
		is.Equal(uint32(46119246), code)
	}

	{
		code := totp.Calculate(1111111109)
		is.Equal(uint32(68084774), code)
	}

	{
		code := totp.Calculate(1111111111)
		is.Equal(uint32(67062674), code)
	}

	{
		code := totp.Calculate(1234567890)
		is.Equal(uint32(91819424), code)
	}

	{
		code := totp.Calculate(2000000000)
		is.Equal(uint32(90698825), code)
	}

	{
		code := totp.Calculate(20000000000)
		is.Equal(uint32(77737706), code)
	}
}

// This test validates the implementation against the RFC 6238
// "TOTP: Time-Based One-Time Password Algorithm"
// The values are available at Appendix B
// https://www.rfc-editor.org/rfc/rfc6238#appendix-B
// It just checks sha512
func Test_Rfc6238_Sha512(t *testing.T) {
	is := is.New(t)
	totp := totp.New(
		[]byte("1234567890123456789012345678901234567890123456789012345678901234"),
		totp.WithAlgorithm(totp.Sha512),
		totp.WithDigits(8),
	)

	{
		code := totp.Calculate(59)
		is.Equal(uint32(90693936), code)
	}

	{
		code := totp.Calculate(1111111109)
		is.Equal(uint32(25091201), code)
	}

	{
		code := totp.Calculate(1111111111)
		is.Equal(uint32(99943326), code)
	}

	{
		code := totp.Calculate(1234567890)
		is.Equal(uint32(93441116), code)
	}

	{
		code := totp.Calculate(2000000000)
		is.Equal(uint32(38618901), code)
	}

	{
		code := totp.Calculate(20000000000)
		is.Equal(uint32(47863826), code)
	}
}
