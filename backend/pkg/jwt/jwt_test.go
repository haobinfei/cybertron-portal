package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret"
const testExpireHours = 24

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(1, "admin", "admin", testSecret, testExpireHours)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "admin", claims.Username)
	assert.Equal(t, "admin", claims.Role)
}

func TestParseToken_Valid(t *testing.T) {
	token, err := GenerateToken(2, "user1", "user", testSecret, testExpireHours)
	require.NoError(t, err)

	claims, err := ParseToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, uint(2), claims.UserID)
	assert.Equal(t, "user1", claims.Username)
	assert.Equal(t, "user", claims.Role)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestParseToken_InvalidSecret(t *testing.T) {
	token, err := GenerateToken(1, "admin", "admin", testSecret, testExpireHours)
	require.NoError(t, err)

	_, err = ParseToken(token, "wrong-secret")
	assert.Error(t, err)
}

func TestParseToken_InvalidFormat(t *testing.T) {
	_, err := ParseToken("not.a.valid.token", testSecret)
	assert.Error(t, err)
}

func TestParseToken_EmptyString(t *testing.T) {
	_, err := ParseToken("", testSecret)
	assert.Error(t, err)
}

func TestParseToken_Expired(t *testing.T) {
	claims := Claims{
		UserID:   1,
		Username: "admin",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)

	_, err = ParseToken(tokenString, testSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token has invalid claims")
}

func TestGenerateToken_DifferentUsers(t *testing.T) {
	t1, err := GenerateToken(10, "alice", "admin", testSecret, testExpireHours)
	require.NoError(t, err)

	t2, err := GenerateToken(20, "bob", "user", testSecret, testExpireHours)
	require.NoError(t, err)

	assert.NotEqual(t, t1, t2)
}
