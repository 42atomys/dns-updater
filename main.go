package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gitlab.com/atomys-universe/dns-updater/pkg/connectors/ovh"
	"gitlab.com/atomys-universe/dns-updater/pkg/manager"
)

const WebhookURL = "https://canary.discord.com/api/webhooks/858713820200566804/NbsedN-G2yzbtM2vM9TyKXODYe4Jw0HVtC_AcZxPk9yTsqA5LhBsAxsBo23SYFJ0hKmK"

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

	manager.RegisterConnector(ovh.New())

	if err := manager.ValidateConfiguration(); err != nil {
		log.Fatal().Err(err).Msg("configuration is invalid")
	}

	go manager.Run()

	log.Info().Msg("DNS Updater is running")
	c := make(chan os.Signal, 2)
	<-c
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// c := Content{
	// 	Content:  fmt.Sprintf("IP: %s", string(bodyBytes)),
	// 	Username: "DNS Updater",
	// }

	// var jsonData []byte
	// jsonData, err = json.Marshal(c)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Cannot Marshall")
	// }

	// log.Print("Post to Discord")
	// _, err = http.Post(WebhookURL, "application/json", bytes.NewReader(jsonData))
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Cannot post")
	// }

	// log.Print("Done !")
}

//  https://discord.com/api/webhooks/858713820200566804/NbsedN-G2yzbtM2vM9TyKXODYe4Jw0HVtC_AcZxPk9yTsqA5LhBsAxsBo23SYFJ0hKmK
