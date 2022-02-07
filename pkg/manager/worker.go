package manager

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/atomys-universe/dns-updater/internal/pkg/dns"
	"gitlab.com/atomys-universe/dns-updater/internal/pkg/ip"
)

// Worker used to pass informations from configuration
// to the worker pool
// Actually only contains the ConfigurationEntry and the provider
// associated with
type Worker struct {
	ConfigurationEntry *ConfigurationEntry
	Provider           Provider
}

// Store workers after spawnWorkers
var workers = []*Worker{}

/**
 * Initialize the workers pools witk all entries and providers
 */
func initializeWorkers() {
	log.Info().Msgf("Initializing workers for %d entries", len(Config.ConfigurationEntries))
	for _, entry := range Config.ConfigurationEntries {

		var worker = &Worker{
			ConfigurationEntry: entry,
			Provider:           entry.Provider,
		}
		workers = append(workers, worker)
	}
}

/**
 * Run the manger and start to fetch IP changes 🎉
 */
func Run() {
	var ipChanges = make(chan ip.IPChangeState)
	go ip.FetchIPChangeRoutine(ipChanges, Config.IPFetchInterval)

	initializeWorkers()

	for {
		receiveIpChanges(<-ipChanges)
	}
}

/**
 * Function called each time we receive a new IP change state. This function
 * makes the decision to update or not the DNS if the currentIP change.
 */
func receiveIpChanges(state ip.IPChangeState) {
	log.Debug().Msgf("New Ip State arrived: %+v", state)
	for _, w := range workers {
		var err error
		switch w.ConfigurationEntry.Type {
		case dns.TypeA:
			if state.IPV4Change {
				log.Info().Str("domain", w.ConfigurationEntry.Domain).Str("type", string(w.ConfigurationEntry.Type)).Str("provider", w.Provider.Name()).Msg("Detect a new IPv4. Update DNS entry")
				err = w.Provider.UpdateDNS(w.ConfigurationEntry.Domain, w.ConfigurationEntry.SubDomain, dns.TypeA, ip.CurrentIPv4)
			}
		case dns.TypeAAAA:
			if state.IPV6Change {
				log.Info().Str("domain", w.ConfigurationEntry.Domain).Str("type", string(w.ConfigurationEntry.Type)).Str("provider", w.Provider.Name()).Msg("Detect a new IPv6. Update DNS entry")
				err = w.Provider.UpdateDNS(w.ConfigurationEntry.Domain, w.ConfigurationEntry.SubDomain, dns.TypeAAAA, ip.CurrentIPv6)
			}
		}
		if err != nil {
			log.Error().Err(err).Msg("Cannot update DNS")
			continue
		}
		log.Info().Str("domain", w.ConfigurationEntry.Domain).Str("type", string(w.ConfigurationEntry.Type)).Str("provider", w.Provider.Name()).Msg("DNS Entry updated successfully")
	}
}
