package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const sourcePath = "/Volumes/Media/Emulation/ps1"
const debug = true

func unzip(source string, dest string) error {
	cmd := exec.Command("/usr/bin/unzip", source, "-d", dest)
	log.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func compressToCHD(source string, dest string) error {
	cmd := exec.Command(
		"chdman",
		"createcd",
		"-i",
		source,
		"-o",
		dest,
	)
	log.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {

	files, err := os.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	for idx, f := range files {
		filename := f.Name()

		if !f.Type().IsRegular() ||
			!strings.HasSuffix(filename, ".zip") {
			// ignore non zip files
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
