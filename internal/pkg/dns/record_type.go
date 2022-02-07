package dns

// RecordType:
//DNS Record Types allowed to be updated by the dns-updater
type RecordType string

const (
	// Type A Record (address)‍
	// Most commonly used to map a fully qualified domain name (FQDN) to an IPv4
	// address and acts as a translator by converting domain names to IP addresses.
	TypeA RecordType = "A"
	// Type AAAA Record (quad A)
	// Similar to A Records but maps to an IPv6 address
	// (smartphones prefer IPv6, if available).
	TypeAAAA RecordType = "AAAA"
)
