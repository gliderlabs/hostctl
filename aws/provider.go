package aws

import (
	"fmt"
	"os"
	"time"

	"github.com/MattAitchison/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/gliderlabs/hostctl/providers"
)

func init() {
	providers.Register(&awsProvider{}, "aws")
}

type awsProvider struct {
	client *ec2.EC2
}

// Setup ec2.Client using aws supported credentials (env, credential file)
func (p *awsProvider) Setup() error {
	region := os.Getenv("HOSTCTL_REGION")
	if region == "" {
		return fmt.Errorf("HOSTCTL_REGION required")
	}
	config := aws.NewConfig().WithRegion(region)
	p.client = ec2.New(config)
	return nil
}

func (p *awsProvider) Env() *env.EnvSet {
	var envSet = env.NewEnvSet("aws")
	envSet.Secret("AWS_ACCESS_KEY", "access key for AWS")
	envSet.Secret("AWS_SECRET_KEY", "secret key for AWS")
	envSet.String("AWS_AVAILABILITY_ZONE", "", "availability zone for AWS; eg: us-west-2a")
	return envSet
}

// Create an instance based on a Host, poll until instance.Status=running
func (p *awsProvider) Create(host providers.Host) error {
	zone := os.Getenv("AWS_AVAILABILITY_ZONE")
	res, err := p.client.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(host.Image),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		UserData:     aws.String(host.Userdata),
		KeyName:      aws.String(host.Keyname),
		InstanceType: aws.String(host.Flavor),
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(zone),
		},
	})
	if err != nil {
		return err
	}
	if res == nil || len(res.Instances) == 0 {
		return fmt.Errorf("no instances created")
	}
	instanceID := res.Instances[0].InstanceId
	_, err = p.client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{instanceID},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(host.Name),
			},
		},
	})
	if err != nil {
		return err
	}
	return p.pollFor(*instanceID, "running")
}

// Destroy the first instance with tag:Name=name
func (p *awsProvider) Destroy(name string) error {
	for id, host := range p.list(name) {
		if host.Name == name {
			return p.destroy(id)
		}
	}
	return nil
}

// destroy instance by ID polling until instance.State=terminated
func (p *awsProvider) destroy(id string) error {
	_, err := p.client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	})
	if err != nil {
		return err
	}
	return p.pollFor(id, "terminated")
}

func (p *awsProvider) List(pattern string) []providers.Host {
	var hosts []providers.Host
	for _, host := range p.list(pattern) {
		hosts = append(hosts, host)
	}
	return hosts
}

// list all instances with tag:Name matching pattern.
// Will NOT return any instances where PublicIpAddress is nil.
// Map key will be set to the instanceID.
func (p *awsProvider) list(pattern string) map[string]providers.Host {
	res, err := p.client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(pattern)},
			},
		},
	})
	if err != nil {
		return nil
	}
	hosts := make(map[string]providers.Host)
	for i := range res.Reservations {
		for _, instance := range res.Reservations[i].Instances {
			// Note: an instance that is terminated/terminating will not have a public IP.
			if instance != nil && instance.PublicIpAddress != nil {
				id := *instance.InstanceId
				hosts[id] = providers.Host{
					Name:    nametag(instance),
					IP:      *instance.PublicIpAddress,
					Region:  *instance.Placement.AvailabilityZone,
					Image:   *instance.ImageId,
					Keyname: *instance.KeyName,
					Flavor:  *instance.InstanceType,
				}
			}
		}
	}
	return hosts
}

// Get the first instance with tag:Name=name
func (p *awsProvider) Get(name string) *providers.Host {
	hosts := p.List(name)
	for _, host := range hosts {
		if host.Name == name {
			return &host
		}
	}
	return nil
}

// pollFor will poll every 2 seconds for an instance by ID,
// until instance.State=state, timeout or an error.
func (p *awsProvider) pollFor(id string, state string) error {
	timeout := time.After(1 * time.Minute)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("aws provider: timed out wating for state: %s", state)
		default:
		}
		res, err := p.client.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(id)},
		})
		if err != nil {
			return err
		}
		for i := range res.Reservations {
			for _, instance := range res.Reservations[i].Instances {
				if *instance.State.Name == state {
					return nil
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func nametag(instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}
