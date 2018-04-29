package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type file struct {
	name          string
	isHtml, isDir bool
}

func (f file) String() string {
	return fmt.Sprintf("[%v]\nisHtml: %v\nisDir: %v\n", f.name, f.isHtml, f.isDir)
}

func main() {
	var root string

	if len(os.Args) == 1 {
		root = "."
	} else {
		root = os.Args[1]
	}

//	files, _ := readDir(root)

//	for _, file := range files {
//		fmt.Printf("%v\n\n", file)
//	}

	buildIndices(root)
}

func buildIndices(dir string) (err error) {
	files, err := readDir(dir)
	if err != nil {
		return
	}
	if len(files) == 1 && files[0].isHtml {
		return
	}
	index := fmt.Sprintf("%s/index.html", dir)
	for _, f := range files {
		switch {
		case f.isHtml:
			appendHyperlinkToFile(index, f.name)
		case f.isDir:
			var empty bool // getting a shadowed return error without this line
			if empty, err = dirIsEmpty(f.name); err != nil {
				return
			} else if empty {
				continue
			}
			var files []file
			files, err = readDir(f.name)
			if err != nil {
				return
			}
			if len(files) == 1 && files[0].isHtml {
				appendHyperlinkToFile(index, fmt.Sprintf("%s/%s", f.name, files[0].name))
				continue
			}
			appendHyperlinkToFile(index, fmt.Sprintf("%s/index.html", f.name))
			if err = buildIndices(f.name); err != nil {
				return
			}
		}
	}
	return
}

func appendHyperlinkToFile(receivingFile string, fileToHyperlink string) {
	fmt.Printf("              index = %s\nhyperlink to append = %s\n\n", receivingFile, fileToHyperlink)
}

func readDir(dir string) (files []file, err error) {
	allFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, f := range allFiles {
		if f.IsDir() {
			files = append(files, file{name: f.Name(), isDir: true})
		} else {
			name := f.Name()
			ext := strings.ToLower(name[len(name)-5:])
			if ext == ".html" || ext[1:] == ".htm" {
				files = append(files, file{name: f.Name(), isHtml: true})
			}
		}
	}
	return
}

func dirIsEmpty(dir string) (answer bool, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	answer = len(files) == 0
	return
}
