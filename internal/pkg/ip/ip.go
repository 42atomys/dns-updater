package ip

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

var (
	// CurrentIPv4 of the machine
	CurrentIPv4 net.IP = net.ParseIP("0.0.0.0")
	// CurrentIPv6 of the machine - only if machine have ipv6
	CurrentIPv6 net.IP = net.ParseIP("0:0:0:0:0:0:0:0")
)

/**
 * Fetch the IPv4 of current machine with https://ipconfig.co service
 * Thanks to him
 */
func fetchIPv4() net.IP {
	return fetch("https://v4.ifconfig.co/ip")
}

/**
 * Fetch the IPv6 of current machine with https://ipconfig.co service
 * Thanks to him
 */
func fetchIPv6() net.IP {
	return fetch("https://v6.ifconfig.co/ip")
}

/**
 * Fetch ip of current machine with given url
 * url needs to return ip on text/plain Content-Type
 */
func fetch(url string) net.IP {
	log.Debug().Msg("Getting current IP")
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch current IP")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch current IP")
	}

	return net.ParseIP(string(bodyBytes))
}
