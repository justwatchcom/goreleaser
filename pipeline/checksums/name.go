package checksums

import (
	"bytes"
	"text/template"

	"github.com/goreleaser/goreleaser/context"
)

type nameData struct {
	Os      string
	Arch    string
	Arm     string
	Version string
	Tag     string
	Binary  string
}

func nameFor(ctx *context.Context) (string, error) {
	var data = nameData{
		Version: ctx.Version,
		Tag:     ctx.Git.CurrentTag,
		Binary:  ctx.Config.Build.Binary,
	}

	var out bytes.Buffer
	t, err := template.New(data.Binary).Parse(ctx.Config.Checksum.NameTemplate)
	if err != nil {
		return "", err
	}
	err = t.Execute(&out, data)
	return out.String(), err
}
