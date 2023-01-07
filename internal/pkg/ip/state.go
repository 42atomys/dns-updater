package ip

import (
	"time"
)

type IPChangeState struct {
	IPV4Change bool
	IPV6Change bool
}

/**
 * This function will fetch the difference of ip and send change in the given
 * chan
 */
func FetchIPChangeRoutine(c chan IPChangeState, interval time.Duration) {
	for {
		var state = IPChangeState{}

		var oldIpv4 = CurrentIPv4
		var fetchedIPv4 = fetchIPv4()
		if fetchedIPv4 != nil && !oldIpv4.Equal(*fetchedIPv4) {
			CurrentIPv4 = *fetchedIPv4
			state.IPV4Change = true
		}

		if ipv6Available() {
			var oldIpv6 = CurrentIPv6
			var fetchedIPv6 = fetchIPv6()
			if fetchedIPv6 != nil && !oldIpv6.Equal(*fetchedIPv6) {
				CurrentIPv6 = *fetchedIPv6
				state.IPV6Change = true
			}
		}

		c <- state
		time.Sleep(interval)
	}
}
