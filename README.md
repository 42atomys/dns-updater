# DNS-Updater

Simple application in Golang that retrieves your ip and updates your DNS entries automatically each time your IP changes.

## Motivation

Having a lab at home. behind an ip which is dynamic, I could not keep a static DNS entry, since at each change of IP I have to change the DNS entries of all my domains.

So I wrote this little program one afternoon in order to get rid of this manual task which, when it occurred when I was not at home, was very annoying, imagine you can no longer watch plex while you are with your date. Unthinkable

## Usage

To use it, its really simple (3 steps)

You can use official Docker image : `atomys/dns-updater` (https://hub.docker.com/r/atomys/dns-updater) to test it or install it un your servers or Kubernetes Clusters. You can found an example to deploy this program on a Kubernetes cluster on the `examples/kubernetes` folder.

### Step 1 : Configure your updater.yaml
In config folder, edit the updater.yaml

```yaml
# Interval at which the program will check if your IP has changed
ipFetchInterval: 30s

# Records entries define the rules to follow when found an ip change
entries:
- #  name of the provider to use. Must be registered
  provider: ovh
  # domain to update
  domain: example.space
  # when your record is on a subdomain enter the name of
  # the sub domain here without domain extension
  # `awesome.example.space` must be `awesome`
  subDomain: null
  # Type of the record. Actually must be A or AAAA only
  type: A
```

### Step 2: Launch it
```
$ ./dns-updater
```

### Step 3: Configure provider authentication data through env

| Provider Name | Provider Site         | Environment Variables                                                                                                                                |
| ------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ovh`         | https://ovh.com       | `OVH_APPLICATION_KEY` Your OVH Application Key<br /> `OVH_APPLICATION_SECRET` Your OVH Apllication Secret<br /> `OVH_CONSUMER_KEY` Your Consumer Key |
| `gandi`       | https://www.gandi.net | `GANDI_APPLICATION_KEY` Your Gandi API key (https://account.gandi.net)                                                                               |
|               |                       |                                                                                                                                                      |

# Contribution

**This project is maintained on GitLab** : https://gitlab.com/Atomys/dns-updater

MR must be on GitLab to be accepted. All pull requests on GitHub will rewrite on Gitlab if necessary.

## How to add provider 

A provider must be have two function
```go
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
```

### Add a provider (built-in method)
Want to use this program but your DNS provider is not integrated? You can participate by creating your provider and a merge request

To add a new provider go to `pkg/providers` package and create a new file named by the name of DNS provider.

After that, Register the new Provider on `pkg/manager/providers.go` file on the `Registration Section`
and add your new built-in provider in the switch case
```go
	case "awesome":
		record.Provider = providers.NewAwesomeProvider()
```

### Add a provider (runtime method)
In some specific case, like custom dns server or custom dns service, private internal dns or somethings like that, 
you can't add this provider to the main repository. 

In this specific case, you can add a custom provider to your app and use the manager.

This is an example of custom dns-updater 
```go
package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"gitlab.com/atomys-universe/dns-updater/internal/pkg/dns"
	"gitlab.com/atomys-universe/dns-updater/pkg/manager"
)

type NinjaProvider struct{}

func (p NinjaProvider) Name() string {
	return "ninja"
}

func (p NinjaProvider) UpdateDNS(domainName, subDomain string, fieldType dns.RecordType, ip net.IP) error {
	// Do update stuff
	return nil
}

func main() {
	// Register your private Custom Provider before validate
	// the configuration
	manager.RegisterCustomProvider(NinjaProvider{})
	// Validate the configuration
	if err := manager.ValidateConfiguration(); err != nil {
		log.Fatal().Err(err).Msg("configuration is invalid")
	}
	// Start the manager
	go manager.Run()
	// Some log
	log.Info().Msg("DNS Updater is running")
	// Ctrl+C catch
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	defer func() {
		log.Info().Msg("Stoping DNS Updater...")
	}()

	<-c
}
```
You can now use `ninja` as provider in your configuration file ! 🎉

All contributions are welcome :)

## Thanks
Thanks to ipconfig.co to provide simple curl service to know our ip