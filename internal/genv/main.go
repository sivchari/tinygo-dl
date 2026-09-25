// The genv command generates a tinygo<version> wrapper directory.
package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed main.go.tmpl
var mainGo string

var tmpl = template.Must(template.New("main.go").Parse(mainGo))

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: genv <version> [<version> ...]")
		os.Exit(1)
	}
	for _, v := range os.Args[1:] {
		if err := generate(v); err != nil {
			fmt.Fprintf(os.Stderr, "genv: %v\n", err)
			os.Exit(1)
		}
	}
}

func generate(version string) error {
	dir := "tinygo" + version
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(dir, "main.go"))
	if err != nil {
		return err
	}
	if err := tmpl.Execute(f, version); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
