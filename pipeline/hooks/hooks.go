package hooks

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/goreleaser/goreleaser/context"
)

// Pipe for hooks
type Pipe struct{}

// Description of the pipe
func (Pipe) Description() string {
	return "Running hooks"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	return runHook(ctx.Config.Build.Env, ctx.Config.Build.Hooks.Post)
}

func runHook(env []string, hook string) error {
	if hook == "" {
		return nil
	}
	log.Println("Running hook", hook)
	cmd := strings.Fields(hook)
	return run(cmd, env)
}

func run(command, env []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, env...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("hook failed: %s\n%v", command, string(out))
	}
	return nil
}
