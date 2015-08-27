# Quickstart

This is a short, simple tutorial intended to get you started with Hostctl as
quickly as possible. Alternatively, you can skip ahead to the [Command Reference](reference.md).

## Overview

Hostctl lets you spin up and down VMs with cloud providers like DigitalOcean and EC2. You can also use Hostctl to easily get their IP, SSH to them, and scale them. The goal is to make it easier when doing development and experiments.

In this tutorial, we're going to use Hostctl with DigitalOcean to tour all the functionality it provides.

## Installing

If you haven't already, you want to download Hostctl for your platform. You can find the latest version on the [releases page](https://github.com/gliderlabs/hostctl/releases), or you can install from this URL providing your platform. Using `curl` and `tar` you can install `hostctl` right into a directory in your `$PATH`:
```
$ curl https://dl.gliderlabs.com/qs/hostctl/latest/$(uname -sm|tr \  _).tgz \
    | tar -zxC /usr/local/bin
```
You can change `/usr/local/bin` as necessary. Now you should be able to run `hostctl` and see usage and available commands:

    $ hostctl

## Environment

VMs require a lot of parameters, which is one of the annoying parts of booting a VM that slows you down. Hostctl addresses this by letting you define parameters upfront in your environment, which you can put into loadable profiles.

You can see everything you can configure in the environment with `hostctl env`:
```
$ hostctl env
HOSTCTL_IMAGE=""                         # vm image
HOSTCTL_REGION=""                        # vm region
HOSTCTL_KEYNAME=""                       # vm keyname
HOSTCTL_USER="progrium"                  # ssh user
HOSTCTL_PROVIDER="digitalocean"          # cloud provider
HOSTCTL_NAME=""                          # optional default name
HOSTCTL_NAMESPACE=""                     # optional namespace for names
HOSTCTL_FLAVOR=""                        # vm flavor
HOSTCTL_USERDATA=""                      # vm user data
DO_TOKEN=""                              # token for DigitalOcean API v2
```
We see they're mostly empty except for some defaults. `HOSTCTL_USER` defaults to your system's logged in user. `HOSTCTL_PROVIDER` defaults to DigitalOcean.

## Cloud Provider Setup

The values you want for the rest are going to depend on the provider. We're going to focus on DigitalOcean. You can see our [Provider Reference](providers.md) to see what values you want for EC2, for example.

There's two values you'll have to lookup for this to work: `DO_TOKEN` which is a personal access token for the API, and `HOSTCTL_KEYNAME` which is typically the fingerprint (or ID) of the public key you want to use.

If you don't have a personal access token, or to make a new one, you can go to [Applications & API](https://cloud.digitalocean.com/settings/applications) once logged into DigitalOcean. The fingerprint values for your SSH keys are under [Security](https://cloud.digitalocean.com/settings/security) in your account settings.

It's easiest to write these to the global Hostctl profile, a simple shell config script at `~/.hostctl` always sourced by `hostctl`. We'll write two lines to it:

    $ echo "export DO_TOKEN=your-token" >> ~/.hostctl
    $ echo "export HOSTCTL_KEYNAME=your-ssh-key-fingerprint" >> ~/.hostctl

## Base VM Attributes

The rest of the required environment defines your VM. `HOSTCTL_IMAGE`, `HOSTCTL_FLAVOR`, and `HOSTCTL_REGION`. For DigitalOcean, these are slug values from the API. In our case, we'll just use these:

    $ export HOSTCTL_IMAGE=ubuntu
    $ export HOSTCTL_FLAVOR=512mb
    $ export HOSTCTL_REGION=nyc1
    $ export HOSTCTL_USER=root

We also set `HOSTCTL_USER` since DigitalOcean's default user is `root`. We're setting these in our terminal session environment, but they could also be defined anywhere else in your environment. Later we'll see how we can make them into loadable profiles.

We can run `hostctl env` again to see its current configuration.

## Provisioning

Let's make a VM called `demo`:

    $ hostctl up demo

The command will wait until the VM is ready to go, showing a progress bar while you wait. Take a moment to think about how easy that was.

## Everything Else

We can list our VMs and see `demo`, as well as any other VMs you might have running on this account:

    $ hostctl ls
    demo

We can get the IP for `demo` very easily:

    $ hostctl ip demo
    198.5.101.164

We could use this with SSH to connect to it by name:

    $ ssh root@$(hostctl ip demo)

Although it's easier to use the builtin convenience command:

    $ hostctl ssh demo

We can boot more VMs with `hostctl up` or we could create a cluster with `hostctl scale`:

    $ hostctl scale node 3
    ..........................................
    $ hostctl ls
    demo
    node.0
    node.1
    node.2

We can see it made 3 hosts called `node`. We could use scale again to resize the cluster or just scale to nothing:

    $ hostctl scale node 0
    ..........................................
    $ hostctl ls
    demo

But we can also just shutdown any hosts with `hostctl down`:

    $ hostctl down demo

## Profiles

With our VM attributes in the environment, we can iteratively change them. Maybe I want the same VM, but using the `docker` image instead of `ubuntu`. Just set it:

    $ export HOSTCTL_IMAGE=docker
    $ hostctl up docker-vm1

However, if we close this session and come back tomorrow, we won't have this environment. We'd have to set it all up again. But we can write our current environment to a profile that we can use later. We do this with `hostctl env --export`:

    $ hostctl env --export > docker.profile

If you look at the contents of `docker.profile`, it's a shell script that exports everything that was in your current hostctl environment. Now tomorrow we can use the profile without setting any environment:

    $ hostctl -p docker.profile up docker-vm2

Profiles make it easy to iteratively build up a VM configuration and then save it to a file you can use later. You can also write profiles from scratch. They're just shell scripts, so the above is basically the same as:

    $ source docker.profile
    $ hostctl up docker-vm2

## Next Steps

And that's not all. Check out the [Command Reference](reference.md) for what else you can do with these commands, or [Provider Reference](providers.md) to try EC2.
