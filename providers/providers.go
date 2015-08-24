package providers

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/MattAitchison/env"
)

var providers = make(map[string]HostProvider)

func Register(provider HostProvider, name string) {
	providers[name] = provider
}

func Get(name string, setup bool) (HostProvider, error) {
	p, found := providers[name]
	if !found {
		return nil, fmt.Errorf("Provider not registered: %s", name)
	}
	if setup {
		return p, p.Setup()
	} else {
		return p, nil
	}
}

type HostProvider interface {
	Setup() error
	Create(host Host) error
	Destroy(name string) error
	List(pattern string) []Host
	Get(name string) *Host
	Env() *env.EnvSet
}

type Host struct {
	Name     string
	IP       string
	Region   string
	Image    string
	Keyname  string
	Flavor   string
	Userdata string
}

type TestProvider struct {
	mu    sync.Mutex
	Hosts []Host
}

func (p *TestProvider) Setup() error {
	return nil
}

func (p *TestProvider) Create(host Host) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Hosts = append(p.Hosts, host)
	return nil
}

func (p *TestProvider) Destroy(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	var hosts []Host
	for i := range p.Hosts {
		if p.Hosts[i].Name != name {
			hosts = append(hosts, p.Hosts[i])
		}
	}
	p.Hosts = hosts
	return nil
}

func (p *TestProvider) List(pattern string) []Host {
	p.mu.Lock()
	defer p.mu.Unlock()
	var hosts []Host
	for i := range p.Hosts {
		if ok, _ := filepath.Match(pattern, p.Hosts[i].Name); ok {
			hosts = append(hosts, p.Hosts[i])
		}
	}
	return hosts
}

func (p *TestProvider) Get(name string) *Host {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.Hosts {
		if p.Hosts[i].Name == name {
			return &p.Hosts[i]
		}
	}
	return nil
}

func (p *TestProvider) Env() *env.EnvSet {
	return nil
}
