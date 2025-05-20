// pkg/middleware/jwt_interceptor.go
package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type JWTInterceptor struct {
	jwkSetURL string
	keySet    jwk.Set
	lastFetch time.Time
}

func NewJWTInterceptor(jwkSetURL string) (*JWTInterceptor, error) {
	keySet, err := jwk.Fetch(context.Background(), jwkSetURL, jwk.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, err
	}
	return &JWTInterceptor{
		jwkSetURL: jwkSetURL,
		keySet:    keySet,
		lastFetch: time.Now(),
	}, nil
}

func (ji *JWTInterceptor) refreshIfNeeded() {
	if time.Since(ji.lastFetch) > 10*time.Minute {
		keySet, err := jwk.Fetch(context.Background(), ji.jwkSetURL, jwk.WithHTTPClient(http.DefaultClient))
		if err == nil {
			ji.keySet = keySet
			ji.lastFetch = time.Now()
		}
	}
}

func (ji *JWTInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ji.refreshIfNeeded()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, errors.New("authorization header missing")
		}

		tokenStr := authHeader[0]
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		token, err := jwt.ParseString(tokenStr, jwt.WithKeySet(ji.keySet))
		if err != nil {
			//fmt.Print("Error from Praseing")
			return nil, err
		}

		// You may add custom claim validations here (e.g., expiration, issuer, roles)
		ctx = context.WithValue(ctx, "user", token)
		return handler(ctx, req)
	}
}
