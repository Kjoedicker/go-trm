package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type file struct {
	name       string
	alias      string
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

func exists(name string, path string, newName string, count int) string {
	if _, err := os.Stat(filepath.Join(path, newName)); err != nil {
		if os.IsNotExist(err) {
			return newName
		}
	}
	return exists(name, path, fmt.Sprintf("%v.%v", name, count), count+1)
}

func genPaths(parent string, files []string) []file {
	trash := filepath.Join(os.Getenv("XDG_DATA_HOME"), "Trash/Trash")
	fFolder := filepath.Join(trash, "files")
	iFolder := filepath.Join(trash, "info")

	FILES := make([]file, len(files))
	for idx, afile := range files {

		FILES[idx].name = afile
		FILES[idx].alias = exists(afile, fFolder, afile, 0)
		FILES[idx].currentPWD = filepath.Join(parent, FILES[idx].name)
		FILES[idx].filePWD = filepath.Join(fFolder, FILES[idx].alias)
		FILES[idx].infoPWD = filepath.Join(iFolder, FILES[idx].alias+".info")
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
