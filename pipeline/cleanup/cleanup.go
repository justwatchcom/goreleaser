package cleanup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goreleaser/goreleaser/context"
)

type Pipe struct{}

func (Pipe) Description() string {
	return "Prepare release"
}

func (Pipe) Run(ctx *context.Context) error {
	for _, step := range ctx.Config.Cleanup.Hooks {
		args := strings.Fields(step)
		cmd := exec.Command(args[0], args[1:]...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("hook failed: %s\n%v", step, string(out))
		}
	}
	dir, err := filepath.Abs(ctx.Config.Dist)
	if err != nil {
		return err
	}
	if dir == "" || dir == "/" {
		return fmt.Errorf("Sneaky dir: %s", dir)
	}
	return os.RemoveAll(dir)
}
