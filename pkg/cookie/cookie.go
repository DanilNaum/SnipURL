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
		secure:   false,
		httpOnly: false,
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

// Set creates a new cookie with the given value and sets it in the HTTP response.
// The cookie value is signed using HMAC-SHA256 to ensure integrity.
// The signature is base64 URL-encoded and appended to the value with a dot separator.
// The cookie settings (name, path, secure, httpOnly, sameSite) are taken from the cookieManager options.
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

// Get retrieves and validates a cookie from the HTTP request.
// It returns the cookie value and any error encountered.
// The method performs the following steps:
// 1. Retrieves the cookie by name from the request
// 2. Returns ErrNoCookie if the cookie is not found
// 3. Splits the cookie value into the actual value and signature parts
// 4. Validates the signature using HMAC-SHA256
// 5. Returns the original value if signature is valid
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
