package http

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/kit/er"
	"net/http"
)

const (
	ErrCodeHttpTest                                                        = "HTTP-000"
	ErrCodeHttpSrvListen                                                   = "HTTP-001"
	ErrCodeDecodeRequest                                                   = "HTTP-002"
	ErrCodeHttpUrlVar                                                      = "HTTP-003"
	ErrCodeHttpCurrentUser                                                 = "HTTP-004"
	ErrCodeHttpUrlVarEmpty                                                 = "HTTP-005"
	ErrCodeHttpUrlFormVarEmpty                                             = "HTTP-006"
	ErrCodeHttpUrlFormVarNotInt                                            = "HTTP-007"
	ErrCodeHttpUrlFormVarNotTime                                           = "HTTP-008"
	ErrCodeHttpMultipartParseForm                                          = "HTTP-009"
	ErrCodeHttpMultipartEmptyContent                                       = "HTTP-010"
	ErrCodeHttpMultipartNotMultipart                                       = "HTTP-011"
	ErrCodeHttpMultipartParseMediaType                                     = "HTTP-012"
	ErrCodeHttpMultipartWrongMediaType                                     = "HTTP-013"
	ErrCodeHttpMultipartMissingBoundary                                    = "HTTP-014"
	ErrCodeHttpMultipartEofReached                                         = "HTTP-015"
	ErrCodeHttpMultipartNext                                               = "HTTP-016"
	ErrCodeHttpMultipartFormNameFileExpected                               = "HTTP-017"
	ErrCodeHttpMultipartFilename                                           = "HTTP-018"
	ErrCodeHttpCurrentClient                                               = "HTTP-019"
	ErrCodeHttpUrlFormVarNotFloat                                          = "HTTP-020"
	ErrCodeHttpUrlFormVarNotBool                                           = "HTTP-021"
	ErrCodeHttpUrlWrongSortFormat                                          = "HTTP-022"
	ErrCodeHttpUrlVarInvalidUUID                                           = "HTTP-023"
	ErrCodeHttpUrlMaxPageSizeExceeded                                      = "HTTP-024"
	ErrCodeRouteBuilderUrlEmpty                                            = "HTTP-025"
	ErrCodeRouteBuilderVerbEmpty                                           = "HTTP-026"
	ErrCodeRouteBuilderBothHandleFuncAndHandlerEmpty                       = "HTTP-027"
	ErrCodeRouteBuilderAuthorizationPoliciesSpecifiedWithoutAuthentication = "HTTP-028"
	ErrCodeRouteBuilderSpecialMiddlewaresRequireSubrouting                 = "HTTP-029"
	ErrCodeUserLoginFail                                                   = "HTTP-030"
	ErrCodeSecurityPermissionsDenied                                       = "HTTP-031"
	ErrCodeSecurityNoResourceRequested                                     = "HTTP-032"
)

