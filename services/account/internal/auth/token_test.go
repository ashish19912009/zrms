package auth_test

import (
	"testing"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	issuer := auth.NewJWTIssuer("test-secret-key", time.Minute*15)

	accountID := "test-account-id"
	role := "admin"

	tokenStr, err := issuer.GenerateToken(accountID, role)
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)
}
