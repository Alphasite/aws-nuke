package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2IpamPool struct {
	svc      *ec2.EC2
	ipamPool *ec2.IpamPool
	tags     []*ec2.Tag
}

func init() {
	register("EC2IpamPool", ListEC2IpamPools)
}

func ListEC2IpamPools(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeIpamPoolsInput{MaxResults: aws.Int64(100)}
	resp, err := svc.DescribeIpamPools(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, ipamPool := range resp.IpamPools {
		resources = append(resources, &EC2IpamPool{
			svc:      svc,
			ipamPool: ipamPool,
			tags:     ipamPool.Tags,
		})

	}

	return resources, nil
}

func (i *EC2IpamPool) Remove() error {
	params := &ec2.DeleteIpamPoolInput{
		IpamPoolId: i.ipamPool.IpamPoolId,
	}

	_, err := i.svc.DeleteIpamPool(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2IpamPool) String() string {
	return *i.ipamPool.IpamPoolId
}

func (i *EC2IpamPool) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Id", i.ipamPool.IpamPoolId)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
