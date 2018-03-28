package debug

import (
	"fmt"

	"github.com/goreleaser/goreleaser/context"
)

type Pipe struct{}

func (Pipe) String() string {
	return "Printing context"
}

func (Pipe) Default(ctx *context.Context) error {
	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	fmt.Printf("Context:\n---\n%+v\n---\n", ctx)
	return nil
}
