package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/retinotopic/pokerGO/pkg/randfuncs"
)

var Secretkey = []byte("cbaTd3Dx9_dfknwPsc5T0rQMx34SvJJf5xvxf7nab")

func WriteHashCookie(w http.ResponseWriter, key []byte) *http.Cookie {
	mac := hmac.New(sha256.New, key)
	r0 := rand.New(rand.NewSource(time.Now().Unix()))
	time.Sleep(time.Millisecond * 25)
	r1 := rand.New(rand.NewSource(time.Now().Unix()))
	r1.Seed(time.Now().UnixNano())
	cookie := &http.Cookie{Name: randfuncs.RandomString(15, r0), Value: randfuncs.RandomString(20, r1), Secure: true, Path: "/"}
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	cookie.Value = cookie.Value + signature
	http.SetCookie(w, cookie)
	return cookie
}
func ReadHashCookie(r *http.Request, key []byte, cookies []*http.Cookie) (*http.Cookie, error) {
	if len(cookies) == 0 {
		return nil, errors.New("zero cookies")
	}
	c := cookies[0]
	name := c.Name
	valueHash := c.Value
	signature := valueHash[20:]
	value := valueHash[:20]

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, errors.New("ValidationErr")
	}
	return c, nil
}
