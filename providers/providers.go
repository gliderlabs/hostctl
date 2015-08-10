package providers

import "fmt"

var providers = make(map[string]HostProvider)

func Register(provider HostProvider, name string) {
	providers[name] = provider
}

func Get(name string) (HostProvider, error) {
	p, found := providers[name]
	if !found {
		return nil, fmt.Errorf("Provider not registered: %s", name)
	}
	return p, p.Setup()
}

type HostProvider interface {
	Setup() error
	Create(host Host) error
	Destroy(name string) error
	List(pattern string) []Host
	Get(name string) *Host
}

type Host struct {
	Name    string
	IP      string
	Region  string
	Image   string
	Keyname string
	Flavor  string
}
