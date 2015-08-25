package aws

import (
	"fmt"

	"github.com/MattAitchison/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/gliderlabs/hostctl/providers"
)

var envSet = env.NewEnvSet("aws")

func init() {
	providers.Register(&awsProvider{}, "aws")
}

type awsProvider struct {
	client *ec2.EC2
}

func (p *awsProvider) Setup() error {
	p.client = ec2.New(nil)
	return nil
}

func (p *awsProvider) Env() *env.EnvSet {
	return envSet
}

func (p *awsProvider) Create(host providers.Host) error {
	res, err := p.client.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(host.Image),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),
		UserData:     aws.String(host.Userdata),
		KeyName:      aws.String(host.Keyname),
		InstanceType: aws.String(host.Flavor),
		Placement: &ec2.Placement{
			AvailabilityZone: aws.String(host.Region),
		},
	})
	if err != nil {
		return err
	}

	if res != nil && len(res.Instances) > 0 {
		return fmt.Errorf("no instances created")
	}
	_, err = p.client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{
			res.Instances[0].InstanceId,
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(host.Name),
			},
		},
	})
	return err
}

func (p *awsProvider) Destroy(name string) error {
	for id, host := range p.list(name) {
		if host.Name == name {
			return p.destroy(id)
		}
	}
	return nil
}

func (p *awsProvider) destroy(id string) error {
	_, err := p.client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	})
	return err
}

func (p *awsProvider) List(pattern string) []providers.Host {
	var hosts []providers.Host
	for _, host := range p.list(pattern) {
		hosts = append(hosts, host)
	}
	return hosts
}

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
			if instance != nil && instance.PublicIpAddress != nil {
				id := *instance.InstanceId
				hosts[id] = providers.Host{
					Name: nametag(instance),
					IP:   *instance.PublicIpAddress,
				}
			}
		}
	}
	return hosts
}

func (p *awsProvider) Get(name string) *providers.Host {
	hosts := p.List(name)
	for _, host := range hosts {
		if host.Name == name {
			return &host
		}
	}
	return nil
}

func nametag(instance *ec2.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}
