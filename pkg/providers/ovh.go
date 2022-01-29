package providers

import (
	"fmt"
	"net"
	"os"

	"github.com/ovh/go-ovh/ovh"
	"github.com/rs/zerolog/log"
	"gitlab.com/atomys-universe/dns-updater/pkg/manager"
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
	target    string
	zone      string
	id        int
	fieldType string
	ttl       int
	subDomain string
}

func (c OVHProvider) UpdateDNS(domainName, subDomain string, fieldType manager.RecordType, ip net.IP) error {
	var recordIds = []int{}
	c.client.Get(fmt.Sprintf("/domain/zone/%s/record?fieldType=%s&subDomain=%s", domainName, fieldType, subDomain), &recordIds)

	for _, recordID := range recordIds {
		var rec = record{}
		c.client.Get(fmt.Sprintf("/domain/zone/%s/record/%d", domainName, recordID), &rec)

		rec.target = ip.String()
		if err := c.client.Put(fmt.Sprintf("/domain/zone/%s/record/%d", domainName, recordID), &rec, nil); err != nil {
			return err
		}
	}

	if err := c.client.Post(fmt.Sprintf("/domain/zone/%s/refresh", domainName), struct{ zoneName string }{domainName}, nil); err != nil {
		return err
	}

	return nil
}
