package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type logic struct {
	trash      string
	fileFolder string
	infoFolder string
}

func exists(paths []string) {
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				err := os.MkdirAll(path, 0755)
				if err != nil {
					exitOnError(err)
				}
			}
		}
	}
}

func genLogic() logic {
	trash := os.Getenv("XDG_TRASH_HOME")
	ffolder := filepath.Join(trash, "files")
	ifolder := filepath.Join(trash, "info")

	exists([]string{trash, ffolder, ifolder})
	return logic{
		fileFolder: ffolder,
		infoFolder: ifolder,
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
		FILES[idx].infoPWD = filepath.Join(logistics.infoFolder, FILES[idx].alias+".trashinfo")
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

func rename(name string, path string, newName string, count int) string {
	if _, err := os.Stat(filepath.Join(path, newName)); err != nil {
		if os.IsNotExist(err) {
			return newName
		}
	}
	return rename(name, path, fmt.Sprintf("%v.%v", name, count), count+1)
}

func move(filename string, fro string, to string) error {
	err := os.Rename(fro, to)
	if err != nil {
		return fmt.Errorf("%v: no such file or directory", filename)
	}
	return nil
}

func restore() {
	files := genPaths(getPWD(), os.Args[2:], false)

	for _, file := range files {
		err := move(file.name, file.filePWD, file.currentPWD)
		if err != nil {
			fmt.Println(err)
		}

		os.Remove(file.infoPWD)
	}
}

func delete() {
	files := genPaths(getPWD(), os.Args[1:], true)

	for _, file := range files {
		err := move(file.name, file.currentPWD, file.filePWD)
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.Create(file.infoPWD)
		check(err)
		defer f.Close()

		currentTime := time.Now()
		body := []byte("[Trash Info]\n" + "Path=" + file.currentPWD + "\n" + "DeletionDate=" + currentTime.Format(time.RFC3339))

		_, err = f.Write(body)
		check(err)
	}
}

func listdir() {
	path := filepath.Join(os.Getenv("XDG_TRASH_HOME"), "files")

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// TODO(#5): find a way to format cleaner output
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func parseFlags(flag string) {
	defer func() {
		if r := recover(); r != nil {
			delete()
		}
	}()

	flags := map[string]func(){
		"-r":        restore,
		"--restore": restore,
		"-l":        listdir,
		"--listdir": listdir,
	}

	flags[flag]()
}

func main() {
	if (len(os.Args)) < 2 {
		exitOnError(fmt.Errorf("Error: invalid number of arguments [%v]", len(os.Args)))
	}

	parseFlags(os.Args[1])
}
