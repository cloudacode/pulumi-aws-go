package vpc

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func VpcRunIpam(ctx *pulumi.Context, prefixName, ipamID, ipamPoolId string) error {
	testVpcIpamPool, err := ec2.GetVpcIpamPool(ctx, ipamID, pulumi.ID(ipamPoolId), &ec2.VpcIpamPoolState{
		AddressFamily: pulumi.String("ipv4"),
	})
	if err != nil {
		return err
	}
	_, err = ec2.NewVpc(ctx, prefixName+"-vpc", &ec2.VpcArgs{
		Ipv4IpamPoolId:    testVpcIpamPool.ID(),
		Ipv4NetmaskLength: pulumi.Int(28),
		Tags:              pulumi.StringMap{"Name": pulumi.String(prefixName + "-vpc")},
	})
	if err != nil {
		return err
	}
	return nil
}
