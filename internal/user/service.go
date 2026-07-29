package user

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"
)

// Empty string is allowed: means no display name set
const maxDisplayNameLength = 100

var ErrDisplayNameTooLong = errors.New("user: display_name is too long")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfile(ctx context.Context, publicUUID string) (Profile, error) {
	return s.repo.GetProfileByPublicUUID(ctx, publicUUID)
}

// Empty trimmed value clears the name; over max is rejected
func (s *Service) UpdateProfile(ctx context.Context, publicUUID, displayName string) (Profile, error) {
	trimmed := strings.TrimSpace(displayName)
	if utf8.RuneCountInString(trimmed) > maxDisplayNameLength {
		return Profile{}, ErrDisplayNameTooLong
	}
	return s.repo.UpdateDisplayName(ctx, publicUUID, trimmed)
}
