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

First, you can to use `go get` to fetch the package through git tag version. The git tags are available in [release page](https://github.com/cloudacode/pulumi-aws-go/releases). Take `v0.4.0` as an example:

```bash
go get github.com/cloudacode/pulumi-aws-go@v0.4.0
```

Now, open `main.go` and start coding!

## Example code

### VPC Provisioning with Dynamic IP allocation

To automate VPC provisioning via IPAM, you need to:
1. Use the [`aws.VpcRunIpam`](https://pkg.go.dev/github.com/cloudacode/pulumi-aws-go/aws#VpcRunIpam) function to provision VPC resource
2. Set the arguments as [`ctx *pulumi.Context, prefixName, ipamID, ipamPoolId, netMaskLength`](https://github.com/cloudacode/pulumi-aws-go/blob/main/aws/vpc.go#L9) on the function.

```go
package main

import (
	"github.com/cloudacode/pulumi-aws-go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := aws.VpcRunIpam(ctx, "test", "ipam-01f9c4c97064eb14b", "ipam-pool-08ecf378574aa542e", 18)
		if err != nil {
			return err
		}
		return nil
	})
}
```

### ECS Fargete Deployment

To deploy ECS Fargate via Pulumi, you need to:
1. Use the [`aws.FargateRun`](https://pkg.go.dev/github.com/cloudacode/pulumi-aws-go/aws#FargateRun) function to deploy ECS Fargate resource
2. Set the arguments as [`ctx *pulumi.Context, vpcId, prefixName string`](https://github.com/cloudacode/pulumi-aws-go/blob/main/aws/fargate.go#L11) on the function.

```go
package main

import (
	"github.com/cloudacode/pulumi-aws-go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := aws.FargateRun(ctx, "vpc-948b7cfd", "test")
		if err != nil {
			return err
		}
		return nil
	})
}
```

🌟🌟 Enjoy!!!
