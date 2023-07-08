package config

import (
	"github.com/mikhailbolshakov/cryptocare/src/kit/config/configuro"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"os"
	"path/filepath"
	"reflect"
)

// Loader loads config file to a custom struct
type Loader interface {
	// WithEnvPath specifies a path to .env file
	WithEnvPath(envPath string) Loader
	// WithConfigPath specifies a path to config file
	WithConfigPath(configPath string) Loader
	// WithConfigPathFromEnv specifies an env var from which path to config is taken
	WithConfigPathFromEnv(env string) Loader
	// LoadPath loads config based on given parameters and puts it to a target struct
	Load(target interface{}) error
}

type configLoaderImpl struct {
	logger            log.CLoggerFunc
	envPath           string
	configPath        string
	configPathEnvName string
}

func NewConfigLoader(logger log.CLoggerFunc) Loader {
	return &configLoaderImpl{
		logger: logger,
	}
}

func (s *configLoaderImpl) l() log.CLogger {
	return s.logger().Cmp("config-loader")
}

func (s *configLoaderImpl) WithEnvPath(envPath string) Loader {
	s.envPath = envPath
	return s
}

func (s *configLoaderImpl) WithConfigPath(configPath string) Loader {
	s.configPath = configPath
	return s
}

func (s *configLoaderImpl) WithConfigPathFromEnv(env string) Loader {
	s.configPathEnvName = env
	return s
}

func (s *configLoaderImpl) Load(target interface{}) error {
	l := s.l().Mth("load").F(log.FF{"cfg-path": s.configPath, "env-path": s.envPath, "env": s.configPathEnvName})

	if reflect.ValueOf(target).Kind() != reflect.Ptr || reflect.TypeOf(target).Elem().Kind() != reflect.Struct {
		return ErrConfigTargetObjectInvalidType()
	}

	var path string
	if s.configPath != "" {
		path = s.configPath
	} else if s.configPathEnvName != "" {
		path = os.Getenv(s.configPathEnvName)
	}

	if path == "" {
		return ErrConfigPathEmpty()
	}

	absPath, _ := filepath.Abs(path)
	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			return ErrConfigFileNotFound(absPath)
		}
		return ErrConfigFileOpen(err, absPath)
	}
	l.DbgF("config file loaded: %s", absPath)

	var absEnvPath string
	if s.envPath != "" {
		absEnvPath, _ = filepath.Abs(s.envPath)
		if _, err := os.Stat(absEnvPath); err != nil {
			if os.IsNotExist(err) {
				return ErrEnvFileNotFound(absPath)
			}
			return ErrEnvFileOpen(err, absPath)
		}
		l.DbgF("env file loaded: %s", absEnvPath)
	}

	// build options
	opts := []configuro.ConfigOptions{configuro.WithLoadFromConfigFile(absPath, true)}
	if absEnvPath != "" {
		opts = append(opts, configuro.WithLoadDotEnv(absEnvPath))
	}

	// create a new config loader
	Loader, err := configuro.NewConfig(opts...)
	if err != nil {
		return ErrConfigInit(err)
	}

	// load config
	err = Loader.Load(target)
	if err != nil {
		return ErrConfigLoad(err)
	}

	l.TrcObj("%v", target)

	return nil

}
