package source

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goreleaser/goreleaser/context"
	"github.com/goreleaser/goreleaser/internal/artifact"
)

const SourceNameTemplate = "{{.Binary}}-{{.Version}}"

// Pipe for archive
type Pipe struct{}

// Description of the pipe
func (Pipe) String() string {
	return "Creating source archives"
}

func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Source.NameTemplate == "" {
		ctx.Config.Source.NameTemplate = SourceNameTemplate
	}
	return nil
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	return create(ctx)
}

// Archive represents a compression archive files from disk can be written to.
type Archive interface {
	Close() error
	Add(name, path string) error
}

func create(ctx *context.Context) error {
	if err := os.MkdirAll(ctx.Config.Dist, 0755); err != nil {
		return err
	}
	name, err := nameFor(ctx)
	if err != nil {
		return err
	}
	filename := name + ".tar.gz"
	path := filepath.Join(ctx.Config.Dist, filename)
	fmt.Printf(" -> Building %s\n", filename)
	if err := ioutil.WriteFile("COMMIT", []byte(ctx.Git.Commit), 0644); err != nil {
		return err
	}
	args := make([]string, 0, len(ctx.Config.Source.Excludes)+4)
	for _, ex := range ctx.Config.Source.Excludes {
		args = append(args, "--exclude="+ex)
	}
	args = append(args, fmt.Sprintf("--transform=s|/|/%s-%s/|", ctx.Config.Builds[0].Binary, ctx.Version))
	args = append(args, "-Pczf")
	args = append(args, path)
	args = append(args, ".")
	cmd := exec.Command("tar", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Source step failed: %s %+v\n%v", cmd.Path, cmd.Args, string(out))
	}
	ctx.Artifacts.Add(artifact.Artifact{
		Type: artifact.Source,
		Name: filename,
		Path: path,
	})
	ctx.Config.Brew.SourceTarball = filename
	os.Remove("COMMIT")
	return nil
}
