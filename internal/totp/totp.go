package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	defaultStep   = 30
	defaultDigits = 6
)

func GenerateSecret() (string, error) {

	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}

	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	return enc.EncodeToString(secret), nil
}

func OTPAuthURL(issuer string, account string, secret string) string {

	label := fmt.Sprintf("%s:%s", issuer, account)

	v := url.Values{}
	v.Set("secret", secret)
	v.Set("issuer", issuer)
	v.Set("algorithm", "SHA1")
	v.Set("digits", fmt.Sprintf("%d", defaultDigits))
	v.Set("period", fmt.Sprintf("%d", defaultStep))

	return "otpauth://totp/" + url.PathEscape(label) + "?" + v.Encode()
}

func Validate(secret, code string, now time.Time) bool {

	code = strings.ReplaceAll(strings.TrimSpace(code), " ", "")
	if len(code) != defaultDigits {
		return false
	}

	counter := now.Unix() / defaultStep
	for offset := int64(-1); offset <= 1; offset++ {
		if code == generateCode(secret, counter+offset) {
			return true
		}
	}

	return false
}

func generateCode(secret string, counter int64) string {

	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	key, err := enc.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return ""
	}

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(counter))

	hash := hmac.New(sha1.New, key)
	_, _ = hash.Write(buf[:])
	sum := hash.Sum(nil)

	offset := sum[len(sum)-1] & 0x0f
	bin := (int(sum[offset])&0x7f)<<24 |
		(int(sum[offset+1])&0xff)<<16 |
		(int(sum[offset+2])&0xff)<<8 |
		(int(sum[offset+3]) & 0xff)

	otp := bin % pow10(defaultDigits)
	return fmt.Sprintf("%0*d", defaultDigits, otp)
}

func pow10(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}
