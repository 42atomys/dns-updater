package manager

import (
	"errors"
	"net"
)

/**
 * Definition of a Connector interface
 * Used to implement any DNS service to make it compatible for DNS Updater
 *
 * Initialize() =>  Called when program start to initialize any client configuration
 *                  If your connector dont need to be initialize, let this func empty
 * Name()       =>  Return the name of the connector. Must be unique
 */
type Connector interface {
	Initialize()
	Name() string
	UpdateDNS(domainName, subDomain string, fieldType RecordType, ip net.IP) error
}

var (
	connectors []Connector
	// Error : No connector is found in the manager
	ErrConnectorNotFound = errors.New("connector not found")
)

/**
 *  RegisterConnector a new connector into the manager
 */
func RegisterConnector(connector Connector) {
	if ConnectorIsRegistered(connector.Name()) {
		return
	}

	connectors = append(connectors, connector)
}

/**
 *  Retrieve if a connector is registered in the manager or not
 */
func ConnectorIsRegistered(connectorName string) bool {
	for _, connector := range connectors {
		if connector.Name() == connectorName {
			return true
		}
	}
	return false
}

/**
 * Retrieve a connector by her name
 */
func Get(connectorName string) (Connector, error) {
	for _, connector := range connectors {
		if connector.Name() == connectorName {
			return connector, nil
		}
	}

	return nil, ErrConnectorNotFound
}
