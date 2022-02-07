package manager

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gitlab.com/atomys-universe/dns-updater/internal/pkg/dns"
)

type Configuration struct {
	IPFetchInterval      time.Duration
	ConfigurationEntries []*ConfigurationEntry `mapstructure:"entries"`
}

type ConfigurationEntry struct {
	ProviderName string   `mapstructure:"provider"`
	Provider     Provider `mapstructure:"-"`
	Domain       string
	SubDomain    string
	Type         dns.RecordType
}

var (
	Config               = &Configuration{}
	ErrIncorrectType     = errors.New("type defined is incorrect. Must be A or AAAA")
	ErrIncorrectInterval = errors.New("interval defined is incorrect. Must be longer or equal than 1 minute")
)

/**
 * Validate the configuration file and her content
 */
func ValidateConfiguration() error {
	err := viper.Unmarshal(&Config)
	if err != nil {
		return err
	}

	for _, record := range Config.ConfigurationEntries {
		if err := registerProvidersFromConfiguration(record); err != nil {
			return err
		}

		switch record.Type {
		case dns.TypeA, dns.TypeAAAA:
		default:
			log.Error().Err(ErrIncorrectType).Str("domain", record.Domain).Str("type", string(record.Type)).Msg("Invalid configuration")
			return ErrIncorrectType
		}
	}

	log.Debug().Msgf("Load %d configurations", len(Config.ConfigurationEntries))
	return nil
}
