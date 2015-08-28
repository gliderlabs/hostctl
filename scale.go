package main

import (
	"bytes"
	"os"
	"strconv"
	"sync"
	"text/template"

	"github.com/gliderlabs/hostctl/providers"
	"github.com/spf13/cobra"
)

type HostNamer interface {
	HostNamePattern() string
}

var defaultHostPattern = "{{.Namespace}}{{.Name}}.{{.Index}}"

func init() {
	Hostctl.AddCommand(scaleCmd)
}

var scaleCmd = &cobra.Command{
	Use:   "scale <name> <count>",
	Short: "Resize host cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if (len(args) < 2 && defaultName == "") ||
			(len(args) < 1 && defaultName != "") {
			cmd.Usage()
			os.Exit(1)
		}
		var name, count string
		if len(args) == 1 {
			name = defaultName
			count = args[0]
		} else {
			name = args[0]
			count = args[1]
		}
		loadStdinUserdata()
		provider, err := providers.Get(providerName, true)
		fatal(err)

		hostPattern := defaultHostPattern
		if namer, ok := provider.(HostNamer); ok {
			hostPattern = namer.HostNamePattern()
		}

		hostTemplate := template.Must(template.New("host").Parse(hostPattern))

		existing := existingHosts(hostTemplate, provider, name)
		desired := desiredHosts(hostTemplate, name, count)
		hosts := append(strSet(existing, desired), namespace+name)
		finished := progressBar(".", 2)
		parallelWait(hosts, func(_ int, host string, wg *sync.WaitGroup) {
			defer wg.Done()
			if !strIn(host, desired) {
				fatal(provider.Destroy(host))
				return
			}
			if strIn(host, desired) && !strIn(host, existing) {
				fatal(provider.Create(newHost(host)))
				return
			}
		})
		finished()
	},
}

func hostName(t *template.Template, name string, index interface{}) string {
	var buf bytes.Buffer
	_ = t.Execute(&buf, struct {
		Namespace string
		Name      string
		Index     interface{}
	}{namespace, name, index})
	return buf.String()
}

func desiredHosts(hostTemplate *template.Template, name, count string) []string {
	c, err := strconv.Atoi(count)
	fatal(err)
	var hosts []string
	for i := 0; i < c; i++ {
		hosts = append(hosts, hostName(hostTemplate, name, i))
	}
	return hosts
}

func existingHosts(hostTemplate *template.Template, provider providers.HostProvider, name string) []string {
	var hosts []string
	for _, h := range provider.List(hostName(hostTemplate, name, "*")) {
		hosts = append(hosts, h.Name)
	}
	return hosts
}

func strIn(str string, list []string) bool {
	for i := range list {
		if str == list[i] {
			return true
		}
	}
	return false
}

func strSet(strs ...[]string) []string {
	m := make(map[string]bool)
	for i := range strs {
		for _, str := range strs[i] {
			m[str] = true
		}
	}
	var set []string
	for k := range m {
		set = append(set, k)
	}
	return set
}
