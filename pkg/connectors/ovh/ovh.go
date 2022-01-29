package ovh

import (
	"fmt"
	"net"
	"os"

	"github.com/ovh/go-ovh/ovh"
	"github.com/rs/zerolog/log"
	"gitlab.com/atomys-universe/dns-updater/pkg/manager"
)

var client *ovh.Client

type OvhConnector struct{}

func New() OvhConnector {
	return OvhConnector{}
}

func (OvhConnector) Name() string {
	return "ovh"
}

func (OvhConnector) Initialize() {
	if os.Getenv("OVH_APPLICATION_KEY") == "" || os.Getenv("OVH_APPLICATION_SECRET") == "" || os.Getenv("OVH_CONSUMER_KEY") == "" {
		log.Fatal().Msg("Missing env configuration for OVH connector")
	}

	client, _ = ovh.NewClient(
		"ovh-eu",
		os.Getenv("OVH_APPLICATION_KEY"),
		os.Getenv("OVH_APPLICATION_SECRET"),
		os.Getenv("OVH_CONSUMER_KEY"),
	)
}

type record struct {
	target    string
	zone      string
	id        int
	fieldType string
	ttl       int
	subDomain string
}

func (c OvhConnector) UpdateDNS(domainName, subDomain string, fieldType manager.RecordType, ip net.IP) error {
	var recordIds = []int{}
	client.Get(fmt.Sprintf("/domain/zone/%s/record?fieldType=%s", domainName, fieldType), &recordIds)

	for _, recordID := range recordIds {
		var rec = record{}
		client.Get(fmt.Sprintf("/domain/zone/%s/record/%d", domainName, recordID), &rec)

		rec.target = ip.String()
		if err := client.Put(fmt.Sprintf("/domain/zone/%s/record/%d", domainName, recordID), &rec, nil); err != nil {
			return err
		}
	}

	if err := client.Post(fmt.Sprintf("/domain/zone/%s/refresh", domainName), struct{ zoneName string }{domainName}, nil); err != nil {
		return err
	}

	return nil
}
