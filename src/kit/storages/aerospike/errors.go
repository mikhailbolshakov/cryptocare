package aerospike

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/kit/er"
)

var (
	ErrCodeAeroConn           = "AERO-001"
	ErrCodeAeroClosed         = "AERO-002"
	ErrCodeAeroNewKey         = "AERO-003"
	ErrCodeAeroInvalidBinType = "AERO-004"
)

var (
	ErrAeroConn = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeAeroConn, "").C(ctx).Err()
	}
	ErrAeroNewKey = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeAeroNewKey, "").C(ctx).Err()
	}
	ErrAeroClosed = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeAeroClosed, "dealing with closed instance").C(ctx).Err()
	}
	ErrAeroInvalidBinType = func(ctx context.Context, bin string) error {
		return er.WithBuilder(ErrCodeAeroInvalidBinType, "invalid bin type").F(er.FF{"bin": bin}).C(ctx).Err()
	}
)
