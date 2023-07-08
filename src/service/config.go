package service

import (
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	kitConfig "github.com/mikhailbolshakov/cryptocare/src/kit/config"
	kitHttp "github.com/mikhailbolshakov/cryptocare/src/kit/http"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	kitAero "github.com/mikhailbolshakov/cryptocare/src/kit/storages/aerospike"
	"github.com/mikhailbolshakov/cryptocare/src/kit/storages/pg"
	"os"
	"path/filepath"
)

// Here are defined all types for your configuration
// You can remove not needed types or add your own

type Storages struct {
	Aero *kitAero.Config
	Pg   *pg.DbClusterConfig
}

type Api struct {
	Username string
	Password string
	Rest     *Address
}

type Address struct {
	Host string
	Port string
}

type ArbitrageNotificationTelegram struct {
	Bot string
}

type ArbitrageNotification struct {
	Telegram *ArbitrageNotificationTelegram
}

type Arbitrage struct {
	Assets                 string
	Depth                  int
	ProcessAssetsPeriodSec int     `config:"process-assets-period-sec"`
	BidProviderPeriodSec   int     `config:"bid-provider-period-sec"`
	MinProfit              float64 `config:"min-profit"`
	CheckLimit             bool    `config:"check-limit"`
	Notification           *ArbitrageNotification
}

type Dev struct {
	Enabled               bool
	BidGeneratorPeriodSec int `config:"bid-gen-period-sec"`
	BidGeneratorBidsCount int `config:"bid-gen-bids-count"`
}

type Config struct {
	Log       *log.Config
	Http      *kitHttp.Config
	Api       *Api
	Storages  *Storages
	Auth      *auth.Config
	Dev       *Dev
	Arbitrage *Arbitrage
}

func LoadConfig() (*Config, error) {

	// get root folder from env
	rootPath := os.Getenv("CRYPTOCAREROOT")
	if rootPath == "" {
		return nil, kitConfig.ErrEnvRootPathNotSet("CRYPTOCAREROOT")
	}

	// config path
	configPath := filepath.Join(rootPath, "trading", "config.yml")

	// .env path
	envPath := filepath.Join(rootPath, "trading", ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		envPath = ""
	}

	// load config
	config := &Config{}
	err := kitConfig.NewConfigLoader(LF()).
		WithConfigPath(configPath).
		WithEnvPath(envPath).
		Load(config)

	if err != nil {
		return nil, err
	}
	return config, nil
}
