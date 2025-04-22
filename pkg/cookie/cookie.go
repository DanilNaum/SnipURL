package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

var (
	ErrNoCookie = errors.New("no cookie")
)

type cookieManager struct {
	secret  []byte
	options *options
}
type options struct {
	name     string
	path     string
	secure   bool
	httpOnly bool
	sameSite http.SameSite
}
type opt func(o *options)

func WithName(name string) opt {
	return func(o *options) {
		o.name = name
	}
}
func WithPath(path string) opt {
	return func(o *options) {
		o.path = path
	}
}
func WithSecure(secure bool) opt {
	return func(o *options) {
		o.secure = secure
	}
}
func WithHTTPOnly(httpOnly bool) opt {
	return func(o *options) {
		o.httpOnly = httpOnly
	}
}
func WithSameSite(sameSite http.SameSite) opt {
	return func(o *options) {
		o.sameSite = sameSite
	}
}
func NewCookieManager(secret []byte, opts ...opt) *cookieManager {
	opt := &options{
		name:     "user",
		path:     "/",
		secure:   true,
		httpOnly: true,
		sameSite: http.SameSiteLaxMode,
	}
	for _, o := range opts {
		o(opt)
	}
	return &cookieManager{
		secret:  secret,
		options: opt,
	}
}
func (c *cookieManager) Set(w http.ResponseWriter, value string) {

	// Create signature for the encoded value
	signature := c.createSignature([]byte(value))
	// Encode the signature to base64 URL-safe
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)
	cookie := &http.Cookie{
		Name:     c.options.name,
		Value:    value + "." + encodedSignature,
		Path:     c.options.path,
		Secure:   c.options.secure,
		HttpOnly: c.options.httpOnly,
		SameSite: c.options.sameSite,
	}
	http.SetCookie(w, cookie)
}

func (c *cookieManager) Get(r *http.Request) (string, error) {
	cookie, err := r.Cookie(c.options.name)
	if err != nil {
		switch errors.Is(err, http.ErrNoCookie) {
		case true:
			return "", ErrNoCookie
		default:
			return "", err
		}
	}
	v := strings.Split(cookie.Value, ".")
	if len(v) != 2 {
		return "", errors.New("invalid cookie")
	}
	signature, err := base64.RawURLEncoding.DecodeString(v[1])
	if err != nil {
		return "", err
	}

	if !c.verifySignature([]byte(v[0]), signature) {
		return "", errors.New("invalid signature")
	}

	return v[0], nil

}

func (c *cookieManager) createSignature(data []byte) []byte {
	mac := hmac.New(sha256.New, c.secret)
	mac.Write(data)
	return mac.Sum(nil)
}

func (c *cookieManager) verifySignature(data, signature []byte) bool {
	expectedSig := c.createSignature(data)
	return hmac.Equal(signature, expectedSig)
}
