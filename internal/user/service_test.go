package user

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type fakeRepo struct {
	profiles map[string]Profile
}

func (f *fakeRepo) GetProfileByPublicUUID(_ context.Context, publicUUID string) (Profile, error) {
	profile, ok := f.profiles[publicUUID]
	if !ok {
		return Profile{}, ErrUserNotFound
	}
	return profile, nil
}

func (f *fakeRepo) UpdateDisplayName(_ context.Context, publicUUID, displayName string) (Profile, error) {
	profile, ok := f.profiles[publicUUID]
	if !ok {
		return Profile{}, ErrUserNotFound
	}
	profile.DisplayName = displayName
	f.profiles[publicUUID] = profile
	return profile, nil
}

const testPublicUUID = "00000000-0000-0000-0000-000000000001"

func newTestService(t *testing.T) *Service {
	t.Helper()
	repo := &fakeRepo{profiles: map[string]Profile{
		testPublicUUID: {Email: "admin@localhost", DisplayName: "Admin"},
	}}
	return NewService(repo)
}

func TestGetProfile(t *testing.T) {
	service := newTestService(t)

	profile, err := service.GetProfile(context.Background(), testPublicUUID)
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if profile.Email != "admin@localhost" || profile.DisplayName != "Admin" {
		t.Fatalf("profile = %+v, want email=admin@localhost displayName=Admin", profile)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	service := newTestService(t)

	_, err := service.GetProfile(context.Background(), "00000000-0000-0000-0000-000000000099")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("err = %v, want ErrUserNotFound", err)
	}
}

func TestUpdateProfile_Trims(t *testing.T) {
	service := newTestService(t)

	profile, err := service.UpdateProfile(context.Background(), testPublicUUID, "  Ada Lovelace  ")
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if profile.DisplayName != "Ada Lovelace" {
		t.Fatalf("displayName = %q, want %q", profile.DisplayName, "Ada Lovelace")
	}
}

func TestUpdateProfile_EmptyAllowed(t *testing.T) {
	service := newTestService(t)

	profile, err := service.UpdateProfile(context.Background(), testPublicUUID, "   ")
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if profile.DisplayName != "" {
		t.Fatalf("displayName = %q, want empty string", profile.DisplayName)
	}
}

func TestUpdateProfile_TooLong(t *testing.T) {
	service := newTestService(t)

	tooLong := "  " + strings.Repeat("a", maxDisplayNameLength+1) + "  "
	_, err := service.UpdateProfile(context.Background(), testPublicUUID, tooLong)
	if !errors.Is(err, ErrDisplayNameTooLong) {
		t.Fatalf("err = %v, want ErrDisplayNameTooLong", err)
	}
}

func TestUpdateProfile_AtMaxLength(t *testing.T) {
	service := newTestService(t)

	exact := strings.Repeat("a", maxDisplayNameLength)
	profile, err := service.UpdateProfile(context.Background(), testPublicUUID, exact)
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if profile.DisplayName != exact {
		t.Fatalf("displayName length = %d, want %d", len(profile.DisplayName), len(exact))
	}
}

func TestUpdateProfile_NotFound(t *testing.T) {
	service := newTestService(t)

	_, err := service.UpdateProfile(context.Background(), "00000000-0000-0000-0000-000000000099", "New Name")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("err = %v, want ErrUserNotFound", err)
	}
}
