package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

const sourcePath = "/Volumes/Media/Emulation/ps1"
const debug = true

// execCmd wraps exec.Command with reasonable defaults and logging.
// takes same arguments as exec.Command.
func execCmd(executable string, args ...string) error {
	cmd := exec.Command(executable, args...)
	log.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// unpack the indicated zip file into the specified directory
// using the unzip command.
func unzip(source string, dest string) error {
	return execCmd("unzip", source, "-d", dest)
}

// compress the unpacked ISO into a CHD file.
func compressToCHD(source string, dest string) error {
	return execCmd(
		"chdman",
		"createcd",
		"-i",
		source,
		"-o",
		dest,
	)
}

func convert(c *cli.Context) {
	// read sourcePath directory
	files, err := os.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	// search for zip files
	for idx, f := range files {
		filename := f.Name()

		if !f.Type().IsRegular() ||
			!strings.HasSuffix(filename, ".zip") {
			// skip non zip files
			continue
		}

		if debug && idx == 1 {
			os.Exit(0)
		}

		zipFilepath := filepath.Clean(fmt.Sprintf("%s/%s", sourcePath, filename))
		rootfilename := strings.Split(zipFilepath, ".")[0]

		fmt.Printf("Converting iso(%s) --> chd(%s)\n", filename, rootfilename)

		tempDir, err := os.MkdirTemp("", "go-chd-")
		if err != nil {
			log.Fatal(err)
		}

		// clean up after ourselves
		defer os.RemoveAll(tempDir)

		log.Printf("Extracting to %s\n", tempDir)

		err = unzip(zipFilepath, tempDir)
		if err != nil {
			log.Fatal(err)
		}

		// read newly unpacked directory
		zipdir, err := os.ReadDir(tempDir)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Found %d files", len(zipdir))

		destName := path.Base(rootfilename)
		os.MkdirAll(destName, 0755)

		for _, zf := range zipdir {

			if !zf.Type().IsRegular() ||
				!strings.HasSuffix(zf.Name(), ".cue") {
				// skip non-cue files
				continue
			}

			log.Printf("Found cue file %s\n", zf.Name())
			src := fmt.Sprintf("%s/%s", tempDir, zf.Name())
			dest := fmt.Sprintf("%s/%s.chd", destName, destName)
			compressToCHD(src, dest)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

func main() {
	app := cli.NewApp()
	app.Name = "convert"
	app.Usage = "convert zip compressed iso game files to CHD"
	app.Action = convert
	app.Run(os.Args)
}
