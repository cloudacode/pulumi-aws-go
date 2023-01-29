package aws

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecs"
	awsxecs "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ecs"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A collection of values returned by FargateRun.
type GetFargetRunResult struct {
	ServiceName pulumi.StringOutput
	Url         pulumi.StringOutput
}

// New Fargate registers a new resource with the given vpc id, unique name, image url, and exposed port.
// https://www.pulumi.com/docs/guides/crosswalk/aws/ecs/#creating-an-ecs-cluster-in-a-vpc
func FargateRun(ctx *pulumi.Context, vpcId, prefixName, imageUrl string, containerPort int) (*GetFargetRunResult, error) {

	var rv GetFargetRunResult

	// Lookup the VPC information
	vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{Id: &vpcId})
	if err != nil {
		return nil, err
	}

	// Get Subents from the VPC
	subnet, err := ec2.GetSubnets(ctx, &ec2.GetSubnetsArgs{Filters: []ec2.GetSubnetsFilter{
		{
			Name:   "vpc-id",
			Values: []string{vpcId},
		},
	}})
	if err != nil {
		return nil, err
	}

	// Set a new security groups
	securityGroup, err := ec2.NewSecurityGroup(ctx, prefixName+"-sg", &ec2.SecurityGroupArgs{
		VpcId: pulumi.String(vpc.Id),
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol: pulumi.String("tcp"),
				FromPort: pulumi.Int(containerPort),
				ToPort:   pulumi.Int(containerPort),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0")},
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			&ec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				Protocol: pulumi.String("-1"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a new ECS cluster
	cluster, err := ecs.NewCluster(ctx, prefixName+"-cluster", nil)
	if err != nil {
		return nil, err
	}

	// Create a LB for the cluster
	lb, err := lb.NewApplicationLoadBalancer(ctx, "lb", &lb.ApplicationLoadBalancerArgs{
		DefaultTargetGroupPort: pulumi.Int(containerPort),
	})
	if err != nil {
		return nil, err
	}

	// Create a Fargate on the ECS Cluster
	service, err := awsxecs.NewFargateService(ctx, prefixName+"-service", &awsxecs.FargateServiceArgs{
		Cluster: cluster.Arn,
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			AssignPublicIp: pulumi.Bool(true),
			Subnets:        toPulumiStringArray(subnet.Ids),
			SecurityGroups: pulumi.StringArray{
				securityGroup.ID(),
			},
		},
		DesiredCount: pulumi.Int(1),
		TaskDefinitionArgs: &awsxecs.FargateServiceTaskDefinitionArgs{
			Container: &awsxecs.TaskDefinitionContainerDefinitionArgs{
				Image:     pulumi.String(imageUrl),
				Cpu:       pulumi.Int(256),
				Memory:    pulumi.Int(512),
				Essential: pulumi.Bool(true),
				PortMappings: &awsxecs.TaskDefinitionPortMappingArray{
					&awsxecs.TaskDefinitionPortMappingArgs{
						ContainerPort: pulumi.Int(containerPort),
						HostPort:      pulumi.Int(containerPort),
						TargetGroup:   lb.DefaultTargetGroup,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Return values
	rv.ServiceName = service.Service.Name()
	rv.Url = lb.LoadBalancer.DnsName()

	return &rv, nil
}

// Format the StringArray from Strings
func toPulumiStringArray(a []string) pulumi.StringArrayInput {
	var res []pulumi.StringInput
	for _, s := range a {
		res = append(res, pulumi.String(s))
	}
	return pulumi.StringArray(res)
}
