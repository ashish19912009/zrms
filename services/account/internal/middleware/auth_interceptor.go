package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ashish19912009/zrms/services/account/internal/client"
	"github.com/ashish19912009/zrms/services/account/internal/helper"
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
		fmt.Printf("rpc name: %s", info.FullMethod)
		claims, err := a.extractAndValidateToken(ctx)
		if err != nil || claims == nil {
			logger.Error("authentication failed", err, map[string]interface{}{
				"layer":  "middleware",
				"method": "Unary",
			})
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed")
		}
		internalCall, rPerm, err := a.extractAndValidatePermission(ctx)
		if err != nil || rPerm == nil {
			logger.Error("no permission request", err, map[string]interface{}{
				"layer":  "middleware",
				"method": "Unary",
			})
			return nil, status.Errorf(codes.PermissionDenied, "no permission requested")
		}
		ctx = context.WithValue(ctx, model.RequestContextKey, &model.RequestContext{
			Claims: claims,
		})
		// bypass Authz if service internal call
		if internalCall != nil && *internalCall {
			return handler(ctx, req)
		}
		isAuthorized, err := a.authz_client.CheckAccess(ctx, claims.RegisteredClaims.Subject, claims.FranchiseID, rPerm.Resource, rPerm.Action)
		if err != nil {
			return nil, err
		}
		if !isAuthorized.Allowed {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}
		if claims.RegisteredClaims.ExpiresAt.AsTime().Before(time.Now()) {
			return nil, status.Error(codes.Unauthenticated, "token expired")
		}
		if claims.RegisteredClaims.Subject == "" || claims.AccountType == "" {
			return nil, status.Error(codes.PermissionDenied, "account not verified")
		}
		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) extractAndValidateToken(ctx context.Context) (*model.AuthClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata in context")
	}

	authHeader := helper.GetMetadataValue(md, "authorization", "Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	authClaims, err := a.authn_client.VerifyToken(ctx, model.Token{Token: token})
	if err != nil {
		logger.Error("token verification failed", err, map[string]interface{}{
			"token": "[redacted]", // Avoid logging the actual token in production
		})
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	return authClaims, nil
}

func (a *AuthInterceptor) extractAndValidatePermission(ctx context.Context) (*bool, *model.Permission, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, nil, fmt.Errorf("missing metadata in context")
	}

	resource := helper.GetMetadataValue(md, "resource", "Resource")
	if resource == "" {
		return nil, nil, fmt.Errorf("missing resource in metadata")
	}

	action := helper.GetMetadataValue(md, "action", "Action")
	if action == "" {
		return nil, nil, fmt.Errorf("missing action in metadata")
	}
	isInternal := helper.GetMetadataValue(md, "trusted", "Trusted")
	if isInternal != "" && isInternal == "true" {
		b := true
		return &b, nil, nil
	}
	b := false
	return &b, &model.Permission{
		Resource: resource,
		Action:   action,
	}, nil
}
