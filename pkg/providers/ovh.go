package providers

import (
	"fmt"
	"net"
	"os"

	"github.com/ovh/go-ovh/ovh"
	"github.com/rs/zerolog/log"
	"gitlab.com/atomys-universe/dns-updater/internal/pkg/dns"
)

type OVHProvider struct {
	client *ovh.Client
}

func NewOvhProvider() OVHProvider {
	if os.Getenv("OVH_APPLICATION_KEY") == "" || os.Getenv("OVH_APPLICATION_SECRET") == "" || os.Getenv("OVH_CONSUMER_KEY") == "" {
		log.Fatal().Msg("Missing env configuration for OVH provider")
	}

	client, _ := ovh.NewClient(
		"ovh-eu",
		os.Getenv("OVH_APPLICATION_KEY"),
		os.Getenv("OVH_APPLICATION_SECRET"),
		os.Getenv("OVH_CONSUMER_KEY"),
	)

	return OVHProvider{
		client: client,
	}
}

func (OVHProvider) Name() string {
	return "ovh"
}

/**
 * record structure used to update DNS entry on OVH API
 */
type record struct {
	Target    string `json:"target"`
	Zone      string `json:"zone"`
	ID        int    `json:"id"`
	FieldType string `json:"fieldType"`
	Ttl       int    `json:"ttl"`
	SubDomain string `json:"subDomain"`
}

func (c OVHProvider) UpdateDNS(domainName, subDomain string, fieldType dns.RecordType, ip net.IP) error {
	var recordIds = []int{}
	if err := c.client.Get(fmt.Sprintf("/domain/zone/%s/record?fieldType=%s&subDomain=%s", domainName, fieldType, subDomain), &recordIds); err != nil {
		return err
	}

	for _, recordID := range recordIds {
		var rec = record{}
		if err := c.client.Get(fmt.Sprintf("/domain/zone/%s/record/%d", domainName, recordID), &rec); err != nil {
			return err
		}

		rec.Target = ip.String()
		if err := c.client.Put(fmt.Sprintf("/domain/zone/%s/record/%d", domainName, recordID), &rec, nil); err != nil {
			return err
		}
	}

	if err := c.client.Post(fmt.Sprintf("/domain/zone/%s/refresh", domainName), struct{ zoneName string }{domainName}, nil); err != nil {
		return err
	}

	return nil
}
