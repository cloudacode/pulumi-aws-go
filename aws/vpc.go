package aws

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	awsxec2 "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func VpcRunIpam(ctx *pulumi.Context, prefixName, ipamID, ipamPoolId string, netMaskLength int) error {

	// Lookup existing IPAM Pool. https://www.pulumi.com/registry/packages/aws/api-docs/ec2/vpcipampool/
	testVpcIpamPool, err := ec2.GetVpcIpamPool(ctx, ipamID, pulumi.ID(ipamPoolId), &ec2.VpcIpamPoolState{
		AddressFamily: pulumi.String("ipv4"),
	})
	if err != nil {
		return err
	}
	// Create a VPC with CIDR from AWS IPAM
	testVpc, err := ec2.NewVpc(ctx, prefixName+"-vpc", &ec2.VpcArgs{
		Ipv4IpamPoolId:    testVpcIpamPool.ID(),
		Ipv4NetmaskLength: pulumi.Int(netMaskLength),
		Tags:              pulumi.StringMap{"Name": pulumi.String(prefixName + "-vpc")},
	})
	if err != nil {
		return err
	}
	ctx.Export("VPC CIDR: ", testVpc.CidrBlock)

	// nextCidr, err := ec2.NewVpcIpamPreviewNextCidr(ctx, "exampleVpcIpamPreviewNextCidr", &ec2.VpcIpamPreviewNextCidrArgs{
	// 	IpamPoolId:    testVpcIpamPool.ID(),
	// 	NetmaskLength: pulumi.Int(netMaskLength),
	// 	DisallowedCidrs: pulumi.StringArray{
	// 		pulumi.String("10.0.0.0/18"),
	// 	},
	// }, nil)
	// if err != nil {
	// 	return err
	// }
	// ctx.Export("Next CIDR: ", nextCidr.Cidr)

	previewNextCidr, err := ec2.GetIpamPreviewNextCidr(ctx, &ec2.GetIpamPreviewNextCidrArgs{
		IpamPoolId:      ipamPoolId,
		NetmaskLength:   pulumi.IntRef(18),
		DisallowedCidrs: []string{"10.0.0.0/18"},
	}, nil)

	if err != nil {
		return err
	}

	ctx.Export("preview CIDR: ", pulumi.StringPtr(previewNextCidr.Cidr))

	newVpc, err := awsxec2.NewVpc(ctx, prefixName+"-vpc", &awsxec2.VpcArgs{
		CidrBlock: &previewNextCidr.Cidr,
	})

	if err != nil {
		return err
	}

	ctx.Export("newVPC CIDR: ", newVpc.Vpc.CidrBlock())

	return nil
}
