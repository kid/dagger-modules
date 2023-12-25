package main

import (
	"context"
	"fmt"
	"path/filepath"
)

const (
	OUT_DIR   = "/out"
	PLAN_FILE = "apply.tfplan"
)

type Terraform struct {
	Version   string
	Directory string
}

func New(directory Optional[string], version Optional[string]) *Terraform {
	return &Terraform{
		Directory: version.GetOr("fixtures"),
		Version:   version.GetOr("1.6.6"),
	}
}

// example usage: "dagger call plan --directory stack"
func (m *Terraform) Plan(ctx context.Context) *Directory {
	exec := m.Base().
		WithExec([]string{"plan", "-input=false", "-out", filepath.Join(OUT_DIR, PLAN_FILE)})

	output, err := exec.Stdout(ctx)
	if err != nil {
		panic(err)
	}

	return exec.
		WithNewFile(filepath.Join(OUT_DIR, "apply.txt"), ContainerWithNewFileOpts{
			Contents: output,
		}).
		Directory(OUT_DIR)
}

func (m *Terraform) PlanOutput(ctx context.Context) *File {
	return m.Plan(ctx).File("apply.txt")
}

// example usage: "dagger call apply --directory stack"
func (m *Terraform) Apply(ctx context.Context, plan *File) *Container {
	return m.Base().
		WithFile(PLAN_FILE, plan).
		WithExec([]string{"apply", "apply.tfplan"})
}

func (m *Terraform) Base() *Container {
	return dag.Container().
		From(fmt.Sprintf("docker.io/hashicorp/terraform:%s", m.Version)).
		WithDirectory(OUT_DIR, dag.Directory()).
		WithMountedDirectory("/src", dag.Host().Directory(m.Directory)).
		WithWorkdir("/src")
}
