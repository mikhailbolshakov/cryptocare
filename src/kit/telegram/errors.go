package telegram

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/kit/er"
)

var (
	ErrCodeTelegramBotEmpty      = "TLG-001"
	ErrCodeTelegramRequestFailed = "TLG-002"
	ErrCodeTelegramResponseError = "TLG-003"
)

var (
	ErrTelegramBotEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeTelegramBotEmpty, "bot empty").Business().C(ctx).Err()
	}
	ErrTelegramRequestFailed = func(ctx context.Context, cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeTelegramRequestFailed, "").C(ctx).Err()
	}
	ErrTelegramResponseError = func(ctx context.Context, status, body string) error {
		return er.WithBuilder(ErrCodeTelegramResponseError, "telegram error").F(er.FF{"status": status, "body": body}).C(ctx).Err()
	}
)
