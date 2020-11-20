package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type file struct {
	name       string
	currentPWD string
	filePWD    string
	infoPWD    string
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getPWD() string {
	path, err := os.Getwd()
	exitOnError(err)

	return path
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// TODO(#2): add a restore function
// func restore() {

// }

// TODO(#3): handle situations where the file already exists
func delete(current []file) {
	for _, file := range current {
		err := os.Rename(file.currentPWD, file.filePWD)
		exitOnError(err)

		f, err := os.Create(file.infoPWD)
		check(err)
		defer f.Close()

		currentTime := time.Now()
		body := []byte("[Trash Info]\n" + "Path=" + file.currentPWD + "\n" + "DeletionDate=" + currentTime.Format(time.RFC3339))

		_, err = f.Write(body)
		check(err)
	}
}

func genPaths(parent string, files []string) []file {
	trash := filepath.Join(os.Getenv("XDG_DATA_HOME"), "Trash/Trash")
	fFolder := filepath.Join(trash, "files")
	iFolder := filepath.Join(trash, "info")

	FILES := make([]file, len(files))
	for idx, afile := range files {
		FILES[idx].name = afile
		FILES[idx].currentPWD = filepath.Join(parent, afile)
		FILES[idx].filePWD = filepath.Join(fFolder, afile)
		FILES[idx].infoPWD = filepath.Join(iFolder, afile+".info")
	}

	return FILES
}

// TODO(#1): add flags in order to parse actions
func main() {
	file := genPaths(getPWD(), os.Args[1:])
	delete(file)
	fmt.Println(file)

	// fmt.Println(getPWD())
}
