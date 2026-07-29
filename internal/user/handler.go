package user

import (
	"context"
	"errors"
	"log"

	"connectrpc.com/connect"

	"project/internal/auth"
	userpb "project/internal/gen/user"
	"project/internal/gen/user/userpbconnect"
)

// Acting user always comes from auth.RequirePrincipal
type Handler struct {
	userpbconnect.UnimplementedUserServiceHandler
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProfile(
	ctx context.Context,
	_ *connect.Request[userpb.GetProfileRequest],
) (*connect.Response[userpb.GetProfileResponse], error) {
	principal, err := auth.RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := h.service.GetProfile(ctx, principal.PublicUUID)
	if errors.Is(err, ErrUserNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user: profile not found"))
	}
	if err != nil {
		log.Printf("user: get profile failed: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("user: request failed"))
	}

	return connect.NewResponse(&userpb.GetProfileResponse{
		Email:       profile.Email,
		DisplayName: profile.DisplayName,
	}), nil
}

func (h *Handler) UpdateProfile(
	ctx context.Context,
	req *connect.Request[userpb.UpdateProfileRequest],
) (*connect.Response[userpb.UpdateProfileResponse], error) {
	principal, err := auth.RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := h.service.UpdateProfile(ctx, principal.PublicUUID, req.Msg.GetDisplayName())
	if errors.Is(err, ErrDisplayNameTooLong) {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if errors.Is(err, ErrUserNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user: profile not found"))
	}
	if err != nil {
		log.Printf("user: update profile failed: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("user: request failed"))
	}

	return connect.NewResponse(&userpb.UpdateProfileResponse{
		Email:       profile.Email,
		DisplayName: profile.DisplayName,
	}), nil
}
