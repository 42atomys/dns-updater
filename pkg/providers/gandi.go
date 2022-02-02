package providers

import (
	"net"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/go-gandi/go-gandi/config"
	"github.com/go-gandi/go-gandi/livedns"
	"gitlab.com/atomys-universe/dns-updater/pkg/manager"
)

type GandiProvider struct {
	client *livedns.LiveDNS
}

func NewGandiProvider() GandiProvider {
	if os.Getenv("GANDI_APPLICATION_KEY") == "" {
		log.Fatal().Msg("Missing env configuration for Gandi provider.")
	}

	newConfig := config.Config{
		APIKey:    os.Getenv("GANDI_APPLICATION_KEY"),
		SharingID: "",
		Debug:     false,
		DryRun:    false,
		APIURL:    config.APIURL,
	}

	client := livedns.New(newConfig)

	return GandiProvider{
		client: client,
	}
}

func (c GandiProvider) Name() string {
	return "gandi"
}

func (c GandiProvider) UpdateDNS(domainName, subDomain string, fieldType manager.RecordType, ip net.IP) error {
	newRecord := livedns.DomainRecord{
		RrsetType:   string(fieldType),
		RrsetName:   subDomain,
		RrsetValues: []string{ip.String()},
	}

	_, err := c.client.UpdateDomainRecordsByName(domainName, subDomain, []livedns.DomainRecord{newRecord})
	if err != nil {
		return err
	}
	return nil
}
