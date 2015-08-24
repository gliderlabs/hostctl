package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/gliderlabs/hostctl/providers"
)

func testRunCmd(t *testing.T, cmdline string, statusExpected int, provider *providers.TestProvider, setupFn func()) (*bytes.Buffer, *bytes.Buffer) {
	if os.Getenv("TEST_CMD") != "" {
		if setupFn != nil {
			setupFn()
		}
		p := registerTestProvider()
		os.Args = strings.Split(os.Getenv("TEST_CMD"), " ")
		fatal(Hostctl.Execute())
		returnTestProvider(p)
		os.Exit(0)
	}
	pc, _, _, _ := runtime.Caller(1)
	callerPath := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	testName := callerPath[len(callerPath)-1]
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append(os.Environ(), "TEST_CMD="+cmdline)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	f := writeTestProvider(provider, testName)
	if f != nil {
		cmd.Env = append(cmd.Env, "TEST_PROVIDER="+f.Name())
	}
	err := cmd.Run()
	status := 0
	if exiterr, ok := err.(*exec.ExitError); ok {
		if s, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			status = s.ExitStatus()
		}
	}
	if status != statusExpected {
		fmt.Println(stderr.String())
		ensure.DeepEqual(t, status, statusExpected)
	}
	readTestProvider(provider, f)
	return &stdout, &stderr
}

func registerTestProvider() *providers.TestProvider {
	if os.Getenv("TEST_PROVIDER") != "" {
		providerFile, err := os.Open(os.Getenv("TEST_PROVIDER"))
		if err != nil {
			panic(err)
		}
		dec := gob.NewDecoder(providerFile)
		var testProvider providers.TestProvider
		err = dec.Decode(&testProvider)
		if err != nil {
			panic(err)
		}
		providerFile.Close()
		providerName = "test"
		providers.Register(&testProvider, "test")
		return &testProvider
	}
	return nil
}

func returnTestProvider(provider *providers.TestProvider) {
	if os.Getenv("TEST_PROVIDER") != "" {
		providerFile, err := os.Create(os.Getenv("TEST_PROVIDER"))
		if err != nil {
			panic(err)
		}
		enc := gob.NewEncoder(providerFile)
		err = enc.Encode(provider)
		if err != nil {
			panic(err)
		}
		providerFile.Close()
	}
}

func writeTestProvider(provider *providers.TestProvider, prefix string) *os.File {
	if provider != nil {
		providerFile, err := ioutil.TempFile("", prefix)
		if err != nil {
			panic(err)
		}
		enc := gob.NewEncoder(providerFile)
		err = enc.Encode(provider)
		if err != nil {
			panic(err)
		}
		providerFile.Close()
		return providerFile
	}
	return nil
}

func readTestProvider(provider *providers.TestProvider, f *os.File) {
	if provider != nil {
		providerFile, err := os.Open(f.Name())
		if err != nil {
			panic(err)
		}
		dec := gob.NewDecoder(providerFile)
		var testProvider providers.TestProvider
		err = dec.Decode(&testProvider)
		if err != nil {
			panic(err)
		}
		providerFile.Close()
		os.Remove(f.Name())
		*provider = testProvider
	}
}

func TestHostctlCmd(t *testing.T) {
	t.Parallel()
	stdout, stderr := testRunCmd(t, "hostctl", 0, nil, nil)
	ensure.StringContains(t, stdout.String(), Hostctl.Short)
	ensure.StringContains(t, stdout.String(), "Usage:")
	ensure.StringContains(t, stdout.String(), "Available Commands:")
	ensure.StringContains(t, stdout.String(), "Flags:")
	ensure.DeepEqual(t, stderr.String(), "")
}

func TestVersionCmd(t *testing.T) {
	t.Parallel()
	stdout, stderr := testRunCmd(t, "hostctl version", 0, nil, func() {
		Version = "0"
	})
	ensure.DeepEqual(t, stdout.String(), "0\n")
	ensure.DeepEqual(t, stderr.String(), "")
}
