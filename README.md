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
records:
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

| Provider Name | Provider Site   | Environment Variables                                                                                                                                |
| ------------- | --------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ovh`         | https://ovh.com | `OVH_APPLICATION_KEY` Your OVH Application Key<br /> `OVH_APPLICATION_SECRET` Your OVH Apllication Secret<br /> `OVH_CONSUMER_KEY` Your Consumer Key |
|               |                 |                                                                                                                                                      |

# Contribution

**This project is maintained on GitLab** : https://gitlab.com/Atomys/dns-updater

MR must be on GitLab to be accepted. All pull requests on GitHub will rewrite on Gitlab if necessary.

## How to add provider
Want to use this program but your DNS provider is not integrated? You can participate by creating your provider and a merge request

To add a new provider go to `pkg/providers` package and create a new file named by the name of DNS provider.

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
	UpdateDNS(domainName, subDomain string, fieldType RecordType, ip net.IP) error
}
```

After that, Register the new Provider on `main.go` file on the `Registration Section`
```go
manager.RegisterProvider(providers.NewAwesomeProvider())
```

All contributions are welcome :)

## Thanks
Thanks to ipconfig.co to provide simple curl service to know our ip