package jwt_test

import (
	"strings"
	"testing"
	"time"

	"github.com/DaanV2/f1-game-dashboards/server/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_JwtManager(t *testing.T) {
	sigInfo, err := jwt.GenerateSigningInfo()
	require.NoError(t, err)

	jwtManager, err := jwt.NewJwtService([]*jwt.SigningInfo{sigInfo})
	require.NoError(t, err)

	t.Run("can generate a new jwt", func(t *testing.T) {
		token, err := jwtManager.Sign(jwt.MapClaims{
			"foo": "bar",
		})
		require.NoError(t, err)

		dots := strings.Count(token, ".")
		require.Equal(t, dots, int(2))

		decoded, err := jwtManager.Verify(token)
		require.NoError(t, err)

		assert.Contains(t, decoded.Header, "typ")
		assert.Contains(t, decoded.Header, "alg")
		assert.Contains(t, decoded.Header, "kid")

		assert.Contains(t, decoded.Claims, "jti")
		assert.Contains(t, decoded.Claims, "foo")
	})

	t.Run("can generate a token, but registered claims are not overriden except sub, jti", func(t *testing.T) {
		now := time.Now().Add(time.Hour * 24 * 14)
		customClaims := jwt.MapClaims{
			"sub": "specific subject",
			"jti": "specific jti",
			"exp": now.Add(time.Hour * 24),
			"nbf": now.Add(time.Hour * 24),
			"iat": now.Add(time.Hour * 24),
			"aud": "specific audience",
			"iss": "specific issuer",
		}

		token, err := jwtManager.Sign(customClaims)
		require.NoError(t, err)
		decoded, err := jwtManager.Verify(token)
		require.NoError(t, err)

		mapClaims, ok := decoded.Claims.(jwt.MapClaims)
		require.True(t, ok)

		assert.Equal(t, mapClaims["sub"], customClaims["sub"])
		assert.Equal(t, mapClaims["jti"], customClaims["jti"])

		assert.NotEqual(t, mapClaims["exp"], customClaims["exp"])
		assert.NotEqual(t, mapClaims["nbf"], customClaims["nbf"])
		assert.NotEqual(t, mapClaims["iat"], customClaims["iat"])
		assert.NotEqual(t, mapClaims["aud"], customClaims["aud"])
		assert.NotEqual(t, mapClaims["iss"], customClaims["iss"])
	})

	t.Run("can refresh the token, will keep the jti and custom claims", func(t *testing.T) {
		customClaims := jwt.MapClaims{
			"sub": "steve",
			"foo": "bar",
			"temp": float64(0),
			"jti": "specific jti",
		}

		jwtToken, err := jwtManager.Sign(customClaims)
		require.NoError(t, err)
		refreshed, err := jwtManager.Refresh(jwtToken)
		require.NoError(t, err)
		decoded, err := jwtManager.Verify(refreshed)
		require.NoError(t, err)

		decodedMap, ok := decoded.Claims.(jwt.MapClaims)
		require.True(t, ok)

		for k, v := range customClaims {
			require.Equal(t, decodedMap[k], v, k)
		}
	})
}
