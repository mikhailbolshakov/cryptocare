package errors

const (
	ErrCodeChainsCalculationAlreadyRun                 = "TRD-001"
	ErrCodeBidProviderAlreadyRun                       = "TRD-002"
	ErrCodeBidStorageScanBidsLight                     = "TRD-003"
	ErrCodeBidStorageScanBidsLightReadRec              = "TRD-004"
	ErrCodeBidStoragePutBids                           = "TRD-005"
	ErrCodeBidStorageGetBidsByIds                      = "TRD-006"
	ErrCodeChainStoragePutChain                        = "TRD-007"
	ErrCodeChainStorageScanChains                      = "TRD-008"
	ErrCodeAuthPwdEmpty                                = "TRD-009"
	ErrCodeAuthPwdPolicy                               = "TRD-010"
	ErrCodeUserPasswordHashGenerate                    = "TRD-011"
	ErrCodeUserEmailEmpty                              = "TRD-012"
	ErrCodeUserNoValidEmail                            = "TRD-013"
	ErrCodeUserNameNotUnique                           = "TRD-014"
	ErrCodeUserNotFound                                = "TRD-015"
	ErrCodeUserNotActive                               = "TRD-016"
	ErrCodeUserLocked                                  = "TRD-017"
	ErrCodeLogoutNoSID                                 = "TRD-018"
	ErrCodeUserStorageCreate                           = "TRD-019"
	ErrCodeUserStorageClearCache                       = "TRD-020"
	ErrCodeUserStorageUpdate                           = "TRD-021"
	ErrCodeUserStorageAeroKey                          = "TRD-022"
	ErrCodeUserStorageGetCache                         = "TRD-023"
	ErrCodeUserStoragePutCache                         = "TRD-024"
	ErrCodeUserStorageGetDb                            = "TRD-025"
	ErrCodeUserStorageCreateIndex                      = "TRD-026"
	ErrCodeUserStorageGetCacheByUsername               = "TRD-027"
	ErrCodeUserStorageGetByIds                         = "TRD-028"
	ErrCodeUserStorageDelete                           = "TRD-029"
	ErrCodeSessionStorageAeroKey                       = "TRD-030"
	ErrCodeSessionStorageGetCache                      = "TRD-031"
	ErrCodeSessionStoragePutCache                      = "TRD-032"
	ErrCodeSessionStorageClearCache                    = "TRD-033"
	ErrCodeSessionStorageGetDb                         = "TRD-034"
	ErrCodeSessionGetByUser                            = "TRD-035"
	ErrCodeSessionStorageUpdateLastActivity            = "TRD-036"
	ErrCodeSessionStorageUpdateLogout                  = "TRD-037"
	ErrCodeSessionStorageCreateSession                 = "TRD-038"
	ErrCodeSessionNotFound                             = "TRD-039"
	ErrCodeSessionLoggedOut                            = "TRD-040"
	ErrCodeSecurityPermissionsDenied                   = "TRD-041"
	ErrCodeSessionAuthorizationInvalidResource         = "TRD-042"
	ErrCodeSidEmpty                                    = "TRD-043"
	ErrCodeNoAuthHeader                                = "TRD-044"
	ErrCodeAuthHeaderInvalid                           = "TRD-045"
	ErrCodeUserInvalidPassword                         = "TRD-046"
	ErrCodeNoUID                                       = "TRD-047"
	ErrCodeChainStorageGetChain                        = "TRD-048"
	ErrCodeSubscriptionMinProfitInvalid                = "TRD-049"
	ErrCodeSubscriptionMaxDepthInvalid                 = "TRD-050"
	ErrCodeSubscriptionNotificationChannelNotSupported = "TRD-051"
	ErrCodeSubscriptionNotificationTelegramInvalid     = "TRD-052"
	ErrCodeSubscriptionIdEmpty                         = "TRD-053"
	ErrCodeSubscriptionNotFound                        = "TRD-054"
	ErrCodeSubscriptionNotActive                       = "TRD-055"
	ErrCodeSubscriptionStoragePut                      = "TRD-056"
	ErrCodeSubscriptionStorageSearch                   = "TRD-057"
	ErrCodeSubscriptionStorageGet                      = "TRD-058"
	ErrCodeSubscriptionStorageDel                      = "TRD-059"
	ErrCodeNotAllowed                                  = "TRD-060"
)
