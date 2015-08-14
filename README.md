# hostctl

Hostctl is an opinionated command line tool for easily provisioning cloud VMs.

Hostctl is ideal for spinning up VMs for development or personal use. It does
nothing more than manage VM hosts, so if you need anything else you should look
at cloud provider specific tools. It's not intended for managing production
clusters, as you should be using a tool like [Terraform](https://terraform.io/) instead.

## Getting hostctl

Until the first release, you can get hostctl with `go get`:

    $ go get github.com/progrium/hostctl

## Usage

```
Usage:
  hostctl [command]

Available Commands:
  down        Terminate host
  env         Show relevant environment
  ip          Show IP for host
  ls          List hosts
  scale       Resize host cluster
  ssh         SSH to host
  up          Provision host, wait until ready
  help        Help about any command
```

## Configuration

```
HOSTCTL_PROVIDER      # what provider backend (digitalocean, ec2)
HOSTCTL_IMAGE         # vm image
HOSTCTL_FLAVOR				# vm flavor
HOSTCTL_REGION				# vm region
HOSTCTL_KEYNAME				# vm keyname
HOSTCTL_USERDATA			# vm userdata
HOSTCTL_NAMESPACE			# optional namespace for names
HOSTCTL_NAME          # optional default name
HOSTCTL_USER          # ssh user
```

## Todo

* move to GL, project infrastructure
* tests
* docs

## License

MIT
