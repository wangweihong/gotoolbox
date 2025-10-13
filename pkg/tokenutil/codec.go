package tokenutil

import (
	"crypto/rsa"
	"fmt"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Options struct {
	URL url.URL
	Key *rsa.PrivateKey
}

var MaxIssueDelay = time.Second * 90

func DefaultJWTTrackedRequestCodec(opts Options) JWTTrackedRequestCodec {
	return JWTTrackedRequestCodec{
		SigningMethod: DefaultJWTSigningMethod,
		Audience:      opts.URL.String(),
		Issuer:        opts.URL.String(),
		MaxAge:        MaxIssueDelay,
		Key:           opts.Key,
	}
}

type TrackedRequest struct {
	// 随机ID, 用于追踪每个请求
	Index string `json:"-"`
	Value any    `json:"value"`
}

// TrackedRequestCodec handles encoding and decoding of a TrackedRequest.
type TrackedRequestCodec interface {
	// Encode returns an encoded string representing the TrackedRequest.
	Encode(value TrackedRequest) (string, error)

	// Decode returns a Tracked request from an encoded string.
	Decode(signed string) (*TrackedRequest, error)
}

var DefaultJWTSigningMethod = jwt.SigningMethodRS256

// JWTTrackedRequestCodec encodes TrackedRequests as signed JWTs
type JWTTrackedRequestCodec struct {
	SigningMethod jwt.SigningMethod
	Audience      string
	Issuer        string
	MaxAge        time.Duration
	Key           *rsa.PrivateKey
}

var _ TrackedRequestCodec = JWTTrackedRequestCodec{}

// JWTClaims represents the JWT claims for a tracked request.
type JWTClaims struct {
	jwt.RegisteredClaims
	TrackedRequest
}

// Encode returns an encoded string representing the TrackedRequest.
func (s JWTTrackedRequestCodec) Encode(value TrackedRequest) (string, error) {
	now := time.Now().UTC()
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{s.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(s.MaxAge)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    s.Issuer,
			NotBefore: jwt.NewNumericDate(now),
			Subject:   value.Index,
		},
		TrackedRequest: value,
	}
	token := jwt.NewWithClaims(s.SigningMethod, claims)
	return token.SignedString(s.Key)
}

// Decode returns a Tracked request from an encoded string.
func (s JWTTrackedRequestCodec) Decode(signed string) (*TrackedRequest, error) {
	parser := jwt.Parser{
		ValidMethods: []string{s.SigningMethod.Alg()},
	}
	claims := JWTClaims{}
	_, err := parser.ParseWithClaims(signed, &claims, func(*jwt.Token) (any, error) {
		return s.Key.Public(), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.VerifyAudience(s.Audience, true) {
		return nil, fmt.Errorf("expected audience %q, got %q", s.Audience, claims.Audience)
	}
	if !claims.VerifyIssuer(s.Issuer, true) {
		return nil, fmt.Errorf("expected issuer %q, got %q", s.Issuer, claims.Issuer)
	}
	claims.TrackedRequest.Index = claims.Subject
	return &claims.TrackedRequest, nil
}
