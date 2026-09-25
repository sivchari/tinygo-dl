// The tinygo0.7.1 command runs TinyGo 0.7.1.
//
// To install, run:
//
//	go install github.com/sivchari/tinygo-dl/tinygo0.7.1@latest
//	tinygo0.7.1 download
//
// And then use the tinygo0.7.1 command as if it were your normal tinygo command.
package main

import "github.com/sivchari/tinygo-dl/internal/version"

func main() {
	version.Run("tinygo0.7.1")
}