var (
	ErrHttpTest          = func() error { return er.WithBuilder(ErrCodeHttpTest, "").Business().Err() }
	ErrHttpSrvListen     = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeHttpSrvListen, "").Err() }
	ErrHttpDecodeRequest = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeDecodeRequest, "invalid request").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlVar = func(ctx context.Context, v string) error {
		return er.WithBuilder(ErrCodeHttpUrlVar, "invalid or empty URL parameter").F(er.FF{"var": v}).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpCurrentUser = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeHttpCurrentUser, `cannot obtain current user`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlVarEmpty = func(ctx context.Context, v string) error {
		return er.WithBuilder(ErrCodeHttpUrlVarEmpty, `URL parameter is empty`).Business().F(er.FF{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlVarInvalidUUID = func(ctx context.Context, v string) error {
		return er.WithBuilder(ErrCodeHttpUrlVarInvalidUUID, `invalid UUID`).Business().F(er.FF{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarEmpty = func(ctx context.Context, v string) error {
		return er.WithBuilder(ErrCodeHttpUrlFormVarEmpty, `URL form value is empty`).Business().F(er.FF{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarNotInt = func(cause error, ctx context.Context, v string) error {
		return er.WrapWithBuilder(cause, ErrCodeHttpUrlFormVarNotInt, "form value must be of int type").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarNotFloat = func(cause error, ctx context.Context, v string) error {
		return er.WrapWithBuilder(cause, ErrCodeHttpUrlFormVarNotFloat, "form value must be of float type").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarNotBool = func(cause error, ctx context.Context, v string) error {
		return er.WrapWithBuilder(cause, ErrCodeHttpUrlFormVarNotBool, "form value must be of bool type").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlFormVarNotTime = func(cause error, ctx context.Context, v string) error {
		return er.WrapWithBuilder(cause, ErrCodeHttpUrlFormVarNotTime, "form value must be of time type in RFC-3339 format").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartParseForm = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeHttpMultipartParseForm, "parse multipart form").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartEmptyContent = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeHttpMultipartEmptyContent, `content is empty`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartNotMultipart = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeHttpMultipartNotMultipart, `content isn't multipart`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartParseMediaType = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeHttpMultipartParseMediaType, "parse media type").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartWrongMediaType = func(ctx context.Context, mt string) error {
		return er.WithBuilder(ErrCodeHttpMultipartWrongMediaType, `wrong media type %s`, mt).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartMissingBoundary = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeHttpMultipartMissingBoundary, `missing boundary`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartEofReached = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeHttpMultipartEofReached, `no parts found`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartNext = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeHttpMultipartNext, "reading part").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartFormNameFileExpected = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeHttpMultipartFormNameFileExpected, `correct part must have name="file" param`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpMultipartFilename = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeHttpMultipartFilename, `filename is empty`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpCurrentClient = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeHttpCurrentClient, `cannot obtain current client`).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlWrongSortFormat = func(ctx context.Context, v string) error {
		return er.WithBuilder(ErrCodeHttpUrlWrongSortFormat, "wrong sort format").Business().F(er.FF{"var": v}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrHttpUrlMaxPageSizeExceeded = func(ctx context.Context, maxPageSize int) error {
		return er.WithBuilder(ErrCodeHttpUrlMaxPageSizeExceeded, "max page size (%d) exceeded", maxPageSize).Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrRouteBuilderUrlEmpty = func() error {
		return er.WithBuilder(ErrCodeRouteBuilderUrlEmpty, "url empty").Err()
	}
	ErrRouteBuilderVerbEmpty = func(url string) error {
		return er.WithBuilder(ErrCodeRouteBuilderVerbEmpty, "verb empty").F(er.FF{"url": url}).Err()
	}
	ErrRouteBuilderBothHandleFuncAndHandlerEmpty = func(url string) error {
		return er.WithBuilder(ErrCodeRouteBuilderBothHandleFuncAndHandlerEmpty, "either handler or handle func must be specified").F(er.FF{"url": url}).Err()
	}
	ErrRouteBuilderAuthorizationPoliciesSpecifiedWithoutAuthentication = func(url string) error {
		return er.WithBuilder(ErrCodeRouteBuilderAuthorizationPoliciesSpecifiedWithoutAuthentication, "authorization requires authentication configured").F(er.FF{"url": url}).Err()
	}
	ErrRouteBuilderSpecialMiddlewaresRequireSubrouting = func(url string) error {
		return er.WithBuilder(ErrCodeRouteBuilderSpecialMiddlewaresRequireSubrouting, "special middlewares require subrouting").F(er.FF{"url": url}).Err()
	}
	ErrSecurityLoginFailed = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeUserLoginFail, "login failed").Business().C(ctx).HttpSt(http.StatusUnauthorized).Err()
	}
	ErrSecurityPermissionsDenied = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSecurityPermissionsDenied, "permissions denied").Business().C(ctx).HttpSt(http.StatusForbidden).Err()
	}
	ErrSecurityNoResourceRequested = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSecurityNoResourceRequested, "no resource requested").Business().C(ctx).HttpSt(http.StatusForbidden).Err()
	}
)
