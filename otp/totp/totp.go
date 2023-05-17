// TODO: Add Secret validation and otp verification
package totp

import (
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bode.fun/otp"
	"bode.fun/otp/hotp"
)

type TotpOptions struct {
	// TODO: Change this to a serialisable format
	Algorithm otp.Algorithm
	Digits    uint
	StepSize  uint
	Account   string
	Issuer    string
}

type TotpOption func(*TotpOptions)

func WithStepSize(stepSize uint) TotpOption {
	return func(to *TotpOptions) {
		to.StepSize = stepSize
	}
}

// A alias for WithStepSize
// TODO: In the url the step size is named period. Is an alias helpful?
func WithPeriod(period uint) TotpOption {
	return WithStepSize(period)
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

func WithAccount(account string) TotpOption {
	return func(to *TotpOptions) {
		to.Account = account
	}
}

func WithIssuer(issuer string) TotpOption {
	return func(to *TotpOptions) {
		to.Issuer = issuer
	}
}

const defaultDigits uint = 6
const defaultStepSize uint = 30

var defaultAlgorithm otp.Algorithm = sha1.New

type Totp struct {
	hotp     *hotp.Hotp
	stepSize uint
	// TODO: Move this to hotp
	account string
	issuer  string
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
		// TODO: Move this to hotp
		account: opts.Account,
		issuer:  opts.Issuer,
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

func (t *Totp) StepSize() uint {
	return t.stepSize
}

func (t *Totp) Period() uint {
	return t.StepSize()
}

// TODO: Move account to hotp
func (t *Totp) Account() string {
	return t.account
}

// TODO: Move issuer to hotp
func (t *Totp) Issuer() string {
	return t.issuer
}

// TODO: Move label to hotp
func (t *Totp) Label() string {
	label := t.Account()

	if t.Issuer() != "" {
		label = label + ":" + t.Issuer()
	}

	return url.PathEscape(label)
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

	// TODO: Add algorithm, currently sha1 is always assumed

	query.Set("digits", fmt.Sprint(t.Digits()))

	if t.issuer != "" {
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
	// TODO: Should this throw when the constructor does not?
	if encodedSecret == "" {
		return nil, fmt.Errorf("the provided secret can not be empty")
	}

	totpOptions := []TotpOption{}

	label := otpUrl.Path
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
		// TODO: Uri decode?
		totpOptions = append(totpOptions, WithIssuer(issuer))
	}

	return NewFromBase32(encodedSecret, totpOptions...)
}
