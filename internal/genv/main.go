// The genv command generates a tinygo<version> wrapper directory.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

var versionRE = regexp.MustCompile(`^\d+(\.\d+)+$`)

var tmpl = template.Must(template.New("main.go").Parse(`// The tinygo{{.}} command runs TinyGo {{.}}.
//
// To install, run:
//
//	go install github.com/sivchari/tinygo-dl/tinygo{{.}}@latest
//	tinygo{{.}} download
//
// And then use the tinygo{{.}} command as if it were your normal tinygo command.
package main

import "github.com/sivchari/tinygo-dl/internal/version"

func main() {
	version.Run("tinygo{{.}}")
}
`))

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: genv <version> [<version> ...]")
		os.Exit(1)
	}
	for _, v := range os.Args[1:] {
		if !versionRE.MatchString(v) {
			fmt.Fprintf(os.Stderr, "genv: invalid version %q: must match %s\n", v, versionRE)
			os.Exit(1)
		}
		dir := "tinygo" + v
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "genv: %v\n", err)
			os.Exit(1)
		}
		f, err := os.Create(filepath.Join(dir, "main.go"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "genv: %v\n", err)
			os.Exit(1)
		}
		if err := tmpl.Execute(f, v); err != nil {
			f.Close()
			fmt.Fprintf(os.Stderr, "genv: %v\n", err)
			os.Exit(1)
		}
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "genv: %v\n", err)
			os.Exit(1)
		}
	}
}
