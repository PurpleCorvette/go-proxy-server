// Package auth provides a gRPC unary interceptor for JWT authentication.
package auth

import (
	"context"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authenticator defines an interface for JWT authentication.
type Authenticator interface {
	AuthUnaryInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor
}

// JWTAuthenticator implements Authenticator interface for JWT authentication.
type JWTAuthenticator struct{}

// NewJWTAuthenticator returns a new instance of JWTAuthenticator.
func NewJWTAuthenticator() *JWTAuthenticator {
	return &JWTAuthenticator{}
}

// signingKey is used to validate the JWT token signature.
var signingKey = []byte(viper.GetString("jwt_secret"))

// AuthUnaryInterceptor returns a gRPC unary interceptor for JWT authentication.
func (a *JWTAuthenticator) AuthUnaryInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		token := getTokenFromMetadata(md)
		if token == "" {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, status.Errorf(codes.Unauthenticated, "unexpected singing method: %v", token.Header["alg"])
			}
			return signingKey, nil
		})
		if err != nil {
			logger.Warnf("Unauthorized access attempt: %v", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

// getTokenFromMetadata extracts the JWT token from the metadata.
func getTokenFromMetadata(md metadata.MD) string {
	values := md["authorization"]
	if len(values) == 0 {
		return ""
	}
	token := values[0]
	if strings.HasPrefix(token, "Bearer ") {
		return token[7:]
	}
	return token
}
