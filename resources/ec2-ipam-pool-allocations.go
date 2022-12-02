package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2IpamPoolAllocation struct {
	svc                *ec2.EC2
	ipamPool           *ec2.IpamPool
	ipamPoolAllocation *ec2.IpamPoolAllocation
	tags               []*ec2.Tag
}

func init() {
	register("EC2IpamPoolAllocation", ListEC2IpamPoolAllocations)
}

func ListEC2IpamPoolAllocations(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeIpamPools(nil)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, ipamPool := range resp.IpamPools {
		allocationsResp, err := svc.GetIpamPoolAllocations(
			&ec2.GetIpamPoolAllocationsInput{
				IpamPoolId: ipamPool.IpamPoolId,
			},
		)
		if err != nil {
			return nil, err
		}

		for _, ipamPoolAllocation := range allocationsResp.IpamPoolAllocations {
			resources = append(resources, &EC2IpamPoolAllocation{
				svc:                svc,
				ipamPool:           ipamPool,
				ipamPoolAllocation: ipamPoolAllocation,
			})
		}

	}

	return resources, nil
}

func (i *EC2IpamPoolAllocation) Remove() error {
	params := &ec2.ReleaseIpamPoolAllocationInput{
		IpamPoolId:           i.ipamPool.IpamPoolId,
		IpamPoolAllocationId: i.ipamPoolAllocation.IpamPoolAllocationId,
		Cidr:                 i.ipamPoolAllocation.Cidr,
	}

	_, err := i.svc.ReleaseIpamPoolAllocation(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2IpamPoolAllocation) String() string {
	return fmt.Sprintf("%s-%s", *i.ipamPool.IpamPoolId, *i.ipamPoolAllocation.IpamPoolAllocationId)
}

func (i *EC2IpamPoolAllocation) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Id", i.ipamPoolAllocation.IpamPoolAllocationId)
	properties.Set("cidr", i.ipamPoolAllocation.Cidr)
	properties.Set("poolId", i.ipamPool.IpamPoolId)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
