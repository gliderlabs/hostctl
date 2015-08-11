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

This is doesn't reflect exactly what's in code. It was used for design:

```
hostctl up <name>             # create instance, blocks until ready
hostctl down <name>           # destroy instance, blocks until terminated
hostctl scale <name> <count>  # scale instance group
hostctl ip <name>             # get instance ip
hostctl ls <pattern>          # list instances
hostctl env                   # show current relevant environment

HOSTCTL_PROVIDER      # what provider backend (digitalocean, ec2)
HOSTCTL_NAMESPACE			# optional namespace for names
HOSTCTL_NAME          # optional default name
HOSTCTL_IMAGE         # vm image
HOSTCTL_FLAVOR				# vm flavor
HOSTCTL_REGION				# vm region
HOSTCTL_KEYNAME				# vm keyname
HOSTCTL_USERDATA			# vm userdata
HOSTCTL_OPTIONS				# vm options
```

## Todo

* move to GL, project infrastructure
* userdata, stdin override
* scale cmd
* up,down support multiple names
* status output for up/down/scale
* flesh out vm options for DO
* global profile option
* version
* docs

## License

MIT
