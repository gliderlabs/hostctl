# hostctl

Hostctl is an opinionated command line tool for just provisioning cloud VMs easily.

```
hostctl up <name>					# create instance
hostctl down <name>					# destroy instance
hostctl scale <name> <count>		# scale instance group
hostctl ip <name>					# get instance ip
hostctl ls <filter>					# list instances
hostctl env							# show environment/config


HOSTCTL_PROVIDER			# what provider backend (digitalocean, ec2)
HOSTCTL_NAME				# default name
HOSTCTL_IMAGE				# vm image
HOSTCTL_FLAVOR				# vm flavor
HOSTCTL_REGION				# vm region
HOSTCTL_KEYNAME				# vm keyname
HOSTCTL_USERDATA			# vm userdata
HOSTCTL_NAMESPACE			# vm namespace for friendly name
HOSTCTL_OPTIONS				# extra options
```

## License

MIT
