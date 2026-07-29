package auth

import (
	"context"
	"errors"
	"time"
)

// Same generic error for unknown users and bad codes
var errUnauthenticated = errors.New("auth: invalid credentials")

type Service struct {
	repo                Repository
	delivery            LoginCodeSender
	jwtManager          *JWTManager
	now                 func() time.Time
	isGuestRegistration bool
}

// TOTP secret derived from jwtManager's signing secret
func NewService(repo Repository, delivery LoginCodeSender, jwtManager *JWTManager, isGuestRegistration bool) *Service {
	return &Service{
		repo:                repo,
		delivery:            delivery,
		jwtManager:          jwtManager,
		now:                 time.Now,
		isGuestRegistration: isGuestRegistration,
	}
}

// Never reveals whether the account already existed
func (s *Service) RequestLogin(ctx context.Context, email string) error {
	normalized := normalizeEmail(email)

	user, err := s.repo.GetActiveUserByEmail(ctx, normalized)
	if errors.Is(err, ErrUserNotFound) {
		if !s.isGuestRegistration {
			return nil
		}
		user, err = s.repo.CreateUser(ctx, normalized)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	code := s.jwtManager.generateLoginCode(user.PublicUUID, normalized, s.now())
	// Failure not surfaced: must not leak account existence
	_ = s.delivery.SendLoginCode(ctx, user.Email, code)
	return nil
}

// Unknown users and bad codes return the same error
func (s *Service) SubmitLogin(ctx context.Context, email, code string) (string, error) {
	normalized := normalizeEmail(email)

	user, err := s.repo.GetActiveUserByEmail(ctx, normalized)
	if errors.Is(err, ErrUserNotFound) {
		return "", errUnauthenticated
	}
	if err != nil {
		return "", err
	}

	if !s.jwtManager.verifyLoginCode(user.PublicUUID, normalized, code, s.now()) {
		return "", errUnauthenticated
	}

	token, err := s.jwtManager.Issue(user.PublicUUID)
	if err != nil {
		return "", err
	}
	return token, nil
}
