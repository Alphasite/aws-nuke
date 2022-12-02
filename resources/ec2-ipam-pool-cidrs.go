package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2IpamPoolCidr struct {
	svc          *ec2.EC2
	ipamPool     *ec2.IpamPool
	ipamPoolCidr *ec2.IpamPoolCidr
	tags         []*ec2.Tag
}

func init() {
	register("ListEC2IpamPoolCidrs", ListEC2IpamPoolCidrs)
}

func ListEC2IpamPoolCidrs(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeIpamPools(nil)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, ipamPool := range resp.IpamPools {
		allocationsResp, err := svc.GetIpamPoolCidrs(
			&ec2.GetIpamPoolCidrsInput{
				IpamPoolId: ipamPool.IpamPoolId,
			},
		)
		if err != nil {
			return nil, err
		}

		for _, ipamPoolCidr := range allocationsResp.IpamPoolCidrs {
			resources = append(resources, &EC2IpamPoolCidr{
				svc:          svc,
				ipamPool:     ipamPool,
				ipamPoolCidr: ipamPoolCidr,
			})
		}

	}

	return resources, nil
}

func (i *EC2IpamPoolCidr) Remove() error {
	params := &ec2.DeprovisionIpamPoolCidrInput{
		IpamPoolId: i.ipamPool.IpamPoolId,
		Cidr:       i.ipamPoolCidr.Cidr,
	}

	_, err := i.svc.DeprovisionIpamPoolCidr(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2IpamPoolCidr) String() string {
	return fmt.Sprintf("%s-%s", *i.ipamPool.IpamPoolId, *i.ipamPoolCidr.Cidr)
}

func (i *EC2IpamPoolCidr) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("cidr", i.ipamPoolCidr.Cidr)
	properties.Set("poolId", i.ipamPool.IpamPoolId)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
