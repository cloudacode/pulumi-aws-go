# pulumi-aws-go
Amazon Web Services (AWS) Pulumi resource package provided by cloudacode

## Initalize

in module workspace
```bash
go mod init github.com/cloudacode/pulumi-aws-go/fargate
go mod tidy
```

in root directory
```bash
go work init
go work use . ./fargate ./vpc
```

## Import local module
to import the local module code
```bash
go mod edit -replace github.com/cloudacode/pulumi-aws-go/fargate=../pulumi-aws-go/fargate
```


# Reference
https://www.digitalocean.com/community/tutorials/how-to-distribute-go-modules

https://go.dev/doc/tutorial/workspaces
!!!INFO
    Go 1.18 or later
