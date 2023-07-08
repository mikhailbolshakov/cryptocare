package pg

import "github.com/mikhailbolshakov/cryptocare/src/kit/er"

const (
	ErrCodeGooseMigrationUp     = "DB-001"
	ErrCodeGooseMigrationGetVer = "DB-002"
	ErrCodePostgresOpen         = "DB-003"
	ErrCodeGooseFolderNotFound  = "DB-004"
	ErrCodeGooseFolderOpen      = "DB-005"
	ErrCodeGooseMigrationLock   = "DB-006"
	ErrCodeGooseMigrationUnLock = "DB-007"
)

var (
	ErrGooseMigrationUp     = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeGooseMigrationUp, "").Err() }
	ErrGooseMigrationGetVer = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeGooseMigrationGetVer, "").Err() }
	ErrPostgresOpen         = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodePostgresOpen, "").Err() }
	ErrGooseFolderNotFound  = func(path string) error {
		return er.WithBuilder(ErrCodeGooseFolderNotFound, "folder not found %s", path).Err()
	}
	ErrGooseFolderOpen    = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodeGooseFolderOpen, "folder open").Err() }
	ErrGooseMigrationLock = func(cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeGooseMigrationLock, "locking before migration").Err()
	}
	ErrGooseMigrationUnLock = func(cause error) error {
		return er.WrapWithBuilder(cause, ErrCodeGooseMigrationUnLock, "unlocking after migration").Err()
	}
)
