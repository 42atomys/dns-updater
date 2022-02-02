package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gitlab.com/atomys-universe/dns-updater/pkg/manager"
	"gitlab.com/atomys-universe/dns-updater/pkg/providers"
)

type Content struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	viper.SetConfigName("updater")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Fatal error on reading config file")
	}
}

func main() {
	log.Info().Msg("DNS Updater starting...")

	manager.RegisterProvider(providers.NewOvhProvider())
	manager.RegisterProvider(providers.NewGandiProvider())

	if err := manager.ValidateConfiguration(); err != nil {
		log.Fatal().Err(err).Msg("configuration is invalid")
	}

	go manager.Run()

	log.Info().Msg("DNS Updater is running")
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	defer func() {
		log.Info().Msg("Stoping DNS Updater...")
	}()

	<-c
}
