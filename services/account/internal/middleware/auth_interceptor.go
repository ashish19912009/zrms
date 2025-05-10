package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/ashish19912009/zrms/services/account/internal/client"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	allowedRoles map[string]bool
	jwtSecret    []byte
	authz_client client.AuthZClient
	authn_client client.AuthNClient
	//	validateToken func(token string) (map[string]interface{}, error)
}

func NewAuthInterceptor(secret string, roles ...string) *AuthInterceptor {
	roleMap := make(map[string]bool)
	for _, role := range roles {
		roleMap[role] = true
	}
	return &AuthInterceptor{allowedRoles: roleMap, jwtSecret: []byte(secret)}
	//	return &AuthInterceptor{validateToken: validateToken}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		role, err := a.extractAndValidateToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("unauthorized %w", err)
		}
		if !a.allowedRoles[role] {
			return nil, fmt.Errorf("forbidden: role '%s' is not allowed", role)
		}
		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) extractAndValidateToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	tokenStr := parts[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// verify singing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token class")
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return "", fmt.Errorf("role not found in token")
	}
	return role, nil
}
