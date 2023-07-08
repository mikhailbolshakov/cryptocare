package errors

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/kit/er"
	"net/http"
)

var (
	ErrChainsCalculationAlreadyRun = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeChainsCalculationAlreadyRun, "already run").Business().C(ctx).Err()
	}
	ErrBidProviderAlreadyRun = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeBidProviderAlreadyRun, "already run").Business().C(ctx).Err()
	}
	ErrBidStorageScanBidsLight = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeBidStorageScanBidsLight, "").C(ctx).Err()
	}
	ErrBidStorageScanBidsLightReadRec = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeBidStorageScanBidsLightReadRec, "").C(ctx).Err()
	}
	ErrBidStoragePutBids = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeBidStoragePutBids, "").C(ctx).Err()
	}
	ErrBidStorageGetBidsByIds = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeBidStorageGetBidsByIds, "").C(ctx).Err()
	}
	ErrChainStoragePutChain = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeChainStoragePutChain, "").C(ctx).Err()
	}
	ErrChainStorageGetChain = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeChainStorageGetChain, "").C(ctx).Err()
	}
	ErrChainStorageScanChains = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeChainStorageScanChains, "").C(ctx).Err()
	}
	ErrUserPasswordHashGenerate = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserPasswordHashGenerate, "").C(ctx).Err()
	}
	ErrAuthPwdEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeAuthPwdEmpty, "password isn't specified").Business().C(ctx).Err()
	}
	ErrAuthPwdPolicy = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeAuthPwdPolicy, "password is too simple").Business().C(ctx).Err()
	}
	ErrUserEmailEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeUserEmailEmpty, "email empty").Business().C(ctx).Err()
	}
	ErrUserNoValidEmail = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeUserNoValidEmail, "email invalid").Business().C(ctx).Err()
	}
	ErrUserInvalidPassword = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeUserInvalidPassword, "password invalid").Business().C(ctx).Err()
	}
	ErrUserNameNotUnique = func(ctx context.Context, email string) error {
		return er.WithBuilder(ErrCodeUserNameNotUnique, "user already exists").Business().F(er.FF{"email": email}).C(ctx).Err()
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
	ErrLogoutNoSID = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeLogoutNoSID, "no SID specified").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrUserStorageCreate = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageCreate, "").C(ctx).Err()
	}
	ErrUserStorageUpdate = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageUpdate, "").C(ctx).Err()
	}
	ErrUserStorageClearCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageClearCache, "").C(ctx).Err()
	}
	ErrUserStorageAeroKey = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageAeroKey, "").C(ctx).Err()
	}
	ErrUserStorageGetCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageGetCache, "").C(ctx).Err()
	}
	ErrUserStoragePutCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStoragePutCache, "").C(ctx).Err()
	}
	ErrUserStorageGetDb = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageGetDb, "").C(ctx).Err()
	}
	ErrUserStorageCreateIndex = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageCreateIndex, "").C(ctx).Err()
	}
	ErrUserStorageGetCacheByUsername = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageGetCacheByUsername, "").C(ctx).Err()
	}
	ErrUserStorageGetByIds = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageGetByIds, "").C(ctx).Err()
	}
	ErrUserStorageDelete = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeUserStorageDelete, "").C(ctx).Err()
	}
	ErrSessionStorageAeroKey = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionStorageAeroKey, "").C(ctx).Err()
	}
	ErrSessionStorageGetCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionStorageGetCache, "").C(ctx).Err()
	}
	ErrSessionStoragePutCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionStoragePutCache, "").C(ctx).Err()
	}
	ErrSessionStorageClearCache = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionStorageClearCache, "").C(ctx).Err()
	}
	ErrSessionStorageGetDb = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionStorageGetDb, "").C(ctx).Err()
	}
	ErrSessionGetByUser = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionGetByUser, "").C(ctx).Err()
	}
	ErrSessionStorageUpdateLastActivity = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionStorageUpdateLastActivity, "").C(ctx).Err()
	}
	ErrSessionStorageUpdateLogout = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionStorageUpdateLogout, "").C(ctx).Err()
	}
	ErrSessionStorageCreateSession = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSessionStorageCreateSession, "").C(ctx).Err()
	}
	ErrSessionNotFound = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionNotFound, "session not found").C(ctx).Business().Err()
	}
	ErrSessionLoggedOut = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionLoggedOut, "session isn't active").C(ctx).Business().HttpSt(http.StatusForbidden).Err()
	}
	ErrSecurityPermissionsDenied = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSecurityPermissionsDenied, "permission denied").C(ctx).Business().HttpSt(http.StatusForbidden).Err()
	}
	ErrSessionAuthorizationInvalidResource = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSessionAuthorizationInvalidResource, "invalid resource").C(ctx).Business().HttpSt(http.StatusForbidden).Err()
	}
	ErrSidEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSidEmpty, "sid empty").C(ctx).Business().HttpSt(http.StatusForbidden).Err()
	}
	ErrNoAuthHeader = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeNoAuthHeader, "no authorization header provided").Business().C(ctx).HttpSt(http.StatusUnauthorized).Err()
	}
	ErrAuthHeaderInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeAuthHeaderInvalid, "authorization header invalid").Business().C(ctx).HttpSt(http.StatusUnauthorized).Err()
	}
	ErrNoUID = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeNoUID, "no UID in session context").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrSubscriptionMinProfitInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSubscriptionMinProfitInvalid, "min profit invalid").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrSubscriptionMaxDepthInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSubscriptionMaxDepthInvalid, "max depth invalid").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrSubscriptionNotificationChannelNotSupported = func(ctx context.Context, channelId string) error {
		return er.WithBuilder(ErrCodeSubscriptionNotificationChannelNotSupported, "notification channel not supported").Business().F(er.FF{"channelId": channelId}).C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrSubscriptionNotificationTelegramInvalid = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSubscriptionNotificationTelegramInvalid, "telegram params invalid").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrSubscriptionIdEmpty = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSubscriptionIdEmpty, "subscription id empty").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrSubscriptionNotFound = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSubscriptionNotFound, "subscription not found").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrSubscriptionNotActive = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeSubscriptionNotActive, "subscription not active").Business().C(ctx).HttpSt(http.StatusBadRequest).Err()
	}
	ErrSubscriptionStoragePut = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSubscriptionStoragePut, "").C(ctx).Err()
	}
	ErrSubscriptionStorageSearch = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSubscriptionStorageSearch, "").C(ctx).Err()
	}
	ErrSubscriptionStorageGet = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSubscriptionStorageGet, "").C(ctx).Err()
	}
	ErrSubscriptionStorageDel = func(cause error, ctx context.Context) error {
		return er.WrapWithBuilder(cause, ErrCodeSubscriptionStorageDel, "").C(ctx).Err()
	}
	ErrNotAllowed = func(ctx context.Context) error {
		return er.WithBuilder(ErrCodeNotAllowed, "operation isn't allowed").Business().C(ctx).HttpSt(http.StatusForbidden).Err()
	}
)
