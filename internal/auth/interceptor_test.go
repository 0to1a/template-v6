package auth

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"

	authpb "project/internal/gen/auth"
)

// Unlisted procedures must be denied by default
func TestInterceptor_DefaultDeny(t *testing.T) {
	m := newTestJWTManager(t, fixedTime)
	interceptor := NewInterceptor(m, map[string]bool{})

	next := connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		t.Fatal("handler reached without authentication")
		return nil, nil
	})

	req := connect.NewRequest(&authpb.RequestLoginRequest{})
	_, err := interceptor.WrapUnary(next)(context.Background(), req)

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) || connectErr.Code() != connect.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated, got %v", err)
	}
}

func TestInterceptor_ValidTokenAttachesPrincipal(t *testing.T) {
	m := newTestJWTManager(t, fixedTime)
	interceptor := NewInterceptor(m, map[string]bool{})

	token, err := m.Issue(testUUID)
	if err != nil {
		t.Fatal(err)
	}

	var got Principal
	next := connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		principal, err := RequirePrincipal(ctx)
		if err != nil {
			return nil, err
		}
		got = principal
		return connect.NewResponse(&authpb.RequestLoginResponse{}), nil
	})

	req := connect.NewRequest(&authpb.RequestLoginRequest{})
	req.Header().Set("Authorization", "Bearer "+token)

	if _, err := interceptor.WrapUnary(next)(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	if got.PublicUUID != testUUID {
		t.Fatalf("principal = %s, want %s", got.PublicUUID, testUUID)
	}
}

func TestInterceptor_MalformedTokenRejected(t *testing.T) {
	m := newTestJWTManager(t, fixedTime)
	interceptor := NewInterceptor(m, map[string]bool{})

	next := connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		t.Fatal("handler reached with a malformed token")
		return nil, nil
	})

	req := connect.NewRequest(&authpb.RequestLoginRequest{})
	req.Header().Set("Authorization", "Bearer not-a-jwt")

	if _, err := interceptor.WrapUnary(next)(context.Background(), req); err == nil {
		t.Fatal("malformed token was accepted")
	}
}
