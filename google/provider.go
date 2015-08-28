package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/MattAitchison/env"
	"github.com/gliderlabs/hostctl/providers"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var envSet = env.NewEnvSet("google")

func init() {
	readEnv()
	providers.Register(new(googleCloudProvider), "google")
}

// default config using client credentials from the "gcloud" tool
var defaultOAuthConfig = oauth2.Config{
	ClientID:     "32555940559.apps.googleusercontent.com",
	ClientSecret: "ZmssLNjJy2998hD4CTg2ejr2",
	Endpoint:     google.Endpoint,
	Scopes: []string{
		compute.DevstorageFullControlScope,
		compute.ComputeScope,
	},
}

func readEnv() {
	envSet.Clear()
	envSet.String("GOOGLE_PROJECT", "", "Google project identifier")
	envSet.String("GOOGLE_OAUTH_TOKEN", "", "OAuth token (saved automatically to ~/.hostctl after authentication)")
}

type googleCloudProvider struct {
	service *compute.Service
	project string
	region  string
}

func (p *googleCloudProvider) HostNamePattern() string {
	return "{{.Namespace}}{{.Name}}-{{.Index}}"
}

func (p *googleCloudProvider) Setup() error {
	readEnv()

	p.region = os.Getenv("HOSTCTL_REGION")
	p.project = envSet.Var("GOOGLE_PROJECT").Value.Get().(string)

	if p.region == "" {
		return errors.New("HOSTCTL_REGION required for Google provider")
	}

	if p.project == "" {
		return errors.New("GOOGLE_PROJECT required for Google provider")
	}

	httpClient, err := newOAuthClient(context.Background())
	if err != nil {
		return err
	}

	service, err := compute.New(httpClient)
	if err != nil {
		return err
	}

	p.service = service
	return nil
}

func (p *googleCloudProvider) Env() *env.EnvSet {
	readEnv()
	return envSet
}

func (p *googleCloudProvider) Create(host providers.Host) error {
	image, err := p.resolveAlias(host.Image)
	if err != nil {
		return err
	}

	instance := &compute.Instance{
		Name:        host.Name,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", host.Region, host.Flavor),
		Disks: []*compute.AttachedDisk{{
			// AutoDelete: true,
			Boot: true,
			InitializeParams: &compute.AttachedDiskInitializeParams{
				SourceImage: image,
			},
		}},
		NetworkInterfaces: []*compute.NetworkInterface{{
			Network: "global/networks/default",
			AccessConfigs: []*compute.AccessConfig{{
				Type: "ONE_TO_ONE_NAT",
			}},
		}},
	}

	// TODO ssh key for host.Keyname

	if host.Userdata != "" {
		instance.Metadata = &compute.Metadata{
			Items: []*compute.MetadataItems{
				{Key: "user-data", Value: &host.Userdata},
			},
		}
	}

	op, err := p.service.Instances.Insert(p.project, p.region, instance).Do()

	if err != nil {
		return err
	}

	return p.waitForZoneOp(op)
}

func (p *googleCloudProvider) Destroy(name string) error {
	op, err := p.service.Instances.Delete(p.project, p.region, name).Do()
	if err != nil {
		if err, ok := err.(*googleapi.Error); ok && err.Code == http.StatusNotFound {
			return nil
		}
		return err
	}
	return p.waitForZoneOp(op)
}

func (p *googleCloudProvider) List(pattern string) (hosts []providers.Host) {
	list, err := p.service.Instances.List(p.project, p.region).Do()
	if err != nil {
		return nil
	}

	for _, instance := range list.Items {
		// TODO convert glob pattern into regex for server-side filtering?
		if ok, _ := filepath.Match(pattern, instance.Name); ok {
			// TODO filter hosts without a public IP?
			hosts = append(hosts, providers.Host{
				Name: instance.Name,
			})
		}
	}

	return hosts
}

