package auth

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
)

const bearerPrefix = "Bearer "

var errMissingOrInvalidToken = errors.New("unauthenticated")

// Paths in publicProcedures skip auth; all else needs it
func NewInterceptor(jwtManager *JWTManager, publicProcedures map[string]bool) connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if publicProcedures[req.Spec().Procedure] {
				return next(ctx, req)
			}

			token, ok := bearerToken(req.Header().Get("Authorization"))
			if !ok {
				return nil, connect.NewError(connect.CodeUnauthenticated, errMissingOrInvalidToken)
			}

			principal, err := jwtManager.Parse(token)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errMissingOrInvalidToken)
			}

			return next(WithPrincipal(ctx, principal), req)
		}
	})
}

// Use over PrincipalFromContext for consistent errors
func RequirePrincipal(ctx context.Context) (Principal, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok {
		return Principal{}, connect.NewError(connect.CodeUnauthenticated, errMissingOrInvalidToken)
	}
	return principal, nil
}

func bearerToken(header string) (string, bool) {
	if !strings.HasPrefix(header, bearerPrefix) {
		return "", false
	}
	token := strings.TrimPrefix(header, bearerPrefix)
	if token == "" {
		return "", false
	}
	return token, true
}
