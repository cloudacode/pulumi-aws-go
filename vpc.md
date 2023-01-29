# VPC Quickstart

In this example, you will build a pulumi project, which can provision VPC resources on AWS.

## Inline Program

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

## Pulumi Over HTTP

```go

```
