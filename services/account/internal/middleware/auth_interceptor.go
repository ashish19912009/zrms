package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/client"
	"github.com/ashish19912009/zrms/services/account/internal/logger"
	"github.com/ashish19912009/zrms/services/account/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authz_client client.AuthZClient
	authn_client client.AuthNClient
	//	validateToken func(token string) (map[string]interface{}, error)
}

func NewAuthInterceptor(authZ client.AuthZClient, authN client.AuthNClient) *AuthInterceptor {
	return &AuthInterceptor{authz_client: authZ, authn_client: authN}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		claims, err := a.extractAndValidateToken(ctx)
		if err != nil || claims == nil {
			logger.Error("authentication failed", err, map[string]interface{}{
				"layer":  "middleware",
				"method": "Unary",
			})
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed")
		}
		if claims.RegisteredClaims.ExpiresAt.AsTime().Before(time.Now()) {
			return nil, status.Error(codes.Unauthenticated, "token expired")
		}
		if claims.RegisteredClaims.Subject == "" || claims.AccountType == "" {
			return nil, status.Error(codes.PermissionDenied, "account not verified")
		}
		// Add claims to context
		ctx = context.WithValue(ctx, model.RequestContextKey, &model.RequestContext{
			Claims: claims,
		})
		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) extractAndValidateToken(ctx context.Context) (*model.AuthClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		authHeader = md.Get("Authorization")
	}
	if len(authHeader) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	tokenStr := model.Token{
		Token: strings.TrimSpace(parts[1]),
	}
	if len(tokenStr.Token) == 0 {
		return nil, fmt.Errorf("empty token")
	}

	authClaims, err := a.authn_client.VerifyToken(ctx, tokenStr)
	if err != nil {
		logger.Error("something went wrong while verifying token", err, nil)
		return nil, err
	}
	return authClaims, nil
}
