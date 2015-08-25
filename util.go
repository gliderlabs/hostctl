package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gliderlabs/hostctl/providers"
	"golang.org/x/crypto/ssh/terminal"
)

func newHost(name string) providers.Host {
	return providers.Host{
		Name:     name,
		Flavor:   hostFlavor,
		Image:    hostImage,
		Region:   hostRegion,
		Keyname:  hostKeyname,
		Userdata: hostUserdata,
	}
}

func loadStdinUserdata() {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		data, err := ioutil.ReadAll(os.Stdin)
		fatal(err)
		hostUserdata = string(data)
	}
}

func parallelWait(items []string, fn func(int, string, *sync.WaitGroup)) {
	var wg sync.WaitGroup
	for i := 0; i < len(items); i++ {
		wg.Add(1)
		go fn(i, items[i], &wg)
	}
	wg.Wait()
}

func fatal(err error) {
	if err != nil {
		fmt.Println("!!", err)
		os.Exit(1)
	}
}

func exists(path ...string) bool {
	_, err := os.Stat(filepath.Join(path...))
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	fatal(err)
	return true
}

func expandHome(path string) string {
	if path[:2] == "~/" {
		usr, _ := user.Current()
		path = strings.Replace(path, "~/", usr.HomeDir+"/", 1)
	}
	return path
}

func optArg(args []string, i int, default_ string) string {
	if i+1 > len(args) {
		return default_
	}
	return args[i]
}

func progressBar(unit string, interval time.Duration) func() {
	finished := make(chan bool)
	go func() {
		for {
			select {
			case <-finished:
				return
			case <-time.After(interval * time.Second):
				fmt.Fprint(os.Stderr, unit)
			}
		}
	}()
	return func() {
		finished <- true
		fmt.Fprintln(os.Stderr)
	}
}

func source(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	cmdStr := fmt.Sprintf("source %s 1>&2; env", filepath)
	out, err := exec.Command("bash", "-c", cmdStr).Output()
	if err != nil {
		return err
	}

	fileStr := string(file)
	outLines := strings.Split(string(out), "\n")
	for _, line := range outLines {
		lineSplit := strings.SplitN(line, "=", 2)
		if len(lineSplit) != 2 {
			continue
		}
		if strings.Contains(fileStr, lineSplit[0]) {
			os.Setenv(lineSplit[0], lineSplit[1])
		}
	}
	return nil
}
