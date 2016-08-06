package plugin

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/cloudconfigs"
)

type VMProperties struct {
	MachineType    string `yaml:"machine_type,omitempty"`
	RootDiskSizeGB int    `yaml:"root_disk_size_gb"`
	RootDiskType   string `yaml:"root_disk_type"`
}
type NetworkCloudProperties struct {
	NetworkName    string   `yaml:"network_name"`
	SubNetworkName string   `yaml:"subnetwork_name"`
	Tags           []string `yaml:"tags"`
}
type GCPCloudConfig struct {
	Context *cli.Context
}

func NewGCPCloudConfig(c *cli.Context) cloudconfigs.CloudConfigProvider {

	provider := &GCPCloudConfig{
		Context: c,
	}
	return provider
}

func (c *GCPCloudConfig) CreateNetworks() ([]enaml.DeploymentNetwork, error) {
	context := c.Context
	networks := []enaml.DeploymentNetwork{}
	for i := 1; i <= SupportedNetworkCount; i++ {
		networkFlag := fmt.Sprintf("network-name-%d", i)
		if context.IsSet(networkFlag) {
			network := enaml.ManualNetwork{
				Name: context.String(networkFlag),
				Type: "manual",
			}
			azs := context.StringSlice(fmt.Sprintf("network-az-%d", i))
			if err := cloudconfigs.CheckRequiredLength(len(azs), i, context, "network-cidr-%d", "network-gateway-%d", "network-dns-%d", "network-reserved-%d", "network-static-%d"); err != nil {
				return nil, err
			}
			ranges := context.StringSlice(fmt.Sprintf("network-cidr-%d", i))
			gateways := context.StringSlice(fmt.Sprintf("network-gateway-%d", i))
			dnsServers := context.StringSlice(fmt.Sprintf("network-dns-%d", i))
			reservedRanges := context.StringSlice(fmt.Sprintf("network-reserved-%d", i))
			staticIPs := context.StringSlice(fmt.Sprintf("network-static-%d", i))
			if err := cloudconfigs.CheckRequiredLength(len(azs), i, context, "gcp-network-name-%d", "gcp-subnetwork-name-%d", "gcp-network-tag-%d"); err != nil {
				return nil, err
			}
			gcpNetworkNames := context.StringSlice(fmt.Sprintf("gcp-network-name-%d", i))
			gcpSubNetworkNames := context.StringSlice(fmt.Sprintf("gcp-subnetwork-name-%d", i))
			gcpNetworkTags := context.StringSlice(fmt.Sprintf("gcp-network-tag-%d", i))
			for index, az := range azs {
				subnet := enaml.Subnet{
					AZ:       az,
					Range:    ranges[index],
					Gateway:  gateways[index],
					DNS:      strings.Split(dnsServers[index], ","),
					Reserved: strings.Split(reservedRanges[index], ","),
					Static:   strings.Split(staticIPs[index], ","),
					CloudProperties: NetworkCloudProperties{
						NetworkName:    gcpNetworkNames[index],
						SubNetworkName: gcpSubNetworkNames[index],
						Tags:           strings.Split(gcpNetworkTags[index], ","),
					},
				}
				network.AddSubnet(subnet)
			}
			networks = append(networks, network)
		}
	}
	return networks, nil
}
func (c *GCPCloudConfig) CreateAZs() ([]enaml.AZ, error) {
	azNames := c.Context.StringSlice("az")
	gcpAZNames := c.Context.StringSlice("gcp-availability-zone")

	if len(azNames) != len(gcpAZNames) {
		err := fmt.Errorf("Sorry you need to provide the same number of az and gcp-availability-zone flags")
		return nil, err
	}
	azs := []enaml.AZ{}

	for i, azName := range azNames {
		az := enaml.AZ{
			Name: azName,
			CloudProperties: map[string]string{
				"availability_zone": gcpAZNames[i],
			},
		}
		azs = append(azs, az)
	}
	return azs, nil
}

func (c *GCPCloudConfig) CreateVMTypes() ([]enaml.VMType, error) {
	vmTypes := []enaml.VMType{
		enaml.VMType{
			Name: "small",
			CloudProperties: VMProperties{
				MachineType:    "n1-standard-1",
				RootDiskSizeGB: 30,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.VMType{
			Name: "medium",
			CloudProperties: VMProperties{
				MachineType:    "n1-standard-2",
				RootDiskSizeGB: 50,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.VMType{
			Name: "large.memory",
			CloudProperties: VMProperties{
				MachineType:    "n1-highmem-4",
				RootDiskSizeGB: 50,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.VMType{
			Name: "large.cpu",
			CloudProperties: VMProperties{
				MachineType:    "n1-highcpu-4",
				RootDiskSizeGB: 30,
				RootDiskType:   "pd-standard",
			},
		},
	}
	return vmTypes, nil
}

func (c *GCPCloudConfig) CreateDiskTypes() ([]enaml.DiskType, error) {
	diskTypes := []enaml.DiskType{
		enaml.DiskType{
			Name:     "small",
			DiskSize: 3000,
			CloudProperties: VMProperties{
				RootDiskSizeGB: 3,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.DiskType{
			Name:     "medium",
			DiskSize: 30000,
			CloudProperties: VMProperties{
				RootDiskSizeGB: 50,
				RootDiskType:   "pd-standard",
			},
		},
		enaml.DiskType{
			Name:     "large",
			DiskSize: 50000,
			CloudProperties: VMProperties{
				RootDiskSizeGB: 50,
				RootDiskType:   "pd-standard",
			},
		},
	}
	return diskTypes, nil
}

func (c *GCPCloudConfig) CreateCompilation() (*enaml.Compilation, error) {
	compilation := &enaml.Compilation{
		Workers:             8,
		ReuseCompilationVMs: true,
		AZ:                  c.Context.StringSlice("network-az-1")[0],
		VMType:              "medium",
		Network:             c.Context.String("network-name-1"),
	}
	return compilation, nil
}
