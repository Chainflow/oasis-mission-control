package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type (
	//Telegram bot details struct
	Telegram struct {
		BotToken string `mapstructure:"tg_bot_token"`
		ChatID   int64  `mapstructure:"tg_chat_id"`
	}

	//SendGrid tokens
	SendGrid struct {
		Token          string `mapstructure:"sendgrid_token"`
		EmailAddress   string `mapstructure:"email_address"`
		PagerdutyEmail string `mapstructure:"pagerduty_email"`
	}

	//Scraper time interval
	Scraper struct {
		Rate          string `mapstructure:"rate"`
		Port          string `mapstructure:"port"`
		ValidatorRate string `mapstructure:"validator_rate"`
	}

	//InfluxDB details
	InfluxDB struct {
		Port     string `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	// DailyAlerts struct to send validator status
	DailyAlerts struct {
		AlertTime1 string `mapstructure:"alert_time1"`
		AlertTime2 string `mapstructure:"alert_time2"`
	}

	// EnableAlerts struct which holds option to enable alerts
	EnableAlerts struct {
		EnableTelegramAlerts string `mapstructure:"enable_telegram_alerts"`
		EnableEmailAlerts    string `mapstructure:"enable_email_alerts"`
	}

	// ValidatorDetails struct
	ValidatorDetails struct {
		ValidatorAddress    string `mapstructure:"validator_addr"`
		ValidatorHexAddress string `mapstructure:"validator_hex_addr"`
		ValidatorName       string `mapstructure:"validator_name"`
		NetworkURL          string `mapstructure:"network_url"`
		NetworkNodeName     string `mapstructure:"network_node_name"`
		SocketPath          string `mapstructure:"socket_path"`
	}

	// 
	AlertsThreshold struct {
		VotingPowerThreshold           int64 `mapstructure:"voting_power_threshold"`
		NumPeersThreshold              int64 `mapstructure:"num_peers_threshold"`
		BlockDiffThreshold             int64 `mapstructure:"block_diff_threshold"`
		MissedBlocksThreshold          int64 `mapstructure:"missed_blocks_threshold"`
		EpochDiffThreshold             int64 `mapstructure:"epoch_diff_threshold"`
		EmergencyMissedBlocksThreshold int64 `mapstructure:"emergency_missed_blocks_threshold"`
	}

	//Config
	Config struct {
		Scraper          Scraper          `mapstructure:"scraper"`
		Telegram         Telegram         `mapstructure:"telegram"`
		SendGrid         SendGrid         `mapstructure:"sendgrid"`
		InfluxDB         InfluxDB         `mapstructure:"influxdb"`
		DailyAlerts      DailyAlerts      `mapstructure:"daily_alerts"`
		EnableAlerts     EnableAlerts     `mapstructure:"enable_alerts"`
		ValidatorDetails ValidatorDetails `mapstructure:"validator_details"`
		AlertsThreshold  AlertsThreshold  `mapstructure:"alerts_threshold"`
	}
)

//ReadFromFile to read config details using viper
func ReadFromFile() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./config/")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config.toml: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config.toml to application config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("error occurred in config validation: %v", err)
	}

	return &cfg, nil
}

//Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
