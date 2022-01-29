package manager

import (
	"errors"
	"net"
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
	UpdateDNS(domainName, subDomain string, fieldType RecordType, ip net.IP) error
}

var (
	providers []Provider
	// Error : No provider is found in the manager
	ErrProviderNotFound = errors.New("provider not found")
)

/**
 *  RegisterProvider a new provider into the manager
 */
func RegisterProvider(provider Provider) {
	if ProviderIsRegistered(provider.Name()) {
		return
	}

	providers = append(providers, provider)
}

/**
 *  Retrieve if a provider is registered in the manager or not
 */
func ProviderIsRegistered(providerName string) bool {
	for _, provider := range providers {
		if provider.Name() == providerName {
			return true
		}
	}
	return false
}

/**
 * Retrieve a provider by her name
 */
func Get(providerName string) (Provider, error) {
	for _, provider := range providers {
		if provider.Name() == providerName {
			return provider, nil
		}
	}

	return nil, ErrProviderNotFound
}
