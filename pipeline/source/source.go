package source

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goreleaser/goreleaser/context"
)

// Pipe for archive
type Pipe struct{}

// Description of the pipe
func (Pipe) Description() string {
	return "Creating source archives"
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
	var file = filepath.Join(ctx.Config.Dist, name+".tar.gz")
	fmt.Printf(" -> Building %s\n", file)
	args := make([]string, 0, len(ctx.Config.Source.Excludes)+4)
	for _, ex := range ctx.Config.Source.Excludes {
		args = append(args, "--exclude="+ex)
	}
	args = append(args, fmt.Sprintf("--transform=s|/|/%s-%s/|", ctx.Config.Build.Binary, ctx.Version))
	args = append(args, "-Pczf")
	args = append(args, file)
	args = append(args, ".")
	cmd := exec.Command("tar", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Source step failed: %s %+v\n%v", cmd.Path, cmd.Args, string(out))
	}
	ctx.AddArtifact(file)
	return nil
}
