// Command embed-binaries is used to build and embed Go binaries in Go code.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	buildFrom = flag.String("build-from", "", "path to the directory hosting the main package to build")
	embedTo   = flag.String("embed-to", ".", "the directory to which embedded binaries should be written")
	prefix    = flag.String("prefix", "", "the prefix to use for the embedded binary files")
	pkg       = flag.String("pkg", "main", "package name to use in generated code")
	assetName = flag.String("asset-name", "asset", "the name of the asset in the generated code")
)

var targets = []struct {
	GOOS, GOARCH, suffix string
}{
	{"darwin", "amd64", "darwin"},
	{"linux", "386", "linux_386"},
	{"linux", "amd64", "linux_amd64"},
	{"windows", "amd64", "windows"},
}

func fail(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *buildFrom == "" {
		fail("build-from flag is required")
	}

	tmpDir, err := ioutil.TempDir("", "lantern-diagnostics-binaries")
	if err != nil {
		fail("failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tmpDir)

	for _, target := range targets {
		fmt.Fprintf(os.Stderr, "Building for %s/%s\n", target.GOOS, target.GOARCH)

		binaryPath := filepath.Join(tmpDir, *assetName)
		buildBinary := exec.Command("go", "build", "-o", binaryPath, *buildFrom)
		buildBinary.Env = os.Environ()
		// The last values for each key in the slice are used, so these will override any existing
		// settings for GOOS and GOARCH.
		buildBinary.Env = append(
			buildBinary.Env,
			fmt.Sprintf("%s=%s", "GOOS", target.GOOS),
			fmt.Sprintf("%s=%s", "GOARCH", target.GOARCH),
		)
		buildBinary.Stdout, buildBinary.Stderr = os.Stderr, os.Stderr
		if err := buildBinary.Run(); err != nil {
			fail("failed to build binaries:", err)
		}

		embeddedFileName := target.suffix + ".go"
		if *prefix != "" {
			embeddedFileName = fmt.Sprintf("%s_%s", *prefix, embeddedFileName)
		}
		embedBinary := exec.Command(
			"go-bindata",
			"-o", filepath.Join(*embedTo, embeddedFileName),
			"-pkg", *pkg,
			"-prefix", tmpDir,
			tmpDir,
		)
		embedBinary.Stdout, embedBinary.Stderr = os.Stderr, os.Stderr
		if err := embedBinary.Run(); err != nil {
			fail("failed to embed binaries:", err)
		}
	}
}
