package ip

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

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
func fetchIPv4() *net.IP {
	return fetch("https://ifconfig.co/ip", "tcp4")
}

/**
 * Fetch the IPv6 of current machine with https://ipconfig.co service
 * Thanks to him
 */
func fetchIPv6() *net.IP {
	return fetch("https://ifconfig.co/ip", "tcp6")
}

/**
 * Fetch ip of current machine with given url
 * url needs to return ip on text/plain Content-Type
 */
func fetch(url string, mode string) *net.IP {
	log.Debug().Msg("Getting current IP")
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	if mode == "tcp6" && !ipv6Available() {
		log.Debug().Msg("No IPv6 available")
		return nil
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, mode, addr)
	}
	httpClient.Transport = transport

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch current IP")
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch current IP")
		return nil
	}

	s := strings.Trim(string(bodyBytes), " \n")
	ip := net.ParseIP(s)

	log.Debug().Msgf("Current IP %s: %s", mode, ip.String())
	return &ip
}

func ipv6Available() bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() == nil {
				return true
			}
		}
	}

	return false
}
