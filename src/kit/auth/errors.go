package auth

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/kit/er"
	"net/http"
)

const (
	ErrCodeAccessTokenCreation           = "AUTH-001"
	ErrCodeSessionPasswordValidation     = "AUTH-002"
	ErrCodeUserNotFound                  = "AUTH-003"
	ErrCodeUserNotActive                 = "AUTH-004"
	ErrCodeUserLocked                    = "AUTH-005"
	ErrCodeSessionLoggedOut              = "AUTH-006"
	ErrCodeSessionAuthWrongSigningMethod = "AUTH-007"
	ErrCodeSessionAuthTokenExpired       = "AUTH-008"
	ErrCodeSessionAuthTokenInvalid       = "AUTH-009"
	ErrCodeSessionAuthTokenClaimsInvalid = "AUTH-010"
	ErrCodeSessionTokenInvalid           = "AUTH-011"
	ErrCodeSessionNoSessionFound         = "AUTH-012"
)

var (
	ErrAccessTokenCreation = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeAccessTokenCreation, "").C(ctx).Err()
	}
	ErrSessionPasswordValidation = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionPasswordValidation, "invalid password").Business().HttpSt(http.StatusUnauthorized).C(ctx).Err()
	}
	ErrUserNotFound = func(ctx context.Context, userId string) error {
		return er.WithBuilder(ErrCodeUserNotFound, "user not found").Business().F(er.FF{"userId": userId}).C(ctx).Err()
	}
	ErrUserNotActive = func(ctx context.Context, userId string) error {
		return er.WithBuilder(ErrCodeUserNotActive, "user not active").Business().F(er.FF{"userId": userId}).C(ctx).Err()
	}
	ErrUserLocked = func(ctx context.Context, userId string) error {
		return er.WithBuilder(ErrCodeUserLocked, "user locked").Business().F(er.FF{"userId": userId}).C(ctx).Err()
	}
	ErrSessionLoggedOut = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionLoggedOut, "session is logged out").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionAuthWrongSigningMethod = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthWrongSigningMethod, "wrong signing method").C(ctx).HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionAuthTokenExpired = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthTokenExpired, "token expired").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionAuthTokenInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthTokenInvalid, "invalid token").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionAuthTokenClaimsInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthTokenClaimsInvalid, "invalid token claims").Business().HttpSt(http.StatusUnauthorized).C(ctx).Err()
	}
	ErrSessionTokenInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionTokenInvalid, "session token is invalid").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSessionNoSessionFound = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionNoSessionFound, "no session found").C(ctx).Business().HttpSt(http.StatusUnauthorized).Err()
	}
)
