package fargate

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	awsxecs "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// reference https://www.pulumi.com/docs/guides/crosswalk/aws/ecs/#creating-an-ecs-cluster-in-a-vpc
func FargateRun(ctx *pulumi.Context, vpcId, prefixName string) error {

	vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{Id: &vpcId})
	if err != nil {
		return err
	}
	subnet, err := ec2.GetSubnetIds(ctx, &ec2.GetSubnetIdsArgs{VpcId: vpc.Id})
	if err != nil {
		return err
	}
	securityGroup, err := ec2.NewSecurityGroup(ctx, prefixName+"-sg", &ec2.SecurityGroupArgs{
		VpcId: pulumi.String(vpc.Id),
		Egress: ec2.SecurityGroupEgressArray{
			&ec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				Protocol: pulumi.String("-1"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
				Ipv6CidrBlocks: pulumi.StringArray{
					pulumi.String("::/0"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	cluster, err := ecs.NewCluster(ctx, prefixName+"-cluster", nil)
	if err != nil {
		return err
	}
	_, err = awsxecs.NewFargateService(ctx, prefixName+"-service", &awsxecs.FargateServiceArgs{
		Cluster: cluster.Arn,
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			Subnets: toPulumiStringArray(subnet.Ids),
			SecurityGroups: pulumi.StringArray{
				securityGroup.ID(),
			},
		},
		DesiredCount: pulumi.Int(2),
		TaskDefinitionArgs: &awsxecs.FargateServiceTaskDefinitionArgs{
			Container: &awsxecs.TaskDefinitionContainerDefinitionArgs{
				Image:     pulumi.String("nginx:latest"),
				Cpu:       pulumi.Int(256),
				Memory:    pulumi.Int(512),
				Essential: pulumi.Bool(true),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func toPulumiStringArray(a []string) pulumi.StringArrayInput {
	var res []pulumi.StringInput
	for _, s := range a {
		res = append(res, pulumi.String(s))
	}
	return pulumi.StringArray(res)
}
