package manager

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Configuration struct {
	IPFetchInterval time.Duration
	Records         []*Record
}

type Record struct {
	Connector string
	Domain    string
	SubDomain string
	Type      RecordType
	Interval  time.Duration
}

type RecordType string

const (
	TypeA    RecordType = "A"
	TypeAAAA RecordType = "AAAA"
)

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

	for _, record := range Config.Records {
		if !ConnectorIsRegistered(record.Connector) {
			return ErrConnectorNotFound
		}

		switch record.Type {
		case TypeA, TypeAAAA:
		default:
			log.Error().Err(ErrIncorrectType).Str("domain", record.Domain).Str("type", string(record.Type)).Msg("Invalid configuration")
			return ErrIncorrectType
		}
	}

	return nil
}
