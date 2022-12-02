package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type EC2IpamScope struct {
	svc   *ec2.EC2
	scope *ec2.IpamScope
	tags  []*ec2.Tag
}

func init() {
	register("EC2IpamScope", ListEC2IpamScopes)
}

func ListEC2IpamScopes(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	params := &ec2.DescribeIpamScopesInput{MaxResults: aws.Int64(100)}
	resp, err := svc.DescribeIpamScopes(params)
	if err != nil {
		return nil, err
	}
	var resources []Resource
	for _, ipamScope := range resp.IpamScopes {
		resources = append(resources, &EC2IpamScope{
			svc:   svc,
			scope: ipamScope,
			tags:  ipamScope.Tags,
		})

	}

	return resources, nil
}

func (i *EC2IpamScope) Remove() error {
	params := &ec2.DeleteIpamScopeInput{
		IpamScopeId: i.scope.IpamScopeId,
	}

	_, err := i.svc.DeleteIpamScope(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2IpamScope) String() string {
	return *i.scope.IpamScopeId
}

func (i *EC2IpamScope) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Id", i.scope.IpamScopeId)

	for _, tag := range i.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}