func (p *googleCloudProvider) Get(name string) *providers.Host {
	instance, err := p.service.Instances.Get(p.project, p.region, name).Do()
	if err != nil {
		return nil
	}

	for _, iface := range instance.NetworkInterfaces {
		for _, access := range iface.AccessConfigs {
			return &providers.Host{
				Name: instance.Name,
				IP:   access.NatIP,
			}
		}
	}

	return nil
}

func (p *googleCloudProvider) resolveAlias(image string) (string, error) {
	// assume any image with a '/' is already full-qualified
	if strings.ContainsRune(image, '/') {
		return image, nil
	}

	alias, ok := PublicAliases[image]
	// if there is no public alias, assume this is a valid image name within the project
	if !ok {
		return imagePath(p.project, image), nil
	}

	publicImageList, err := p.service.Images.List(alias.Project).Filter(fmt.Sprintf("name eq ^%s(-.+)*-v.+", alias.Name)).Do()
	if err != nil {
		return "", err
	}

	imageVersion := func(image *compute.Image) string {
		parts := strings.Split(image.Name, "v")
		return parts[len(parts)-1]
	}

	var publicImage *compute.Image
	for _, image := range publicImageList.Items {
		if image.Deprecated == nil && (publicImage == nil || imageVersion(image) > imageVersion(publicImage)) {
			publicImage = image
		}
	}

	if publicImage == nil {
		return "", fmt.Errorf("could not find image for alias %s", image)
	}

	userImageList, err := p.service.Images.List(p.project).Filter(fmt.Sprintf("name eq ^%s$", image)).Do()
	if err != nil {
		return "", err
	}

	var userImage *compute.Image
	for _, image := range userImageList.Items {
		if image.Deprecated == nil {
			userImage = image
			break
		}
	}

	if userImage == nil {
		return publicImage.SelfLink, nil
	}

	return "", fmt.Errorf(`Image name "%s" is ambiguous, please use one of the fully-qualified names:
Your image:
    %s
Public image:
    %s
`, image, imagePath(p.project, userImage.Name), imagePath(alias.Project, publicImage.Name))
}

func imagePath(project, name string) string {
	return fmt.Sprintf("projects/%s/global/images/%s", project, name)
}

func (p *googleCloudProvider) waitForZoneOp(op *compute.Operation) (err error) {
	opName := op.Name

	for op.Status != "DONE" {
		time.Sleep(1 * time.Second)

		op, err = p.service.ZoneOperations.Get(p.project, p.region, opName).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func newOAuthClient(ctx context.Context) (*http.Client, error) {
	var token *oauth2.Token

	tokenJSON := envSet.Var("GOOGLE_OAUTH_TOKEN").Value.Get().(string)

	if tokenJSON != "" {
		token = new(oauth2.Token)
		err := json.Unmarshal([]byte(tokenJSON), token)
		if err != nil {
			return nil, err
		}
	} else {
		config, sdkErr := google.NewSDKConfig("")
		if sdkErr == nil {
			return config.Client(ctx), nil
		}
		fmt.Printf(`Unable to read "gcloud" authentication:
    %s

Attempting OAuth online authentication.

For alternative authentication, please use:
    gcloud auth login

`, sdkErr)

		token = tokenFromWeb(ctx, &defaultOAuthConfig)
		saveToken(token)
	}

	return defaultOAuthConfig.Client(ctx, token), nil
}

func saveToken(token *oauth2.Token) {
	path := "~/.hostctl"
	fmt.Println("Success! Saving OAuth token to", path)

	path, _ = homedir.Expand(path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	// if the file is not empty, make sure it ends with a newline
	f.Seek(-1, os.SEEK_END)
	buf := make([]byte, 1)
	n, err := f.Read(buf)
	if n == 0 {
		if err != io.EOF {
			return
		}
	} else if buf[0] != '\n' {
		f.WriteString("\n")
	}

	tokenJSON, _ := json.Marshal(token)
	fmt.Fprintf(f, "export GOOGLE_OAUTH_TOKEN='%s'\n", tokenJSON)
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		http.Error(rw, "", 500)
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	go openURL(authURL)
	fmt.Printf("Opening authentication at:\n    %s\n", authURL)
	code := <-ch

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	return token
}

func openURL(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
}
