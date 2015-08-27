# Hostctl

The simplest way to boot and manage cloud VMs.

[![Circle CI](https://circleci.com/gh/gliderlabs/hostctl.png?style=shield)](https://circleci.com/gh/gliderlabs/hostctl)
[![IRC Channel](https://img.shields.io/badge/irc-%23gliderlabs-blue.svg)](https://kiwiirc.com/client/irc.freenode.net/#gliderlabs)
<br /><br />

Hostctl is an opinionated CLI tool for basic cloud VM operations, ideal for
development and personal use. It does nothing more than manage VM hosts, so if
you need anything else you should look at cloud provider specific tools. Hostctl
supports pluggable cloud providers, currently including DigitalOcean and Amazon
EC2.

## Getting Hostctl

You can install the Hostctl binary right into a directory in your `$PATH`:
```
$ curl https://dl.gliderlabs.com/gh/hostctl/latest/$(uname -sm|tr \  _).tgz \
    | tar -zxC /usr/local/bin
```
## Using Hostctl

The quickest way to see Hostctl in action is our
[Quickstart](user/quickstart.md) tutorial. After [configuring for a particular
provider](user/providers.md), spinning up a VM is as simple as:

    $ hostctl up <vm-name>

For a full list of command, see the [Command Reference](user/reference.md) in
the User Guide.

## Contributing

Pull requests are welcome! We recommend getting feedback before starting by
opening a [GitHub issue](https://github.com/gliderlabs/hostctl/issues) or
discussing in [Slack](http://glider-slackin.herokuapp.com/).

Also check out our Developer Guide on [Contributing Providers](dev/providers.md)
and [Staging Releases](dev/releases.md).

## License

MIT

<img src="https://ga-beacon.appspot.com/UA-58928488-2/hostctl/readme?pixel" />
