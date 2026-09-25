// The tinygo0.5.0 command runs TinyGo 0.5.0.
//
// To install, run:
//
//	go install github.com/sivchari/tinygo-dl/tinygo0.5.0@latest
//	tinygo0.5.0 download
//
// And then use the tinygo0.5.0 command as if it were your normal tinygo command.
package main

import "github.com/sivchari/tinygo-dl/internal/version"

func main() {
	version.Run("tinygo0.5.0")
}
