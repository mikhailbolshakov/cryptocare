package config

import "github.com/mikhailbolshakov/cryptocare/src/kit/er"

const (
	ErrCodeEnvRootPathNotSet             = "CFG-001"
	ErrCodeEnvFileOpening                = "CFG-002"
	ErrCodeEnvNewConfig                  = "CFG-003"
	ErrCodeEnvLoad                       = "CFG-004"
	ErrCodeServiceConfigNotSpecified     = "CFG-005"
	ErrCodeConfigPortEnvEmpty            = "CFG-006"
	ErrCodeConfigHostEnvEmpty            = "CFG-007"
	ErrCodeConfigTimeout                 = "CFG-008"
	ErrCodeEnvFileLoadingPath            = "CFG-009"
	ErrCodeConfigFileNotFound            = "CFG-010"
	ErrCodeConfigFileOpen                = "CFG-011"
	ErrCodeConfigPathEmpty               = "CFG-012"
	ErrCodeEnvFileNotFound               = "CFG-013"
	ErrCodeEnvFileOpen                   = "CFG-014"
	ErrCodeConfigInit                    = "CFG-015"
	ErrCodeConfigLoad                    = "CFG-016"
	ErrCodeConfigTargetObjectInvalidType = "CFG-017"
)

var (
	ErrEnvRootPathNotSet = func(v string) error {
		return er.WithBuilder(ErrCodeEnvRootPathNotSet, "root path env variable %s isn't set", v).Err()
	}
	ErrEnvFileOpening = func(cause error, path string) error {
		return er.WrapWithBuilder(cause, ErrCodeEnvFileOpening, "").F(er.FF{"path": path}).Err()
	}
	ErrEnvNewConfig       = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeEnvNewConfig, "").Err() }
	ErrEnvLoad            = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeEnvLoad, "").Err() }
	ErrEnvFileLoadingPath = func(cause error, file string) error {
		return er.WrapWithBuilder(cause, ErrCodeEnvFileLoadingPath, ".env loading").F(er.FF{"path": file}).Err()
	}
	ErrServiceConfigNotSpecified = func(svc string) error {
		return er.WithBuilder(ErrCodeServiceConfigNotSpecified, "config for service isn't specified").F(er.FF{"svc": svc}).Err()
	}
	ErrConfigPortEnvEmpty = func() error {
		return er.WithBuilder(ErrCodeConfigPortEnvEmpty, "env var CONFIG_CFG_GRPC_PORT is empty").Err()
	}
	ErrConfigHostEnvEmpty = func() error {
		return er.WithBuilder(ErrCodeConfigHostEnvEmpty, "env var CONFIG_CFG_GRPC_HOST is empty").Err()
	}
	ErrConfigTimeout      = func() error { return er.WithBuilder(ErrCodeConfigTimeout, "not ready within timeout").Err() }
	ErrConfigFileNotFound = func(path string) error {
		return er.WithBuilder(ErrCodeConfigFileNotFound, "config file %s not found", path).Err()
	}
	ErrEnvFileNotFound = func(path string) error {
		return er.WithBuilder(ErrCodeEnvFileNotFound, "env file %s not found", path).Err()
	}
	ErrConfigFileOpen = func(cause error, path string) error {
		return er.WrapWithBuilder(cause, ErrCodeConfigFileOpen, "open file %s", path).Err()
	}
	ErrEnvFileOpen = func(cause error, path string) error {
		return er.WrapWithBuilder(cause, ErrCodeEnvFileOpen, "open file %s", path).Err()
	}
	ErrConfigPathEmpty               = func() error { return er.WithBuilder(ErrCodeConfigPathEmpty, "config path empty").Err() }
	ErrConfigTargetObjectInvalidType = func() error {
		return er.WithBuilder(ErrCodeConfigTargetObjectInvalidType, "target object must be pointer on struct").Err()
	}
	ErrConfigInit = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeConfigInit, "").Err() }
	ErrConfigLoad = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeConfigLoad, "").Err() }
)
