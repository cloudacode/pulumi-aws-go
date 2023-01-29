# Pulumi AWS Go Package by Cloudacode
Amazon Web Services (AWS) Pulumi resource package provided by cloudacode

[![Go Report Card](https://goreportcard.com/badge/github.com/cloudacode/pulumi-aws-go)](https://goreportcard.com/badge/github.com/cloudacode/pulumi-aws-go)
[![GoDoc](https://godoc.org/github.com/cloudacode/pulumi-aws-go?status.svg)](https://godoc.org/github.com/cloudacode/pulumi-aws-go)

## How to use it

### Prerequisites

- [Golang](https://golang.org/dl/) version 1.18 or above. You can follow the instructions in the official [installation page](https://golang.org/doc/install) (check it by `go version`)
- [Pulumi](https://www.pulumi.com/). You can fllow the instructions in the official [installation page](https://www.pulumi.com/docs/get-started/install/)

### Create a Pulumi Project

Create and move to directory which is your project workspace
```bash
mkdir pulumi-aws-project && cd pulumi-aws-project
```

Please login pulumi backend to store your state
```bash
pulumi login
```

### Import Package

First, you can to use `go get` to fetch the `latest` package or any git tag version. The git tags are available in [release page](https://github.com/cloudacode/pulumi-aws-go/releases). Take `v0.4.0` as an example:

```bash
# fetch the latest version
go get github.com/cloudacode/pulumi-aws-go
# or explicit target version
go get github.com/cloudacode/pulumi-aws-go@v0.4.0
```

Now, open `main.go` and start coding!

## Example code

### VPC Provisioning with Dynamic IP allocation

To automate VPC provisioning via IPAM, you need to:
1. Use the [`aws.VpcRunIpam`](https://pkg.go.dev/github.com/cloudacode/pulumi-aws-go/aws#VpcRunIpam) function to provision VPC resource
2. Set the arguments as [`ctx *pulumi.Context, prefixName, ipamID, ipamPoolId, netMaskLength`](https://github.com/cloudacode/pulumi-aws-go/blob/main/aws/vpc.go#L9) on the function.

Please check some examples in the [getting started vpc](./vpc.md).

### ECS Fargete Deployment

To deploy ECS Fargate via Pulumi, you need to:
1. Use the [`aws.FargateRun`](https://pkg.go.dev/github.com/cloudacode/pulumi-aws-go/aws#FargateRun) function to deploy ECS Fargate resource
2. Set the arguments as [`ctx *pulumi.Context, vpcId, prefixName, imageUrl string, containerPort int`](https://github.com/cloudacode/pulumi-aws-go/blob/main/aws/fargate.go#L11) on the function.

Please check some examples in the [getting started fargate](./fargate.md).

🌟🌟 Enjoy!!!
