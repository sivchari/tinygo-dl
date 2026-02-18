package version

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
)

// Run runs the specified version of TinyGo.
func Run(version string) {
	log.SetFlags(0)

	root, err := tinygoRoot(version)
	if err != nil {
		log.Fatalf("%s: %v", version, err)
	}

	if len(os.Args) > 1 && os.Args[1] == "download" {
		if len(os.Args) != 2 {
			log.Fatalf("usage: %s download", version)
		}
		if err := download(version, root); err != nil {
			log.Fatalf("%s: download failed: %v", version, err)
		}
		log.Printf("Success. You may now run '%s'!", version)
		os.Exit(0)
	}

	tinygobin := filepath.Join(root, "bin", "tinygo"+exe())
	if _, err := os.Stat(tinygobin); err != nil {
		log.Fatalf("%s: not downloaded. Run '%s download' to install to %v", version, version, root)
	}

	runTinyGo(root)
}

func tinygoRoot(version string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	return filepath.Join(home, "sdk", version), nil
}

func exe() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func runTinyGo(root string) {
	tinygobin := filepath.Join(root, "bin", "tinygo"+exe())
	cmd := exec.Command(tinygobin, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	// Ignore signals so child process can handle them
	signal.Ignore(os.Interrupt)

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func download(version, root string) error {
	// Parse version: tinygo0.40.1 -> 0.40.1
	ver := strings.TrimPrefix(version, "tinygo")

	// Determine OS and architecture
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Map to TinyGo release naming
	var archName string
	switch goarch {
	case "amd64":
		archName = "amd64"
	case "arm64":
		archName = "arm64"
	case "arm":
		archName = "arm"
	default:
		return fmt.Errorf("unsupported architecture: %s", goarch)
	}

	var ext string
	switch goos {
	case "darwin", "linux":
		ext = "tar.gz"
	case "windows":
		ext = "zip"
	default:
		return fmt.Errorf("unsupported OS: %s", goos)
	}

	// Build download URL
	// https://github.com/tinygo-org/tinygo/releases/download/v0.40.1/tinygo0.40.1.darwin-arm64.tar.gz
	url := fmt.Sprintf("https://github.com/tinygo-org/tinygo/releases/download/v%s/tinygo%s.%s-%s.%s",
		ver, ver, goos, archName, ext)

	log.Printf("Downloading %s...", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(root), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Extract archive
	if ext == "tar.gz" {
		if err := extractTarGz(resp.Body, root); err != nil {
			return fmt.Errorf("failed to extract: %v", err)
		}
	} else {
		return fmt.Errorf("zip extraction not implemented yet")
	}

	log.Printf("Installed to %s", root)
	return nil
}

func extractTarGz(r io.Reader, dest string) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// TinyGo archives have "tinygo/" prefix, strip it
		name, ok := strings.CutPrefix(header.Name, "tinygo/")
		if !ok {
			continue
		}

		if name == "" {
			continue
		}

		target := filepath.Join(dest, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		}
	}

	return nil
}
