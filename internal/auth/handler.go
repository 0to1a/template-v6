package auth

import (
	"context"
	"errors"
	"log"

	"connectrpc.com/connect"

	authpb "project/internal/gen/auth"
	"project/internal/gen/auth/authpbconnect"
)

type Handler struct {
	authpbconnect.UnimplementedAuthServiceHandler
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RequestLogin(
	ctx context.Context,
	req *connect.Request[authpb.RequestLoginRequest],
) (*connect.Response[authpb.RequestLoginResponse], error) {
	if err := h.service.RequestLogin(ctx, req.Msg.GetEmail()); err != nil {
		log.Printf("auth: request login failed: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("auth: request failed"))
	}
	return connect.NewResponse(&authpb.RequestLoginResponse{}), nil
}

func (h *Handler) SubmitLogin(
	ctx context.Context,
	req *connect.Request[authpb.SubmitLoginRequest],
) (*connect.Response[authpb.SubmitLoginResponse], error) {
	token, err := h.service.SubmitLogin(ctx, req.Msg.GetEmail(), req.Msg.GetCode())
	if errors.Is(err, errUnauthenticated) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid email or code"))
	}
	if err != nil {
		log.Printf("auth: submit login failed: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("auth: login failed"))
	}
	return connect.NewResponse(&authpb.SubmitLoginResponse{AccessToken: token}), nil
}
