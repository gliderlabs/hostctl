package google

type ImageAlias struct {
	Project string
	Name    string
}

var PublicAliases = map[string]ImageAlias{
	"centos-6":           ImageAlias{"centos-cloud", "centos-6"},
	"centos-7":           ImageAlias{"centos-cloud", "centos-7"},
	"container-vm":       ImageAlias{"google-containers", "container-vm"},
	"coreos":             ImageAlias{"coreos-cloud", "coreos-stable"},
	"debian-7":           ImageAlias{"debian-cloud", "debian-7-wheezy"},
	"debian-7-backports": ImageAlias{"debian-cloud", "backports-debian-7-wheezy"},
	"debian-8":           ImageAlias{"debian-cloud", "debian-8-jessie"},
	"opensuse-13":        ImageAlias{"opensuse-cloud", "opensuse-13"},
	"rhel-6":             ImageAlias{"rhel-cloud", "rhel-6"},
	"rhel-7":             ImageAlias{"rhel-cloud", "rhel-7"},
	"sles-11":            ImageAlias{"suse-cloud", "sles-11"},
	"sles-12":            ImageAlias{"suse-cloud", "sles-12"},
	"ubuntu-12-04":       ImageAlias{"ubuntu-os-cloud", "ubuntu-1204-precise"},
	"ubuntu-14-04":       ImageAlias{"ubuntu-os-cloud", "ubuntu-1404-trusty"},
	"ubuntu-14-10":       ImageAlias{"ubuntu-os-cloud", "ubuntu-1410-utopic"},
	"ubuntu-15-04":       ImageAlias{"ubuntu-os-cloud", "ubuntu-1504-vivid"},
	"windows-2008-r2":    ImageAlias{"windows-cloud", "windows-server-2008-r2"},
	"windows-2012-r2":    ImageAlias{"windows-cloud", "windows-server-2012-r2"},
}
