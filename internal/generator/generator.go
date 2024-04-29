package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	_ "embed"
)

type ConfigImport struct {
	Package string `json:"package"`
}

type Config struct {
	ModuleName string         `json:"module_name"`
	GoVersion  string         `json:"go_version"`
	Imports    []ConfigImport `json:"imports"`
}

//go:embed templates/go.mod.template
var goModTemplate string

//go:embed templates/main.go.template
var mainGoTemplate string

func (c Config) GenerateInto(ctx context.Context, dir string) error {
	sort.Slice(c.Imports, func(i, j int) bool {
		return c.Imports[i].Package < c.Imports[j].Package
	})

	for k, v := range map[string]string{
		"main.go": mainGoTemplate,
		"go.mod":  goModTemplate,
	} {
		outFile, err := os.Create(filepath.Join(dir, k))
		if err != nil {
			return fmt.Errorf("failed to create %v: %w", k, err)
		}
		outTemplate, err := template.New(k).Parse(v)
		if err != nil {
			return fmt.Errorf("failed to initialise %v template: %w", k, err)
		}
		if err := outTemplate.Execute(outFile, c); err != nil {
			return fmt.Errorf("failed to execute %v template: %w", k, err)
		}
		if err := outFile.Close(); err != nil {
			return fmt.Errorf("failed to close %v file: %w", k, err)
		}
	}
	return nil
}
