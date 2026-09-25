// The tinygo0.31.2 command runs TinyGo 0.31.2.
//
// To install, run:
//
//	go install github.com/sivchari/tinygo-dl/tinygo0.31.2@latest
//	tinygo0.31.2 download
//
// And then use the tinygo0.31.2 command as if it were your normal tinygo command.
package main

import "github.com/sivchari/tinygo-dl/internal/version"

func main() {
	version.Run("tinygo0.31.2")
}
