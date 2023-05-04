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

func main() {
	files, err := os.ReadDir("/Volumes/Media/Emulation/ps1")
	if err != nil {
		log.Fatal(err)
	}
	for idx, f := range files {
		if idx > 3 {
			fmt.Println("done")
			os.Exit(0)
		}
		filename := f.Name()
		if !strings.HasSuffix(filename, ".zip") {
			continue
		}
		zipFilepath := filepath.Clean(fmt.Sprintf("%s/%s", "/Volumes/Media/Emulation/ps1", filename))
		rootfilename := strings.Split(zipFilepath, ".")[0]
		fmt.Printf("%s --> %s\n", filename, rootfilename)

		tempdir, err := os.MkdirTemp("", "go-chd-")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Made %s\n", tempdir)
		defer os.RemoveAll(tempdir)
		cmd := exec.Command("/usr/bin/unzip", zipFilepath, "-d", tempdir)
		log.Println(cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		zipdir, err := os.ReadDir(tempdir)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Found %d", len(zipdir))
		destname := path.Base(rootfilename)
		os.MkdirAll(destname, 0755)
		for _, zf := range zipdir {
			if !strings.HasSuffix(zf.Name(), ".cue") {
				continue
			}
			log.Printf("found %s\n", zf.Name())
			cmd := exec.Command(
				"chdman",
				"createcd",
				"-i",
				fmt.Sprintf("%s/%s", tempdir, zf.Name()),
				"-o",
				fmt.Sprintf("%s/%s.chd", destname, destname),
			)
			log.Println(cmd)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}
