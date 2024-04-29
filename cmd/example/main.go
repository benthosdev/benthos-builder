package main

import (
	"context"
	"log/slog"

	"github.com/benthosdev/benthos-builder/internal/command"
	"github.com/benthosdev/benthos-builder/internal/generator"
)

func main() {
	conf := generator.Config{
		ModuleName: "foo.com/meowdev/woofer",
		GoVersion:  "1.22.2",
		Imports: []generator.ConfigImport{
			{Package: "github.com/benthosdev/benthos/v4/public/components/nanomsg"},
			{Package: "github.com/benthosdev/benthos/v4/public/components/kafka"},
			{Package: "github.com/benthosdev/benthos/v4/public/components/memcached"},
			{Package: "github.com/benthosdev/benthos/v4/public/components/mongodb"},
			{Package: "github.com/benthosdev/benthos/v4/public/components/mqtt"},
			{Package: "github.com/benthosdev/benthos/v4/public/components/maxmind"},
			{Package: "github.com/benthosdev/benthos/v4/public/components/msgpack"},
			{Package: "github.com/benthosdev/benthos/v4/public/components/nats"},
		},
	}

	dir := "./tmp"
	ctx := context.Background()

	slog.Info("Generating module files")
	if err := conf.GenerateInto(ctx, dir); err != nil {
		panic(err)
	}

	c := command.InDir(dir)

	slog.Info("Pulling in module imports")
	if err := c.GoModTidy(ctx); err != nil {
		panic(err)
	}

	slog.Info("Building benthos")
	if err := c.GoBuild(ctx); err != nil {
		panic(err)
	}
}
