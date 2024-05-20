package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"

	"github.com/DaanV2/f1-game-dashboards/server/config"
	"github.com/DaanV2/f1-game-dashboards/server/pkg/randx"
	"github.com/charmbracelet/log"
	go_jwt "github.com/golang-jwt/jwt/v5"
)

var (
	_ json.Marshaler   = &SigningInfo{}
	_ json.Unmarshaler = &SigningInfo{}
)

type SigningInfo struct {
	KeyID      string
	Method     go_jwt.SigningMethod
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

func (s *SigningInfo) MarshalJSON() ([]byte, error) {
	privateKey, err := x509.MarshalPKCS8PrivateKey(s.PrivateKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := x509.MarshalPKIXPublicKey(s.PublicKey)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"kid":     s.KeyID,
		"alg":     s.Method.Alg(),
		"private": privateKey,
		"public":  publicKey,
	})
}

func (s *SigningInfo) UnmarshalJSON(data []byte) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	key, aErr := config.Get[string](m, "kid")
	method, bErr := config.Get[string](m, "alg")
	private, cErr := config.Get[string](m, "private")
	public, dErr := config.Get[string](m, "public")
	err := errors.Join(aErr, bErr, cErr, dErr)
	if err != nil {
		return err
	}

	s.KeyID = key
	if m := go_jwt.GetSigningMethod(method); m != nil {
		s.Method = m
	} else {
		return errors.New("method not found")
	}

	s.PrivateKey, aErr = x509.ParsePKCS8PrivateKey([]byte(private))
	s.PublicKey, bErr = x509.ParsePKIXPublicKey([]byte(public))
	return errors.Join(aErr, bErr)
}

// GenerateSigningInfo creates a new set of SigningInfo using RSA
func GenerateSigningInfo() (*SigningInfo, error) {
	key, err := randx.GenerateBase64(32)
	if err != nil {
		return nil, err
	}

	log.Debug("generating signing info", "keyid", key)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &SigningInfo{
		KeyID:      key,
		Method:     go_jwt.SigningMethodRS512,
		PrivateKey: privateKey,
		PublicKey:  privateKey.Public(),
	}, nil
}
