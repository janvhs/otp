// TODO: Add Secret validation and otp verification
// TODO: Add secret generation
// TODO: Add otp verification
// TODO: Add remaining time calculation
// TODO: Add window for validation
package totp

import (
	"encoding/base32"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bode.fun/otp/hotp"
)

type Algorithm = hotp.Algorithm

const (
	Sha1   Algorithm = hotp.Sha1
	Sha256 Algorithm = hotp.Sha256
	Sha512 Algorithm = hotp.Sha512
)

type totpOptions struct {
	algorithm Algorithm
	digits    uint
	stepSize  uint
	account   string
	issuer    string
}

type TotpOption func(*totpOptions)

func WithStepSize(stepSize uint) TotpOption {
	return func(to *totpOptions) {
		to.stepSize = stepSize
	}
}

// Alias for totp.WithStepSize
func WithPeriod(period uint) TotpOption {
	return WithStepSize(period)
}

func WithDigits(digits uint) TotpOption {
	return func(to *totpOptions) {
		to.digits = digits
	}
}

func WithAlgorithm(algorithm Algorithm) TotpOption {
	return func(to *totpOptions) {
		to.algorithm = algorithm
	}
}

func WithAccount(account string) TotpOption {
	return func(to *totpOptions) {
		to.account = account
	}
}

func WithIssuer(issuer string) TotpOption {
	return func(to *totpOptions) {
		to.issuer = issuer
	}
}

const defaultDigits uint = 6
const defaultStepSize uint = 30

var defaultAlgorithm Algorithm = Sha1

// Totp is a stateless time based One Time Password algorithm.
//
// It uses a the unix time to verify the user.
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
	opts := &totpOptions{
		algorithm: defaultAlgorithm,
		digits:    defaultDigits,
		stepSize:  defaultStepSize,
	}

	for _, option := range options {
		option(opts)
	}

	hotp, err := hotp.NewFromBase32(
		secret,
		hotp.WithDigits(opts.digits),
		hotp.WithAccount(opts.account),
		hotp.WithIssuer(opts.issuer),
		hotp.WithAlgorithm(opts.algorithm),
	)

	if err != nil {
		return nil, err
	}
	return &Totp{
		hotp:     hotp,
		stepSize: opts.stepSize,
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
	opts := &totpOptions{
		algorithm: defaultAlgorithm,
		digits:    defaultDigits,
		stepSize:  defaultStepSize,
	}

	for _, option := range options {
		option(opts)
	}

	hotp := hotp.New(secret,
		hotp.WithDigits(opts.digits),
		hotp.WithAccount(opts.account),
		hotp.WithIssuer(opts.issuer),
		hotp.WithAlgorithm(opts.algorithm),
	)

	return &Totp{
		hotp:     hotp,
		stepSize: opts.stepSize,
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

func (t *Totp) StepSize() uint {
	return t.stepSize
}

func (t *Totp) Period() uint {
	return t.StepSize()
}

func (t *Totp) Account() string {
	return t.hotp.Account()
}

func (t *Totp) Issuer() string {
	return t.hotp.Issuer()
}

func (t *Totp) Label() string {
	return t.hotp.Label()
}

// TODO: Maybe change the output to a string and prepend the result with 0s
func (t *Totp) Calculate(movingFactor uint64) uint32 {
	flooredSeconds := float64(movingFactor)
	movingFactor = uint64(math.Floor(flooredSeconds / float64(t.stepSize)))
	return t.hotp.Calculate(movingFactor)
}

// TODO: Maybe change the output to a string and prepend the result with 0s
func (t *Totp) Now() uint32 {
	unixSeconds := time.Now().Unix()
	return t.Calculate(uint64(unixSeconds))
}

// Alias for totp.Now()
func (t *Totp) CalculateNow() uint32 {
	return t.Now()
}

// TODO: Add tests
// References: https://docs.yubico.com/yesdk/users-manual/application-oath/uri-string-format.html
// TODO: Look up other references (saved a bunch in otp on iPhone)
func (t *Totp) ToUrl() string {
	label := t.Account()

	if t.Issuer() != "" {
		label = label + ":" + t.Issuer()
	}

	otpUrl := &url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   label,
	}

	encodedSecret := base32.StdEncoding.EncodeToString(t.Secret())

	query := otpUrl.Query()

	query.Set("secret", encodedSecret)
	query.Set("period", fmt.Sprint(t.StepSize()))

	query.Set("algorithm", string(t.Algorithm()))

	query.Set("digits", fmt.Sprint(t.Digits()))

	if t.Issuer() != "" {
		query.Set("issuer", t.Issuer())
	}

	otpUrl.RawQuery = query.Encode()

	return otpUrl.String()
}

func NewFromUrl(rawUrl string) (*Totp, error) {
	otpUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	encodedSecret := otpUrl.Query().Get("secret")

	totpOptions := []TotpOption{}

	label := otpUrl.Path
	label = strings.TrimPrefix(label, "/")
	account, _, _ := strings.Cut(label, ":")
	if account != "" {
		totpOptions = append(totpOptions, WithAccount(account))
	}

	periodAsString := otpUrl.Query().Get("period")
	if periodAsString != "" {
		periodAsInt, err := strconv.Atoi(periodAsString)
		if err != nil {
			return nil, err
		}

		period := uint(periodAsInt)
		totpOptions = append(totpOptions, WithStepSize(period))
	}

	digitsAsString := otpUrl.Query().Get("digits")
	if digitsAsString != "" {
		digitsAsInt, err := strconv.Atoi(digitsAsString)
		if err != nil {
			return nil, err
		}

		digits := uint(digitsAsInt)
		totpOptions = append(totpOptions, WithDigits(digits))
	}

	issuer := otpUrl.Query().Get("issuer")
	if issuer != "" {
		totpOptions = append(totpOptions, WithIssuer(issuer))
	}

	algorithm := otpUrl.Query().Get("algorithm")
	if algorithm == "" {
		algorithm = string(Sha1)
	}
	totpOptions = append(totpOptions, WithAlgorithm(hotp.Algorithm(algorithm)))

	return NewFromBase32(encodedSecret, totpOptions...)
}
