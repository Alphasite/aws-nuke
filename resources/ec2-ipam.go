package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2Ipam struct {
	svc  *ec2.EC2
	ipam *ec2.Ipam
	tags []*ec2.Tag
}

func init() {
	register("EC2Ipam", ListEC2Ipam)
}

func ListEC2Ipam(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeIpamsInput{MaxResults: aws.Int64(100)}
	resp, err := svc.DescribeIpams(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, ipam := range resp.Ipams {
		resources = append(resources, &EC2Ipam{
			svc:  svc,
			ipam: ipam,
			tags: ipam.Tags,
		})

	}

	return resources, nil
}

func (i *EC2Ipam) Remove() error {
	params := &ec2.DeleteIpamInput{
		IpamId: i.ipam.IpamId,
	}

	_, err := i.svc.DeleteIpam(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2Ipam) String() string {
	return *i.ipam.IpamId
}

func (i *EC2Ipam) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Id", i.ipam.IpamId)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
