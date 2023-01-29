package aws

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	awsxec2 "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A collection of values returned by FargateRun.
type GetVPCRunResult struct {
	VPCId   pulumi.StringOutput
	VPCCidr pulumi.StringOutput
}

// New VPC resource with the given unique name, ipam id, ipam-pool id, and netmask length for new VPC.
// https://www.pulumi.com/docs/guides/crosswalk/aws/vpc
func VpcRunIpam(ctx *pulumi.Context, prefixName, ipamID, ipamPoolId string, netMaskLength int) (*GetVPCRunResult, error) {

	var rv GetVPCRunResult

	// Lookup next available CIDR range from being returned by the pool
	previewNextCidr, err := ec2.GetIpamPreviewNextCidr(ctx, &ec2.GetIpamPreviewNextCidrArgs{
		IpamPoolId:    ipamPoolId,
		NetmaskLength: pulumi.IntRef(netMaskLength),
		// TODO: Add a feature to exclude a particular CIDR range
		// DisallowedCidrs: []string{"10.0.0.0/18"},
	}, nil)
	if err != nil {
		return nil, err
	}

	// Create a new VPC with CIDR block from AWS IPAM
	newVpc, err := awsxec2.NewVpc(ctx, prefixName, &awsxec2.VpcArgs{
		CidrBlock: &previewNextCidr.Cidr,
	})
	if err != nil {
		return nil, err
	}

	// Return values
	rv.VPCId = newVpc.VpcId
	rv.VPCCidr = pulumi.StringOutput(newVpc.Vpc.CidrBlock())

	return &rv, nil
}
