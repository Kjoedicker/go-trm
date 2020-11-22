package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type logic struct {
	trash      string
	fileFolder string
	infoFolder string
}

func genLogic() logic {
	trash := filepath.Join(os.Getenv("XDG_DATA_HOME"), "Trash/Trash")

	return logic{
		fileFolder: filepath.Join(trash, "files"),
		infoFolder: filepath.Join(trash, "info"),
	}
}

type file struct {
	name       string
	alias      string
	currentPWD string
	filePWD    string
	infoPWD    string
}

func genPaths(parent string, files []string, alias bool) []file {
	logistics := genLogic()

	FILES := make([]file, len(files))
	for idx, afile := range files {
		FILES[idx].name = afile
		FILES[idx].currentPWD = filepath.Join(parent, FILES[idx].name)

		switch alias {
		case true:
			FILES[idx].alias = rename(afile, logistics.fileFolder, afile, 0)
		default:
			FILES[idx].alias = afile
		}

		FILES[idx].filePWD = filepath.Join(logistics.fileFolder, FILES[idx].alias)
		FILES[idx].infoPWD = filepath.Join(logistics.infoFolder, FILES[idx].alias+".info")
	}

	return FILES
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

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func rename(name string, path string, newName string, count int) string {
	fmt.Println(newName)
	if _, err := os.Stat(filepath.Join(path, newName)); err != nil {
		if os.IsNotExist(err) {
			return newName
		}
	}
	return rename(name, path, fmt.Sprintf("%v.%v", name, count), count+1)
}

func move(fro string, to string) {
	fmt.Println(to)
	err := os.Rename(fro, to)
	exitOnError(err)
}

// TODO(#2): add a restore function
func restore(current []file) {
	for _, file := range current {
		move(file.filePWD, file.currentPWD)
		os.Remove(file.infoPWD)
	}
}

func delete(current []file) {
	for _, file := range current {
		move(file.currentPWD, file.filePWD)

		f, err := os.Create(file.infoPWD)
		check(err)
		defer f.Close()

		currentTime := time.Now()
		body := []byte("[Trash Info]\n" + "Path=" + file.currentPWD + "\n" + "DeletionDate=" + currentTime.Format(time.RFC3339))

		_, err = f.Write(body)
		check(err)
	}
}

// TODO(#1): add flags in order to parse actions
func main() {
	file := genPaths(getPWD(), os.Args[1:], false)
	// delete(file)
	restore(file)

	// fmt.Println(file)

	// fmt.Println(getPWD())
}
