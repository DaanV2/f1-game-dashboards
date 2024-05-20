package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/DaanV2/f1-game-dashboards/server/pkg/randx"
	go_jwt "github.com/golang-jwt/jwt/v5"
)

type (
	MapClaims = go_jwt.MapClaims

	JwtService struct {
		signingKey  *SigningInfo
		signingKeys []*SigningInfo

		defaultClaims go_jwt.RegisteredClaims

		parseOptions []go_jwt.ParserOption
	}
)

// NewJwtService creates a new jwt signing and verification services
func NewJwtService(sigs []*SigningInfo) (*JwtService, error) {
	if len(sigs) == 0 {
		return nil, errors.New("no signing methods provided")
	}
	var signingKey *SigningInfo
	methods := make([]string, len(sigs))
	for i, s := range sigs {
		if s.PrivateKey != nil {
			signingKey = s
		}

		methods[i] = s.Method.Alg()
	}

	if signingKey == nil {
		return nil, errors.New("no key provided that can be used to sign")
	}

	result := &JwtService{
		signingKey:  signingKey,
		signingKeys: sigs,

		defaultClaims: go_jwt.RegisteredClaims{
			Audience: []string{"f1-game-dashboards"},
			Issuer:   "f1-game-dashboards",
		},
	}

	result.parseOptions = []go_jwt.ParserOption{
		go_jwt.WithValidMethods(methods),
		go_jwt.WithLeeway(time.Minute * 5),
		go_jwt.WithExpirationRequired(),
		go_jwt.WithIssuer(result.defaultClaims.Issuer),
	}

	return result, nil
}

// Sign generates a JWT from the given claims
func (j *JwtService) Sign(customClaims MapClaims) (string, error) {
	key := j.GetSigningKey()
	now := time.Now()
	jti, err := randx.GenerateBase64(36)
	if err != nil {
		return "", fmt.Errorf("error generating jti: %w", err)
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	claims := combineClaims(
		// Overridable properties
		go_jwt.MapClaims{
			"sub": "unspecified",
			"jti": jti,
		},
		// User specified
		customClaims,
		// These cannot be overriden
		go_jwt.MapClaims{
			"exp": now.Add(time.Hour * 24).Unix(),
			"nbf": now.Add(time.Second * -5).Unix(),
			"iat": now.Unix(),

			"aud": j.defaultClaims.Audience,
			"iss": j.defaultClaims.Issuer,
		},
	)

	token := go_jwt.NewWithClaims(key.Method, claims)
	token.Header["kid"] = key.KeyID

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(key.PrivateKey)
}

// Verify tries to parse and validates the token
func (j *JwtService) Verify(tokenStr string) (*go_jwt.Token, error) {
	// Parse and verify the token
	token, err := go_jwt.Parse(tokenStr, j.getKey, j.parseOptions...)

	if err != nil {
		return token, err
	}
	if !token.Valid {
		return token, errors.New("token is invalid")
	}

	return token, err
}

// Refresh checks if the token is still valid, and then generate a new JWT
func (j *JwtService) Refresh(tokenStr string) (string, error) {
	token, err := j.Verify(tokenStr)
	if err != nil {
		return "", err
	}
	claims := go_jwt.MapClaims{}

	switch c := token.Claims.(type) {
	case go_jwt.MapClaims:
		for k, v := range c {
			claims[k] = v
		}
	default:
		claims["sub"], err = c.GetSubject()
		if err != nil {
			return "", err
		}
	}

	return j.Sign(claims)
}

// GetSigningKey return the key used to sign
func (j *JwtService) GetSigningKey() *SigningInfo {
	return j.signingKey
}

// getKey returns either VerificationKey or VerificationKeySet
func (j *JwtService) getKey(token *go_jwt.Token) (interface{}, error) {
	keys := make([]go_jwt.VerificationKey, 0)
	kid, ok := token.Header["kid"]
	for _, v := range j.signingKeys {
		if ok {
			if v.KeyID == kid {
				keys = append(keys, v.PublicKey)
			}
		} else {
			keys = append(keys, v.PublicKey)
		}
	}

	// If no keys, atleast provide the signing key
	if len(keys) == 0 {
		keys = append(keys, j.GetSigningKey().PublicKey)
	}

	return go_jwt.VerificationKeySet{Keys: keys}, nil
}

// combineClaims combine multiple claims into one
func combineClaims(claims ...MapClaims) MapClaims {
	result := MapClaims{}

	for _, c := range claims {
		for k, v := range c {
			result[k] = v
		}
	}

	return result
}
