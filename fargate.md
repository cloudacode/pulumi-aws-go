# Fargate Quickstart

In this example, you will build a pulumi project, which can provision ECS Fargate on AWS.

## Inline Program

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

## Pulumi Over HTTP

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudacode/pulumi-aws-go/aws"
	"github.com/gorilla/mux"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// define request/response types for various REST ops

type CreateSiteReq struct {
	ID string `json:"id"`
	Image string `json:"image"`
	Port  int    `json:"port"`
}

type UpdateSiteReq struct {
	Image string `json:"image"`
	Port  int    `json:"port"`
}

type SiteResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type ListSitesResponse struct {
	IDs []string `json:"ids"`
}

var project = "fargate-test"

func main() {
	ensurePlugins()

	router := mux.NewRouter()

	// setup our RESTful routes for our Site resource
	router.HandleFunc("/fargate", createHandler).Methods("POST")
	router.HandleFunc("/fargate", listHandler).Methods("GET")
	router.HandleFunc("/fargate/{id}", getHandler).Methods("GET")
	router.HandleFunc("/fargate/{id}", updateHandler).Methods("PUT")
	router.HandleFunc("/fargate/{id}", deleteHandler).Methods("DELETE")

	// define and start our http server
	s := &http.Server{
		Addr:    ":8088",
		Handler: router,
	}
	fmt.Println("starting server on :8088")
	log.Fatal(s.ListenAndServe())
}

// creates new sites
func createHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var createReq CreateSiteReq
	err := json.NewDecoder(req.Body).Decode(&createReq)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "failed to parse create request")
		return
	}

	ctx := context.Background()

	stackName := createReq.ID
	program := createPulumiProgram(createReq.Image, createReq.Port)

	s, err := auto.NewStackInlineSource(ctx, stackName, project, program)
	if err != nil {
		// if stack already exists, 409
		if auto.IsCreateStack409Error(err) {
			w.WriteHeader(409)
			fmt.Fprintf(w, fmt.Sprintf("stack %q already exists", stackName))
			return
		}

		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "eu-north-1"})

	// deploy the stack
	// we'll write all of the update logs to st	out so we can watch requests get processed
	upRes, err := s.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	response := &SiteResponse{
		ID:  stackName,
		URL: upRes.Outputs["ecsUrl"].Value.(string),
	}
	json.NewEncoder(w).Encode(&response)
}

// lists all sites
func listHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()
	// set up a workspace with only enough information for the list stack operations
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:    tokens.PackageName(project),
		Runtime: workspace.NewProjectRuntimeInfo("go", nil),
	}))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	stacks, err := ws.ListStacks(ctx)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	var ids []string
	for _, stack := range stacks {
		ids = append(ids, stack.Name)
	}

	response := &ListSitesResponse{
		IDs: ids,
	}
	json.NewEncoder(w).Encode(&response)
}

// gets info about a specific site
func getHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	stackName := params["id"]
	// we don't need a program since we're just getting stack outputs
	var program pulumi.RunFunc = nil
	ctx := context.Background()
	s, err := auto.SelectStackInlineSource(ctx, stackName, project, program)
	if err != nil {
		// if the stack doesn't already exist, 404
		if auto.IsSelectStack404Error(err) {
			w.WriteHeader(404)
			fmt.Fprintf(w, fmt.Sprintf("stack %q not found", stackName))
			return
		}

		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	// fetch the outputs from the stack
	outs, err := s.Outputs(ctx)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	response := &SiteResponse{
		ID:  stackName,
		URL: outs["ecsURL"].Value.(string),
	}
	json.NewEncoder(w).Encode(&response)
}

// updates the content for an existing site
func updateHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updateReq UpdateSiteReq
	err := json.NewDecoder(req.Body).Decode(&updateReq)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "failed to parse create request")
		return
	}

	ctx := context.Background()
	params := mux.Vars(req)
	stackName := params["id"]
	program := createPulumiProgram(updateReq.Image, updateReq.Port)

	s, err := auto.SelectStackInlineSource(ctx, stackName, project, program)
	if err != nil {
		if auto.IsSelectStack404Error(err) {
			w.WriteHeader(404)
			fmt.Fprintf(w, fmt.Sprintf("stack %q not found", stackName))
			return
		}

		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "eu-north-1"})

	// deploy the stack
	// we'll write all of the update logs to st	out so we can watch requests get processed
	upRes, err := s.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		// if we already have another update in progress, return a 409
		if auto.IsConcurrentUpdateError(err) {
			w.WriteHeader(409)
			fmt.Fprintf(w, "stack %q already has update in progress", stackName)
			return
		}

		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	response := &SiteResponse{
		ID:  stackName,
		URL: upRes.Outputs["ecsURL"].Value.(string),
	}
	json.NewEncoder(w).Encode(&response)
}

// deletes a site
func deleteHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := context.Background()
	params := mux.Vars(req)
	stackName := params["id"]
	// program doesn't matter for destroying a stack
	program := createPulumiProgram("", 0)

	s, err := auto.SelectStackInlineSource(ctx, stackName, project, program)
	if err != nil {
		// if stack doesn't already exist, 404
		if auto.IsSelectStack404Error(err) {
			w.WriteHeader(404)
			fmt.Fprintf(w, fmt.Sprintf("stack %q not found", stackName))
			return
		}

		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: "eu-north-1"})

	// destroy the stack
	// we'll write all of the logs to stdout so we can watch requests get processed
	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	// delete the stack and all associated history and config
	err = s.Workspace().RemoveStack(ctx, stackName)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(200)
}

// this function defines our ECS Fargate resources in terms of the contents that the caller passes in.
// this allows us to dynamically deploy ECS Fargate based on user defined values from the POST body.
func createPulumiProgram(image string, port int) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {

		res, _ := aws.FargateRun(ctx, "vpc-948b7cfd", "test", image, port)

		// export the ecs URL
		ctx.Export("ecsUrl", res.Url)
		return nil
	}
}

// ensure plugins runs once before the server boots up
// making sure the proper pulumi plugins are installed
func ensurePlugins() {
	ctx := context.Background()
	w, err := auto.NewLocalWorkspace(ctx)
	if err != nil {
		fmt.Printf("Failed to setup and run http server: %v\n", err)
		os.Exit(1)
	}
	err = w.InstallPlugin(ctx, "aws", "v4.0.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		os.Exit(1)
	}
}
```
