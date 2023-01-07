package ip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchIpv4(t *testing.T) {
	ip := fetchIPv4()

	assert.NotNil(t, ip)
}

func TestFetchIpv6(t *testing.T) {
	ip := fetchIPv6()

	if ip == nil {
		t.Skip("No IPv6")
	}

	assert.NotNil(t, ip)
	assert.True(t, ip.To4() == nil)
}
