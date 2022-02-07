package manager

import (
	"errors"
	"net"

	"gitlab.com/atomys-universe/dns-updater/internal/pkg/dns"
	"gitlab.com/atomys-universe/dns-updater/pkg/providers"
)

/**
 * Definition of a Provider interface
 * Used to implement any DNS service to make it compatible for DNS Updater
 *
 * Name()       =>  Return the name of the provider. Must be unique
 * UpdateDNS()  =>  Method called when an IP changes is detected and used to update
 *                  DNS Entry on Provider
 */
type Provider interface {
	// Return the name of the provider. Must be unique
	Name() string
	// Method called when an IP changes is detected and used to update
	// DNS Entry on Provider
	UpdateDNS(domainName, subDomain string, fieldType dns.RecordType, ip net.IP) error
}

var (
	customProviders []Provider
	// Error : No provider is found in the manager
	ErrProviderNotFound = errors.New("provider not found")
)

// ##########################
// ## Registration Section ##
// ##########################

/**
 * Register providers which is entered in the configuration
 * and initialize them
 */
func registerProvidersFromConfiguration(record *ConfigurationEntry) error {
	switch record.ProviderName {
	case "ovh":
		record.Provider = providers.NewOvhProvider()
	case "gandi":
		record.Provider = providers.NewGandiProvider()
	default:
		var err error
		record.Provider, err = getCustomProvider(record.ProviderName)
		if err != nil {
			return err
		}
	}

	return nil
}

// ##########################
// ## Registration Section ##
// ##########################

/**
 *  RegisterProvider a custom provider
 */
func RegisterCustomProvider(provider Provider) {
	if customProviderExist(provider.Name()) {
		return
	}

	customProviders = append(customProviders, provider)
}

/**
 *  Retrieve if a custom provider is registered or not
 */
func customProviderExist(providerName string) bool {
	for _, provider := range customProviders {
		if provider.Name() == providerName {
			return true
		}
	}
	return false
}

/**
 * getCustomProvider will search and found if a custom provider with
 * given name is registered or not.
 */
func getCustomProvider(providerName string) (Provider, error) {
	for _, provider := range customProviders {
		if provider.Name() == providerName {
			return provider, nil
		}
	}
	return nil, ErrProviderNotFound
}
