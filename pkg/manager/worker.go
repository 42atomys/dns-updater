package manager

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/atomys-universe/dns-updater/internal/pkg/ip"
)

type Worker struct {
	Record   *Record
	Provider Provider
}

// Store workers after spawnWorkers
var workers = []*Worker{}

/**
 * Initialize the workers pools witk all records and providers
 */
func initializeWorkers() {
	log.Info().Msgf("Initializing workers for %d records", len(Config.Records))
	for _, record := range Config.Records {
		conn, err := Get(record.Provider)
		if err != nil {
			panic("can't finish here. Provider is validate before")
		}

		var worker = &Worker{
			Record:   record,
			Provider: conn,
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
		switch w.Record.Type {
		case TypeA:
			if state.IPV4Change {
				log.Info().Str("domain", w.Record.Domain).Str("type", string(w.Record.Type)).Str("provider", w.Provider.Name()).Msg("Detect a new IPv4. Update DNS entry")
				err = w.Provider.UpdateDNS(w.Record.Domain, w.Record.SubDomain, TypeA, ip.CurrentIPv4)
			}
		case TypeAAAA:
			if state.IPV6Change {
				log.Info().Str("domain", w.Record.Domain).Str("type", string(w.Record.Type)).Str("provider", w.Provider.Name()).Msg("Detect a new IPv6. Update DNS entry")
				err = w.Provider.UpdateDNS(w.Record.Domain, w.Record.SubDomain, TypeAAAA, ip.CurrentIPv6)
			}
		}
		if err != nil {
			log.Error().Err(err).Msg("Cannot update DNS")
			continue
		}
		log.Info().Str("domain", w.Record.Domain).Str("type", string(w.Record.Type)).Str("provider", w.Provider.Name()).Msg("DNS Entry updated successfully")
	}
}
